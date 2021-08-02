package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/catinello/base62"
	"github.com/goware/urlx"
)

type Shortener struct {
	Id            string `json:"shortUrl"` //base62 encoded string used to id the original url
	LongUrl       string
	collisionList []string
}

func (s *Shortener) normalize() error {
	url, err := urlx.Parse(s.LongUrl)
	if err != nil {
		return err
	}
	s.LongUrl, err = urlx.Normalize(url)
	if err != nil {
		return err
	}
	return nil
}

func (s *Shortener) hexToBase62(hash string) error {
	val, err := strconv.ParseInt(hash[:8], 16, 64)
	if err != nil {
		log.Println("error while parsing hex", err)
		return err
	}
	//shortened url
	s.Id = base62.Encode(int(val))
	return nil
}

func (s *Shortener) Calculate() string {
	h := sha256.New()
	h.Write([]byte(s.LongUrl))
	return hex.EncodeToString(h.Sum(nil))
}

type HashCalculator interface {
	Calculate() string
}

var visitCount map[string]int

//computeId computes id that will be used to construct short url. It needs url store
//and hashing function. Use nil for default implementation.
func (s *Shortener) computeId(store URLStore, hasher HashCalculator) error {

	if hasher == nil {
		hasher = s //use default implementation if hash calculator is not provided
	}
	const splitter string = "###urlshortener###"
	var origID string
rehash:
	err := s.hexToBase62(hasher.Calculate())
	if err != nil {
		log.Println("error while calculating hash for", s.LongUrl)
		return err
	}
	s.LongUrl = strings.Split(s.LongUrl, splitter)[0]
	if origID == "" {
		origID = s.Id
	}
	if v, exists := store[s.Id]; exists {
		//check if url exist in collision list
		for _, key := range v.collisionList {
			if found, exists := store[key]; exists && found.LongUrl == s.LongUrl {
				s = found
				return nil
			}
		}
		if v.LongUrl != s.LongUrl {
			//on collision
			s.LongUrl = s.LongUrl + splitter + time.Now().String()
			goto rehash
		}
		//shorturl already exists
		s = store[s.Id]
		visitCount[s.Id] = visitCount[s.Id] + 1
		if visitCount[s.Id] > 3 {
			log.Println("same url hit more than 3 times")
		}
		return nil
	}
	if origID != s.Id {
		(*store[origID]).collisionList = append((*store[origID]).collisionList, s.Id)
	}
	store[s.Id] = s
	return nil
}
