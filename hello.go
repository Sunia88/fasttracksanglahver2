package main

import (
	"html/template"
    "net/http"
	


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

/*
type DataPasien struct {
   NamaPasien   string
   NomorCM      string
}

type DaftarPasien struct {
   NomorCM      string
   Diagnosis    string
   GolIKI       int
   ATS          int
   JamDatang    time.Time
   ShiftJaga    int
}

func tambahData(ctx context.Context) {
   
}

func getCM(w http.ResponseWriter, r *http.Request){
   
}
*/

func mainPage(w http.ResponseWriter, r *http.Request) {
    main, _ := template.ParseFiles("templates/base.html", "templates/main.html")
	main.Execute(w, nil)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
   login, _ := template.ParseFiles("templates/base.html", "templates/index.html")
   login.Execute(w, nil)

}