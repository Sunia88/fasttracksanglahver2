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
	http.HandleFunc("/getcm", getCM)
    http.HandleFunc("/tambahdokter", tambahDataDokter)
	http.HandleFunc("/getinfo", getInfo)
	http.HandleFunc("/inputdatapts", inputPasien)

}


type Dokter struct {
   Username     string
   NamaLengkap  string
   Email        string
   Password     string
   TglDaftar    time.Time
}


type DataPasien struct {
   NamaPasien                   string
   NomorCM, JenKel, Alamat      string
   TglDaftar, Umur              time.Time
}


type KunjunganPasien struct {
   Diagnosis, LinkID            string
   GolIKI, ATS, ShiftJaga       string
   JamDatang                    time.Time
   Dokter                       string
}

func renderPasien(w http.ResponseWriter, data DataPasien, tmp string ){
      tmpl, err := template.New("tempPasien").Parse(tmp)
	  if err != nil {
	  fmt.Fprint(w, "Error Parsing: %v", err)
	  }
	  tmpl.Execute(w, data)
}
   
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}){
   t, _ := template.ParseFiles("templates/base.html", "templates/"+tmpl+".html")
   t.Execute(w, p)
}

func inputPasien(w http.ResponseWriter, r *http.Request){
   if r.Method != "POST" {
      http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	  return
   }
   
   ctx := appengine.NewContext(r)
   
   u := user.Current(ctx)
   doc := u.Email
   
   nocm := r.FormValue("nocm")
   grandParentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   parentKey := datastore.NewKey(ctx, "DataPasien", nocm, 0, grandParentKey)
   pasienKey := datastore.NewIncompleteKey(ctx, "KunjunganPasien", parentKey)
   
   q := datastore.NewQuery("DataPasien").Ancestor(grandParentKey).Filter("__key__ >", parentKey)
   if q.Count(ctx) == 0 {
   
   data := &DataPasien{
      NamaPasien: r.FormValue("namapts"),
   }
   
      if _, err := datastore.Put(ctx, parentKey, data);err != nil{
            fmt.Fprint(w, "Error Database: %v", err)
		    return
	     }
   }
   
   kun := &KunjunganPasien{
	  Diagnosis: r.FormValue("diag"),
	  GolIKI: r.FormValue("iki"),
	  ATS: r.FormValue("ats"),
	  ShiftJaga: r.FormValue("shif"),
	  JamDatang: time.Now(),
	  Dokter: doc,
	  LinkID: pasienKey.Encode(),
   }
   
   if _, err := datastore.Put(ctx, pasienKey, kun); err != nil {
      fmt.Fprint(w, "Error Database: %v", err)
	  return
   }
   
   fmt.Fprint(w, kun.LinkID)
   
   
}
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
    adaPasien := `
    <label for="namapasien">Nama Pasien:</label><br>
    {{with .NamaPasien}}<input type="text" name="namapasien" id="namapts" class="form-control text-capitalize" value={{.}}>&nbsp&nbsp<div id="errorpts"></div><br>
	{{else}}
	<input type="text" name="namapasien" id="namapts" class="form-control text-capitalize">&nbsp&nbsp<div id="errorpts"></div><br>
	{{end}}
	<label for="diagnosis">Diagnosis:</label><br>
	<input type="text" name="diagnosis" id="diag" class="form-control text-capitalize">&nbsp&nbsp<div id="errordiag"></div><br><br>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="1">ATS 1</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="2">ATS 2</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="3">ATS 3</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="4">ATS 4</label>
    <label for="" class="radio-inline"><input type="radio" name="ats" id="" value="5">ATS 5</label>&nbsp&nbsp<div id="errorats"></div><br><br>
	<label for="" class="radio-inline"><input type="radio" name="iki" id="" value="1">IKI 1</label>
	<label for="" class="radio-inline"><input type="radio" name="iki" id="" value="2">IKI 2</label>&nbsp&nbsp<div id="erroriki"></div><br><br>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="1">Pagi</label>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="2">Sore</label>
	<label for="" class="radio-inline"><input type="radio" name="shift" id="" value="3">Malam</label>&nbsp&nbsp<div id="errorshift"></div><br><br>
	
	<button type="submit" class="btn btn-primary btn-md" id="btnsub">Tambahkan Pasien</button><br><div id="errorbtn"></div>
   `

   ctx := appengine.NewContext(r)
   
   nocm := r.FormValue("nocm");
   
   parentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   pasienKey := datastore.NewKey(ctx, "DataPasien", nocm, 0, parentKey)
   
   q := datastore.NewQuery("DataPasien").Ancestor(parentKey).Filter("__key__ >", pasienKey)
   var pasien []DataPasien
   t, err := q.GetAll(ctx, &pasien)
   if err != nil {
      fmt.Fprint(w, "Error Database: %v", err)
	  return
   }
   if len(t) == 0{
         pts := DataPasien{}
         renderPasien(w, pts, adaPasien)
	  }else{
   for _, pts := range pasien {
            renderPasien(w, pts, adaPasien)
	     }
	  }
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

func mainPage(w http.ResponseWriter, r *http.Request){
   renderTemplate(w, "main", nil)
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