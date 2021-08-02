package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"time"
)

const apiServerPort string = ":3000"
const maintenanceServerPort string = ":3001"

func main() {
	urlStore := make(URLStore)
	urlStore.Load()

	visitCount = make(map[string]int)

	maintenanceServer := http.NewServeMux()
	maintenanceServer.HandleFunc("/backup", handleBackup(urlStore))
	go func() {
		log.Println("Starting maintenance server on port", maintenanceServerPort)
		http.ListenAndServe(maintenanceServerPort, maintenanceServer)
	}()

	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				for k, _ := range urlStore {
					delete(urlStore, k)
				}
				log.Println("url store cleared at ", t)
			}
		}
	}()

	http.HandleFunc("/shorten", handleShortenRequest(urlStore))
	http.HandleFunc("/", handleShortUrl(urlStore))
	log.Println("Starting server on port", apiServerPort)
	http.ListenAndServe(apiServerPort, nil)
}

func handleShortUrl(store URLStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.RequestURI[1:]
			if _, ok := store[id]; ok {
				http.Redirect(w, r, store[id].LongUrl, http.StatusSeeOther)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleShortenRequest(store URLStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var shortener Shortener
			err := json.NewDecoder(r.Body).Decode(&shortener)
			if err != nil {
				log.Println("error decoding shorten request json", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = shortener.normalize()
			if err != nil {
				log.Println("error while trying to normalize,", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = shortener.computeId(store, nil)
			if err != nil {
				log.Println("error: failed to generate short url,", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var response Shortener
			response = shortener
			response.Id = path.Join(r.Host, shortener.Id)
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				log.Println("error while marshalling json", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleBackup(store URLStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if err := store.Backup(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("error while taking backup", err)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
