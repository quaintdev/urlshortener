package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerShortenRequest(t *testing.T) {
	jsonStr := []byte(`{"LongUrl":"https://marketplace.visualstudio.com/items?itemName=humao.rest-client"}`)
	req, err := http.NewRequest("POST", "/shortenurl", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	urlStore := make(URLStore)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleShortenRequest(urlStore))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expectedShortener, actualShortener Shortener

	expected := `{"LongUrl":"https://marketplace.visualstudio.com/items?itemName=humao.rest-client","ShortUrl":"3YLBCD"}`
	json.Unmarshal([]byte(expected), &expectedShortener)
	json.NewDecoder(rr.Body).Decode(&actualShortener)
	if expectedShortener.Id != actualShortener.Id {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleShortenUrl(t *testing.T) {

	urlStore := make(URLStore)
	urlStore["3YLBCD"] = &Shortener{
		LongUrl: "https://marketplace.visualstudio.com/items?itemName=humao.rest-client",
	}

	rr := httptest.NewRecorder()

	req1, err := http.NewRequest("GET", "/3YLBCD", nil)
	if err != nil {
		t.Fatal(err)
	}
	req1.RequestURI = "/3YLBCD"
	handler := http.HandlerFunc(handleShortUrl(urlStore))
	handler.ServeHTTP(rr, req1)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}

func Test_computeId_id_doesNotExist(t *testing.T) {
	s := Shortener{
		LongUrl: "https://marketplace.visualstudio.com/items?itemName=humao.rest-client",
	}
	store := make(URLStore)
	s.computeId(store, nil)
	if s.Id != "3YLBCD" {
		t.Errorf("actual:%v", s)
		t.FailNow()
	}
	if _, exists := store[s.Id]; !exists {
		t.Errorf("store: %v", store)
		t.FailNow()
	}
}

func Test_computeId_id_Exist(t *testing.T) {
	s := Shortener{
		LongUrl: "https://marketplace.visualstudio.com/items?itemName=humao.rest-client",
	}
	store := make(URLStore)
	store["3YLBCD"] = &Shortener{
		LongUrl: s.LongUrl,
		Id:      "3YLBCD",
	}
	s.computeId(store, nil)
	if s.Id != "3YLBCD" {
		t.Errorf("actual:%v", s)
		t.FailNow()
	}
	if _, exists := store[s.Id]; !exists {
		t.Errorf("store: %v", store)
		t.FailNow()
	}
}

type DummyHasher struct {
	*Shortener
}

func (d *DummyHasher) Calculate() string {
	if strings.Contains(d.LongUrl, "###urlshortener###") {
		return "1234567890abcdef"
	}
	return "112233abcdef"
}

func Test_computeId_on_collision(t *testing.T) {
	store := make(URLStore)
	s := Shortener{
		LongUrl: "https://marketplace.visualstudio.com/items?itemName=humao.rest-client",
	}
	s1 := Shortener{
		LongUrl: "https://github.com/drone/drone",
	}
	s2 := Shortener{
		LongUrl: "https://github.com/drone/drone",
	}

	d1 := DummyHasher{
		Shortener: &s,
	}
	d2 := DummyHasher{
		Shortener: &s1,
	}
	d3 := DummyHasher{
		Shortener: &s2,
	}

	s.computeId(store, &d1)
	s1.computeId(store, &d2)
	s2.computeId(store, &d3)
	if len(store) != 2 && len(store[s.Id].collisionList) != 1 {
		t.FailNow()
	}
}
