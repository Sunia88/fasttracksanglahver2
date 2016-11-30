package main

import (
	"html/template"
    "net/http"
	"time"
	

	"appengine"
	"appengine/datastore"
	"appengine/user"

    "fmt"

)


func init() {
    
    http.HandleFunc("/", index)
	http.HandleFunc("/mainpage", mainPage)
	http.HandleFunc("/registrasi", lamanRegistrasi)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/getcm", getCM)
    http.HandleFunc("/tambahdokter", tambahDataDokter)
	http.HandleFunc("/getinfo", getInfo)

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

func getCM(w http.ResponseWriter, r *http.Request){
  /* takAdaPasien := `
    <label for="namapasien">Nama Pasien:</label><br>
    <input type="text" name="nocm" id="nocm" class="form-control text-capitalize"><br>   
    <label for="diagnosis">Diagnosis:</label><br>
	<input type="text" name="diagnosis" id="diag" class="form-control text-capitalize"><br>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="1">ATS 1</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="2">ATS 2</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="3">ATS 3</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="4">ATS 4</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="5">ATS 5</label><br>
	<label for="" class="radio-inline"><input type="radio" name="iki" id="" value="1"></label>
	<label for="" class="radio-inline"><input type="radio" name="iki" id="" value="2"></label><br>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="1">Pagi</label>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="2">Sore</label>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="3">Malam</label><br>
	
	<button type="submit"></button>
   `
   
   adaPasien := `
    <label for="namapasien">Nama Pasien:</label><br>
    <input type="text" name="namapasien" id="nocm" class="form-control text-capitalize" value={{.NamaPasien}}><br>   
    <label for="diagnosis">Diagnosis:</label><br>
	<input type="text" name="diagnosis" id="diag" class="form-control text-capitalize"><br>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="1">ATS 1</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="2">ATS 2</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="3">ATS 3</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="4">ATS 4</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="5">ATS 5</label><br>
	<label for="" class="radio-inline"><input type="radio" name="iki" id="" value="1"></label>
	<label for="" class="radio-inline"><input type="radio" name="iki" id="" value="2"></label><br>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="1">Pagi</label>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="2">Sore</label>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="3">Malam</label><br>
	
	<button type="submit"></button>
   `
   
   ctx := appengine.NewContext(r)*/
   
   nocm := r.FormValue("nocm");
   
   fmt.Fprint(w, nocm)
   /*parentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   pasienKey := datastore.NewKey(ctx, "DataPasien", nocm, 0, parentKey)
   
   q := datastore.NewQuery("DataPasien").Ancestor(parentKey).Filter("__key__ >", pasienKey)
   
   t := q.Run(ctx)
   
   for {
      var p DataPasien
	  k, err := t.Next(&p)
	  if err == datastore.Done {
	     fmt.Fprint(w, takAdaPasien)
	  }
	  
	  if err != nil {
	     tmpl := template.Must(template.New("ada").Parse(adaPasien))
		 tmpl.Execute(w, k)
	  }
   }*/
}

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
   renderTemplate(w, "main", nil)
}

func loginPage(w http.ResponseWriter, r *http.Request){
   renderTemplate(w, "login", nil)
}

func getInfo(w http.ResponseWriter, r *http.Request){

      type Person struct {
	     NamaLengkap   string
		 Logout        string
	  }
	  
      ctx := appengine.NewContext(r)
      
      var logout, email string
      u := user.Current(ctx)	  
	  logout, _ = user.LogoutURL(ctx, "/")
	  email = u.Email
      
	  p := &Person{
	     NamaLengkap: email,
		 Logout: logout,
	  }
	  
      fmt.Fprint(w, "<p>Selamat datang "+p.NamaLengkap+"<br>Klik <a href="+p.Logout+">di sini</a> untuk Logout.")
}