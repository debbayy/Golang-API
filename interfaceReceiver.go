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
	Telp	 string	`json:"telp"`
}

type KaryawanHandler struct {
	db *sql.DB
}

func (kh *KaryawanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		kh.getKaryawan(w, r)
	case "POST":
		kh.newKaryawan(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (kh *KaryawanHandler) getKaryawan(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var karyawan Karyawan

	err := kh.db.QueryRow("SELECT id, name, email, password,telp FROM karyawan WHERE id = ?", id).Scan(&karyawan.ID, &karyawan.Name, &karyawan.Email, &karyawan.Password,&karyawan.Telp)
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

	result, err := stmt.Exec(karyawan.Name, karyawan.Email, karyawan.Password,karyawan.Telp)
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
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	karyawanHandler := &KaryawanHandler{db}

	http.Handle("/karyawan", karyawanHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
