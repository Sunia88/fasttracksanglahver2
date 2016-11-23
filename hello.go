package main

import (
	"html/template"
    "net/http"
	"time"
	

	"appengine"
	"appengine/datastore"
//	"appengine/user"


)


func init() {
    
//  http.HandleFunc("/", mainPage)
	http.HandleFunc("/login", loginPage)
//	http.HandleFunc("/getcm", getCM)
    http.HandleFunc("/tambahdokter", tambahDataDokter)

}


type Dokter struct {
   Username     string
   NamaLengkap  string
   Email        string
   Password     string
   TglDaftar    time.Time
}
/*
type DataPasien struct {
   NamaPasien   string
   NomorCM      string
   TglDaftar    time.Time
}

type KunjunganPasien struct {
   NomorCM      string
   Diagnosis    string
   GolIKI       int
   ATS          int
   JamDatang    time.Time
   ShiftJaga    int
   Dokter       *datastore.Key
}
*/


/*
//entity group fasttrack
func parentKey(ctx context.Context) *datastore.Key{
   ctx := appengine.NewContext(r *http.Request)
   return datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)

}


//setiap KunjunganPasien akan ada data dokter yang merawat
func dokterKey(ctx context.Context, username string) *datastore.Key {
   
   return datastore.NewKey(ctx, "Dokter", username, 0, parentKey)
   
}

//setiap KunjunganPasien akan disimpan dibawah DataPasien
func pasienKey(ctx context.Context, noCM string) *datastore.Key {

   return datastore.NewKey(ctx, "DataPasien", noCM, 0, parentKey)
   
   
}

*/
func tambahDataDokter(w http.ResponseWriter, r *http.Request) {

   if r.Method != "POST" {
      http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	  return
   }
   
   ctx := appengine.NewContext(r)

   dr := &Dokter{
      Username: r.FormValue("username"),
	  NamaLengkap: r.FormValue("namalengkap"),
	  Email: r.FormValue("email"),
	  Password: r.FormValue("pwd1"),
	  TglDaftar: time.Now(),
   }
   
   parentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   
   dok := datastore.NewKey(ctx, "Dokter", dr.Username, 0, parentKey)
   if _, err := datastore.Put(ctx, dok, dr); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
   }
   
   http.Redirect(w, r, "/login", http.StatusSeeOther)
}
/*
func getCM

func mainPage(w http.ResponseWriter, r *http.Request) {

    ctx := appengine.NewContext(r)
	
    main, _ := template.ParseFiles("templates/base.html", "templates/main.html")
	main.Execute(w, nil)
}


*/


func loginPage(w http.ResponseWriter, r *http.Request) {
   login, _ := template.ParseFiles("templates/base.html", "templates/registrasi.html")
   login.Execute(w, nil)

}