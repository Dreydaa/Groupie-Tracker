package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

const URL = "https://groupietrackers.herokuapp.com/api"

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	var data JSON
	resp, err := http.Get(URL)
	if err != nil {
		panic("erreur 500.css")
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&data)
	url := data.Artists
	Art, err := http.Get(url)
	if err != nil {
		panic("Une erreur est survenue:,JERem:")
	}
	defer Art.Body.Close()
	var Artists []Artists
	err = json.NewDecoder(Art.Body).Decode(&Artists)
	err = tmpl.Execute(w, Artists)
	if err != nil {
		panic("Une erreur est survenue,Jer:" + err.Error())
	}
}
func PageArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistID := r.FormValue("ID")
	fmt.Println("artistsID:", artistID)
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("static/band.html"))
		var data JSON
		resp, err := http.Get(URL)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&data)
		urlArtist := data.Artists + "/" + artistID
		Art, err := http.Get(urlArtist)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		}
		defer Art.Body.Close()
		var Locations Locations
		urlLocations := data.Locations + "/" + artistID
		Loc, err := http.Get(urlLocations)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		}
		err = json.NewDecoder(Loc.Body).Decode(&Locations)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		}
		defer Loc.Body.Close()
		var artist Artist
		err = json.NewDecoder(Art.Body).Decode(&artist)
		urlRelation := "https://groupietrackers.herokuapp.com/api/relation/" + artistID
		Rel, err := http.Get(urlRelation)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		} else {
			fmt.Println("Tentative d'ouverture:", urlRelation)
		}
		var relation Relation
		err = json.NewDecoder(Rel.Body).Decode(&relation)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		}
		defer Rel.Body.Close()
		var outputBand OutputBand
		outputBand.ID = artist.ID
		outputBand.Image = artist.Image
		outputBand.FirstAlbum = artist.FirstAlbum
		outputBand.Name = artist.Name
		outputBand.Members = artist.Members
		outputBand.CreationDate = artist.CreationDate
		outputBand.Locations = Format_Locations_From_Array(Locations.Lcs)
		outputBand.RelationA = Format_Date(relation.DatesLocations)
		err = tmpl.Execute(w, outputBand)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		}
	case "POST":
		//do nothing
	default:
		http.Redirect(w, r, "400", 400)
	}
}

func Format_Date(m map[string][]string) []string {
	/*
		Utilisé pour formater vers une chaîne les concerts avec les dates
	*/
	a := make([]string, 0)
	var line string
	for k, v := range m {
		line = "\n" + Format_Location_From_String(k) + " : "
		for i, d := range v {
			line += d
			if i < len(v)-1 {
				line += "  " // Ajoutez un espace entre les dates
			}
		}
		a = append(a, line)
	}
	return a
}

func Format_Location_From_String(s string) string {
	var result string = "   "
	var isSpaceBefore = false
	result = "  "
	for index, char := range s {
		if index == 0 {
			result += string(rune(char - 32))
		} else if isSpaceBefore {
			result += string(rune(char - 32))
			isSpaceBefore = false
		} else if char != '_' && char != '-' {
			result += string(char)
		} else if char == '_' || char == '-' {
			result += "   "
			isSpaceBefore = true
		}
	}
	return result
}
func Format_Locations_From_Array(l []string) []string {
	var result string = "    "
	var isSpaceBefore = false
	for indexE, element := range l {
		result = "   "
		for index, char := range element {
			if index == 0 {
				result += string(rune(char - 32))
			} else if isSpaceBefore {
				result += string(rune(char - 32))
				isSpaceBefore = false
			} else if char != '_' && char != '-' {
				result += string(char)
			} else if char == '_' || char == '-' {
				result += "   "
				isSpaceBefore = true
			}
		}
		l[indexE] = result
	}
	return l
}
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl := template.Must(template.ParseFiles("templates/error.html"))
	tmpl.Execute(w, nil)
}
