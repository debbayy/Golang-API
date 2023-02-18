package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Karyawan struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Telp     string `json:"telp"`
}

type KaryawanHandler struct {
	db *sql.DB
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	getKaryawan(w http.ResponseWriter, r *http.Request)
	newKaryawan(w http.ResponseWriter, r *http.Request)
	getKaryawans(w http.ResponseWriter, r *http.Request)
}

func (kh *KaryawanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		id := r.URL.Query().Get("id")
		if id == "" {
			kh.getKaryawans(w, r)
		} else {
			kh.getKaryawan(w, r)
		}
	case "POST":
		kh.newKaryawan(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (kh *KaryawanHandler) getKaryawan(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var karyawan Karyawan

	err := kh.db.QueryRow("SELECT id, name, email, password,telp FROM karyawan WHERE id = ?", id).Scan(&karyawan.ID, &karyawan.Name, &karyawan.Email, &karyawan.Password, &karyawan.Telp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(karyawan)
}

func (kh *KaryawanHandler) newKaryawan(w http.ResponseWriter, r *http.Request) {
	var karyawan Karyawan

	err := json.NewDecoder(r.Body).Decode(&karyawan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := kh.db.Prepare("INSERT INTO karyawan (name, email, password,telp) VALUES (?, ?, ?,	?	)")
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

func (kh *KaryawanHandler) getKaryawans(w http.ResponseWriter, r *http.Request) {
	fmt.Println("masuk")

	rows, err := kh.db.Query("SELECT * from karyawan")
	fmt.Println(err, "ini error 1")

	if err != nil {
		fmt.Println("masuk nill")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var karyawanList []Karyawan
	for rows.Next() {
		var karyawan Karyawan
		fmt.Println(karyawan, "ini apaaan sad karyaq=wasa")
		err := rows.Scan(&karyawan.ID, &karyawan.Name, &karyawan.Email, &karyawan.Password, &karyawan.Telp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		karyawanList = append(karyawanList, karyawan)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode karyawanList as JSON and send it in the response
	json.NewEncoder(w).Encode(karyawanList)
}

func main() {
	// Connect to database
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Create KaryawanHandler instance
	karyawanHandler := &KaryawanHandler{db}

	// Register KaryawanHandler to handle requests to /karyawan
	http.Handle("/karyawan", karyawanHandler)

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))

}
