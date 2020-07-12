package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.HandleFunc("/currents", Recently).Methods(http.MethodPost)
	r.HandleFunc("/checkin", CheckIn).Methods(http.MethodPost)
	r.HandleFunc("/checkout", CheckOut).Methods(http.MethodPost)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type Check struct {
	ID      int64 `json:id`
	PlaceID int64 `json:place_id`
}

type Location struct {
	Lat  float64
	Long float64
}

// Recently returns currently visited
func Recently(w http.ResponseWriter, r *http.Request) {

}

var InsertCheckIn = func(ID, placeID int64) error {
	db, err := sql.Open("sqlite3", "thaichana.db")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO visits VALUES(?, ?);", ID, placeID)
	return err
}

// CheckIn check-in to place, returns density (ok, too much)
func CheckIn(w http.ResponseWriter, r *http.Request) {
	var chk Check
	if err := json.NewDecoder(r.Body).Decode(&chk); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer r.Body.Close()

	if err := InsertCheckIn(chk.ID, chk.PlaceID); err == nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
}

// CheckOut check-out from place
func CheckOut(w http.ResponseWriter, r *http.Request) {

}
