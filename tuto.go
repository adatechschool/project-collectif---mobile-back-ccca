package main

import(
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
    "encoding/json"
	"github.com/gorilla/mux"
)

type event struct{
	ID string `json:"ID"`
	Title string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID: "1",
		Title: "Intro to Golang",
		Description: "Go is better with maple syrup",
	},

}

func createEvent(w http.ResponseWriter, r *http.Request){
	var newEvent event
	reqBody, err:= ioutil.ReadAll(r.Body)
	if err != nil{
		fmt.Fprintf(w, "ça marche pas")
	}
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func homeLink(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome Home !")
}

func main(){
	router:= mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))
}