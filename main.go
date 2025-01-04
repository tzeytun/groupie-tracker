package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		ArtistName     string              `json:"artistName"`
		DatesLocations map[string][]string `json:"datesLocations"`
		Image          string              `json:"image"`
	} `json:"index"`
}

func main() {
	http.HandleFunc("/", artistHandler)
	http.HandleFunc("/relation", RelationHandler)
	http.ListenAndServe(":8080", nil)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Getirme hatası", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		http.Error(w, "JSON dönüştürme hatası", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("static/main.html")
	if err != nil {
		http.Error(w, "Template dosyasını açma hatası", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, "Template yürütme hatası", http.StatusInternalServerError)
		return
	}
}

func RelationHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/relation"

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Getirme hatası", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var cal Relation
	err = json.NewDecoder(response.Body).Decode(&cal)
	if err != nil {
		http.Error(w, "JSON dönüştürme hatası", http.StatusInternalServerError)
		return
	}

	artistURL := "https://groupietrackers.herokuapp.com/api/artists"

	artistResponse, err := http.Get(artistURL)
	if err != nil {
		http.Error(w, "Getirme hatası", http.StatusInternalServerError)
		return
	}
	defer artistResponse.Body.Close()

	var artists []Artist
	err = json.NewDecoder(artistResponse.Body).Decode(&artists)
	if err != nil {
		http.Error(w, "JSON dönüştürme hatası", http.StatusInternalServerError)
		return
	}

	for i := range cal.Index {
		for _, artist := range artists {
			if cal.Index[i].ID == artist.ID {
				cal.Index[i].ArtistName = artist.Name
				break
			}
		}
	}
	for i := range cal.Index {
		for _, artist := range artists {
			if cal.Index[i].ID == artist.ID {
				cal.Index[i].Image = artist.Image
				break
			}
		}
	}

	tmpl, err := template.ParseFiles("static/relation.html")
	if err != nil {
		http.Error(w, "Template dosyasını açma hatası", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, cal)
	if err != nil {
		http.Error(w, "Template yürütme hatası", http.StatusInternalServerError)
		return
	}
}
