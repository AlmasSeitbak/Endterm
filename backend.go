package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Barber struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func main() {
	// Open database connection
	db, err := sql.Open("sqlite3", "./barbershops.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create barber table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS barber (id INTEGER PRIMARY KEY, name TEXT, address TEXT, phone TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Create router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/barber", getBarbers(db)).Methods("GET")
	router.HandleFunc("/barber/{id}", getBarber(db)).Methods("GET")
	router.HandleFunc("/barber", createBarber(db)).Methods("POST")
	router.HandleFunc("/barber/{id}", updateBarber(db)).Methods("PUT")
	router.HandleFunc("/barber/{id}", deleteBarber(db)).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all barbers
func getBarbers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM barber")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		barbers := []Barber{}
		for rows.Next() {
			var b Barber
			if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.Phone); err != nil {
				log.Fatal(err)
			}
			barbers = append(barbers, b)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(barbers)
	}
}

// Get a single barber
func getBarber(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			log.Fatal(err)
		}

		var b Barber
		row := db.QueryRow("SELECT * FROM barber WHERE id = ?", id)
		if err := row.Scan(&b.ID, &b.Name, &b.Address, &b.Phone); err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(b)
	}
}

// Create a new barber
func createBarber(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b Barber
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			log.Fatal(err)
		}

		result, err := db.Exec("INSERT INTO barber (name, address, phone) VALUES (?, ?, ?)", b.Name, b.Address, b.Phone)
		if err != nil {
			log.Fatal(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		b.ID = int(id)

		w.Header().Set("Content-Type", "application
