package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const urlStorePath string = "url_store.db"

type URLStore map[string]*Shortener

func (store URLStore) Backup() error {
	var buffer bytes.Buffer
	for _, s := range store {
		line := s.Id + "["
		for _, v := range s.collisionList {
			line = line + v + ","
		}
		line = line + "]:" + s.LongUrl + "\n"
		buffer.WriteString(line)
	}

	err := ioutil.WriteFile(urlStorePath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (store URLStore) Load() {
	f, err := os.OpenFile(urlStorePath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("error reading data store from file system", err)
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		s := Shortener{
			Id:      strings.SplitN(line, "[", 2)[0],
			LongUrl: strings.SplitN(line, ":", 2)[1],
		}
		shortUrls := strings.SplitN(strings.SplitN(line, "[", 2)[1], "]", 2)[0]
		for _, v := range strings.Split(shortUrls, ",") {
			s.collisionList = append(s.collisionList, v)
		}
		store[s.Id] = &s
	}
}
