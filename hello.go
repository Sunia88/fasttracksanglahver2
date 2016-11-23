package main

import (
	"html/template"
    "net/http"
	
	"time"
	
	"appengine"
	"appengine/datastore"
	"appengine/user"


)

/*
	"appengine"
	"appengine/datastore"
	"appengine/user"
*/

func init() {
    http.HandleFunc("/", mainPage)
	http.HandleFunc("/login", loginPage)
//	http.HandleFunc("/getcm", getCM)

}

type Dokter struct {
   Username     string
   NamaLengkap  string
   Email        string
   Password     string
   TglDaftar    time.Time
}

type DataPasien struct {
   NamaPasien   string
   NomorCM      string
   TglDaftar    time.Time
}

type DaftarPasien struct {
   NomorCM      string
   Diagnosis    string
   GolIKI       int
   ATS          int
   JamDatang    time.Time
   ShiftJaga    int
}

func tambahDataDokter(w http.ResponseWriter, r *http.Request) {
   ctx := appengine.NewContext(r)

   k := datastore.IncompleteKey(ctx, "Dokter", nil)   
}

func getCM(w http.ResponseWriter, r *http.Request){
   
}


func mainPage(w http.ResponseWriter, r *http.Request) {
    main, _ := template.ParseFiles("templates/base.html", "templates/main.html")
	main.Execute(w, nil)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
   login, _ := template.ParseFiles("templates/base.html", "templates/index.html")
   login.Execute(w, nil)

}