package main

import(
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
    "encoding/json"
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func test() {
    // Open up our database connection.
    db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/surfspot")

    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    }

    // defer the close till after the main function has finished
    // executing
    defer db.Close()

	// perform a db.Query insert
	insert, err := db.Query("INSERT INTO spots VALUES ( 1, 'Pipeline', 'Reef Break', 4, false, 'Oahu, Hawaii', 'Pipeline, Oahu, Hawaii', 'https://magicseaweed.com/Pipeline-Backdoor-Surf-Report/616/', 'https://dl.airtable.com/ZuXJZ2NnTF40kCdBfTld_thomas-ashlock-64485-unsplash.jpg', '2018-07-22', '2018-08-31', '2018-05-31T00:16:16.000Z' )")

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

}

type spot struct{
	ID int `json:"ID"`
	Name string `json:"Name"`
	SurfBreak string `json:"Surf Break"`
	DifficultyLevel int `json:"Difficulty Level"`
	Favorite bool `json:"Favorite"`
	StateCountry string `json:"State/Country"`
	Address string `json:"Address"`
	Link string `json:"Link"`
	Photos string `json:"Photos"`
	SeasonStart string `json:"Season Start"`
	SeasonEnd string `json:"Season End"`
	CreatedTime string `json:"createdTime"`
}


type allSpots []spot

type partialSpot struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	SurfBreak string `json:"Surf Break"`
}

type allPartialSpots []partialSpot

var spots = allSpots{
	{
        ID: 1,
        Name: "Pipeline",
        SurfBreak: "Reef Break",
        DifficultyLevel: 4,
        Favorite: false,
        StateCountry: "Oahu, Hawaii",
        Address: "Pipeline, Oahu, Hawaii",
        Link: "https://magicseaweed.com/Pipeline-Backdoor-Surf-Report/616/",
        Photos: "https://dl.airtable.com/ZuXJZ2NnTF40kCdBfTld_thomas-ashlock-64485-unsplash.jpg",
        SeasonStart: "2018-07-22",
        SeasonEnd: "2018-08-31",
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
	spotID, _ := strconv.Atoi(mux.Vars(r)["id"])

	for i, singleSpot := range spots {
		if singleSpot.ID == spotID {
			spots = append(spots[:i], spots[i+1:]...)
			fmt.Fprintf(w, "The spot with ID %v has been deleted successfully", spotID)
		}
	}
}

func getOneSpot(w http.ResponseWriter, r *http.Request) {
	spotID, _ := strconv.Atoi(mux.Vars(r)["id"])

	for _, singleSpot := range spots {
		if singleSpot.ID == spotID {
			json.NewEncoder(w).Encode(singleSpot.ID)
		}
	}
}

func getAllSpots(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(spots)
}

func updateSpot(w http.ResponseWriter, r *http.Request) {
	spotID, _ := strconv.Atoi(mux.Vars(r)["id"])
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
	router.HandleFunc("/", getAllSpots).Methods("GET")
	router.HandleFunc("/spot", createSpot).Methods("POST")
	router.HandleFunc("/spots/{id}", deleteSpot).Methods("DELETE")
	router.HandleFunc("/spots/{id}", getOneSpot).Methods("GET")
	router.HandleFunc("/spots", getAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", updateSpot).Methods("PATCH")

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create
	}
	test()
	log.Fatal(srv.ListenAndServe())
}