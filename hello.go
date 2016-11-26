package main

import (
	"html/template"
    "net/http"
	"time"
	

	"appengine"
	"appengine/datastore"
	"appengine/user"


)


func init() {
    
    http.HandleFunc("/", index)
	http.HandleFunc("/mainpage", mainPage)
	http.HandleFunc("/registrasi", lamanRegistrasi)
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
type CurrentUser struct {
   Logout       string
   NamaLengkap  string
}

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

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}){
   t, _ := template.ParseFiles("templates/base.html", "templates/"+tmpl+".html")
   t.Execute(w, p)
}

/*
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
   
   http.Redirect(w, r, "/", http.StatusSeeOther)
}

//func getCM

func index(w http.ResponseWriter, r *http.Request) {
   if r.Method != "GET" {
      http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
	  return
   }
   
   if r.URL.Path != "/" {
      http.NotFound(w, r)
	  return
   }
   
   ctx := appengine.NewContext(r)

   
   var login string
   
   if u := user.Current(ctx); u != nil {
      http.Redirect(w, r, "/mainpage", http.StatusSeeOther)
	   
   } else {
      login, _ = user.LoginURL(ctx, "/mainpage")
      http.Redirect(w, r, login, http.StatusSeeOther)
      
   }

}

func lamanRegistrasi(w http.ResponseWriter, r *http.Request) {
   renderTemplate(w, "registrasi", nil)

}

func mainPage(w http.ResponseWriter, r *http.Request){

      type person struct {
	     NamaLengkap   string
		 Logout        string
	  }
      //ctx := appengine.NewContext(r)
      //parentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
	  /*
      var logout, email string
      u := user.Current(ctx)	  
	  logout, _ = user.LogoutURL(ctx, "/")
	  email = u.Email
	  */
	  cur := person{
	     NamaLengkap: "I Wayan Surya Sedana",
	     Logout: "logout",
	  }

   //Problem: cara menambahkan logout ke type CurrentUser	 
	  /*  
	  q := datastore.NewQuery("Dokter").Ancestor(parentKey).Filter("Email =", email).Project("NamaLengkap")
	  
	  res := q.Run(ctx)
	  
	  for {
		 _, err := res.Next(&cur)
		 if err == datastore.Done {
		    break
		 }
	  }
	  cur.Logout = logout */
	  
   renderTemplate(w, "main", cur)
}

func loginPage(w http.ResponseWriter, r *http.Request){
   renderTemplate(w, "login", nil)
}