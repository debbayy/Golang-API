package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Karyawan struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Telp	 string 	`json:"telp"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getKaryawan(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	id := r.URL.Query().Get("id")

	var karyawan Karyawan

	err := db.QueryRow("SELECT id, name, email, password, telp FROM karyawan WHERE id = ?", id).Scan(&karyawan.ID, &karyawan.Name, &karyawan.Email, &karyawan.Password, &karyawan.Telp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(karyawan)
}

func getKaryawans(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	rows, err := db.Query("SELECT * from karyawan")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var karyawanList []Karyawan
	for rows.Next() {
		var karyawan Karyawan
		err := rows.Scan(&karyawan.ID, &karyawan.Name, &karyawan.Email, &karyawan.Password, &karyawan.Telp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		karyawanList = append(karyawanList, karyawan)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(karyawanList)
}


func newKaryawan(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	var karyawan Karyawan

	err := json.NewDecoder(r.Body).Decode(&karyawan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO karyawan (name, email, password, telp) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(karyawan.Name, karyawan.Email, karyawan.Password, karyawan.Telp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	karyawan.ID = int(newID)

	json.NewEncoder(w).Encode(karyawan)
}


func main() {
	http.HandleFunc("/karyawan", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getKaryawan(w, r)
		case "POST":
			newKaryawan(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	//get all data karyawan
	http.HandleFunc("/karyawans", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getKaryawans(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
		
	log.Fatal(http.ListenAndServe(":8080", nil))
}