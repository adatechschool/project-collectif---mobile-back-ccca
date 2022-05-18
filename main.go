package main

import(
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
    "encoding/json"
	"github.com/gorilla/mux"
)

type spot struct{
	ID string `json:"ID"`
	Name string `json:"Name"`
	SurfBreak string `json:"Surf Break"`
	DifficultyLevel int `json:"Difficulty Level"`
	Favorite bool `json:"Favorite"`
	StateCountry string `json:"State/Country"`
	Address string `json:"Address"`
	Link string `json:"Link"`
	Photos string `json:"Photos"`
	StartSeason string `json:"Start Season"`
	EndSeason string `json:"End Season"`
	CreatedTime string `json:"createdTime"`
}

type allSpots []spot

var spots = allSpots{
	{
        ID: "rec5aF9TjMjBicXCK",
        Name: "Pipeline",
        SurfBreak: "Reef Break",
        DifficultyLevel: 4,
        Favorite: true,
        StateCountry: "Oahu, Hawaii",
        Address: "Pipeline, Oahu, Hawaii",
        Link: "https://magicseaweed.com/Pipeline-Backdoor-Surf-Report/616/",
        Photos: "https://dl.airtable.com/ZuXJZ2NnTF40kCdBfTld_thomas-ashlock-64485-unsplash.jpg",
        StartSeason: "2018-07-22",
        EndSeason: "2018-08-31",
        CreatedTime: "2018-05-31T00:16:16.000Z",
    },

}

func homeLink(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome Home !")
}

func createSpot(w http.ResponseWriter, r *http.Request){
	var newSpot spot
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "ça marche pas")
	}
	json.Unmarshal(reqBody, &newSpot)
	spots = append(spots, newSpot)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newSpot)
}

func deleteSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for i, singleSpot := range spots {
		if singleSpot.ID == spotID {
			spots = append(spots[:i], spots[i+1:]...)
			fmt.Fprintf(w, "The spot with ID %v has been deleted successfully", spotID)
		}
	}
}

func getOneSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for _, singleSpot := range spots {
		if singleSpot.ID == spotID {
			json.NewEncoder(w).Encode(singleSpot)
		}
	}
}

func getAllSpots(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(spots)
}

func updateSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]
	var updatedSpot spot

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the spot title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedSpot)

	for i, singleSpot := range spots {
		if singleSpot.ID == spotID {
			singleSpot.Name = updatedSpot.Name
			spots = append(append(spots[:i], singleSpot), spots[i+1:]...)
			json.NewEncoder(w).Encode(singleSpot)
		}
	}
}

func main(){
	router:= mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/spot", createSpot).Methods("POST")
	router.HandleFunc("/spots/{id}", deleteSpot).Methods("DELETE")
	router.HandleFunc("/spots/{id}", getOneSpot).Methods("GET")
	router.HandleFunc("/spots", getAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", updateSpot).Methods("PATCH")
	log.Fatal(http.ListenAndServe(":8080", router))
}