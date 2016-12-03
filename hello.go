package main

import (
	"html/template"
    "net/http"
	"time"
	

	"appengine"
	"appengine/datastore"
	"appengine/user"
//    "appengine/memcache"
	
    "fmt"

)


func init() {
    
    http.HandleFunc("/", index)
	http.HandleFunc("/mainpage", mainPage)
	http.HandleFunc("/getcm", getCM)

	http.HandleFunc("/getinfo", getInfo)
	http.HandleFunc("/inputdatapts", inputPasien)
//	http.HandleFunc("/getlist", listPasien)

}

var PasienAda bool = false


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

type ListPasien struct {
   DataPasien
   KunjunganPasien
}

func ubahTanggal(tgl time.Time, shift string) string{
   jam := tgl.Hour()
   
   if jam < 12 && shift == "3"{
	     tgl.AddDate(0,0,-1)
		 }
   final := tgl.Format("02-01-2006")
   return final
}

func renderPasien(w http.ResponseWriter, data interface{}, tmp string ){
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

/*   
   html := `
      <tr>
      	<td id="notabel"></td>
      	<td>{{.JamDatang}}</td>
      	<td>{{.NomorCM}}</td>
      	<td>{{.NamaPasien}}</td>
      	<td>{{.Diagnosis}}</td>
      	<td>{{.ATS}}</td>
      	<td>{{.GolIKI}}</td>
      </tr>
   `
*/
   ctx := appengine.NewContext(r)
   
   u := user.Current(ctx)
   doc := u.Email
   
   nocm := r.FormValue("nocm")
   grandParentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   parentKey := datastore.NewKey(ctx, "DataPasien", nocm, 0, grandParentKey)
   pasienKey := datastore.NewIncompleteKey(ctx, "KunjunganPasien", parentKey)
   
   data := &DataPasien{
      NamaPasien: r.FormValue("namapts"),
   }
   loc, err := time.LoadLocation("Asia/Makassar") 
   if err != nil{
      fmt.Println("err: ", err.Error())
   }
   
   kun := &KunjunganPasien{
	  Diagnosis: r.FormValue("diag"),
	  GolIKI: r.FormValue("iki"),
	  ATS: r.FormValue("ats"),
	  ShiftJaga: r.FormValue("shift"),
	  JamDatang: time.Now().In(loc),
	  Dokter: doc,
	  LinkID: pasienKey.Encode(),
   }


//tunda memcache dulu   
   
   var res ListPasien
   res.JamDatang = kun.JamDatang
   res.NomorCM = nocm
   res.NamaPasien = data.NamaPasien
   res.Diagnosis = kun.Diagnosis
   res.ATS = kun.ATS
   res.GolIKI = kun.GolIKI
   res.LinkID = kun.LinkID
/*   
   item1 := &memcache.Item{
      Key: res.LinkID,
	  Object: res,
   } */
   
   if PasienAda == false {
       if _, err := datastore.Put(ctx, parentKey, data);err != nil{
            fmt.Fprint(w, "Error Database: %v", err)
		    return
	     }
       if _, err := datastore.Put(ctx, pasienKey, kun); err != nil {
           fmt.Fprint(w, "Error Database: %v", err)
	       return
         }
		 
/*	   if err := memcache.Add(ctx, item1); err == memcache.ErrNotStored {
           if err := memcache.Set(ctx, item1); err != nil{
		      fmt.Printf("error setting item; %v", err)
		   }
         }*/
		 
      }else{
	   if _, err := datastore.Put(ctx, pasienKey, kun); err != nil {
           fmt.Fprint(w, "Error Database: %v", err)
	       return
         }
/*	   if err := memcache.Add(ctx, item1); err == memcache.ErrNotStored {
           if err := memcache.Set(ctx, item1); err != nil{
		      fmt.Printf("error setting item; %v", err)
		   }
	  }*/
		 
   }
/*   
   res2 := new(ListPasien)
   if item2, err := memcache.Get(ctx, res.LinkID, res2); err == memcache.ErrCacheMiss{
      fmt.Printf("Item tidak tersimpan dalam cache")
   }else if err != nil {
      fmt.Printf("Tidak bisa mengambil item: %v", err)
   } else {
      renderPasien(w, res2, html)
   }
*/   

fmt.Fprint(w, "Yeeeeee")   
//   renderPasien(w, &res, html)

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
   
   adakah := &PasienAda
   ctx := appengine.NewContext(r)
   
   nocm := r.FormValue("nocm");
   
   parentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   pasienKey := datastore.NewKey(ctx, "DataPasien", nocm, 0, parentKey)
   
   q := datastore.NewQuery("DataPasien").Ancestor(pasienKey)
   var pasien []DataPasien
   t, err := q.GetAll(ctx, &pasien)
   if err != nil {
      fmt.Fprint(w, "Error Database: %v", err)
	  return
   }
   if len(t) == 0{
   
         pts := DataPasien{}
         renderPasien(w, pts, adaPasien)
		    *adakah = false
	  }else{
   for _, pts := range pasien {
            renderPasien(w, pts, adaPasien)
			*adakah = true
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

/*
func listPasien(w http.ResponseWriter, r *http.Request){
   ctx := appengine.NewContext(r)
   
   u := user.Current(ctx)
   email := u.Email
   
   parentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   
   q := datastore.NewQuery("KunjunganPasien").Ancestor(parentKey).Filter("Dokter =", email).Limit(10).Order("-JamDatang")
   
   
   
}*/