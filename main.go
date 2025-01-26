package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	Id           string    `json:"id"`
	OriginalUrl  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalUrl string) string {
	hasher := md5.New()

	//save string as bytes
	hasher.Write([]byte(OriginalUrl))
	fmt.Println("hasher : ", OriginalUrl)

	// print byte data
	fmt.Println("hasher : ", hasher)

	// calcualtes the md5 hash of the input data
	//The nil argument indicates that no additional data is being appended to the hash
	data := hasher.Sum(nil)
	fmt.Println("hasher data: ", data)

	//EncodeToString returns a string that represents the hexadecimal encoding of the data.
	hash := hex.EncodeToString(data)
	fmt.Println("Encoded to string : ", hash)

	// we return the first 8 characters only .
	// as there is almost no chance of creating the same hash again. so first 8 characters are enough
	fmt.Println("final String ", hash[:8])
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL // use the shortURL as the id for simplicity
	urlDB[id] = URL{
		Id:           id,
		OriginalUrl:  originalURL,
		CreationDate: time.Now(),
	}

	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get method")
	//Fprintf() writes output to other places except terminal . for eg- webpage, file etc
	fmt.Fprintf(w, "Hello world")
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL_ := createURL(data.URL)
	//fmt.Fprintf(w, shortURL_)

	response := struct {
		ShortUrl string `json:"short_url"`
	}{ShortUrl: shortURL_}

	//set header as response type is json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//redirect function
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)

}

func main() {
	// fmt.Println("starting url shortener")
	// generateShortURL("fdsfdsfsdfdsf")

	//the
	http.HandleFunc("/", handler)

	http.HandleFunc("/shorten", shortURLHandler)

	http.HandleFunc("/redirect/", redirectURLHandler)

	fmt.Println("starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error on starting server: ", err)
	}
}
