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
	
//	"encoding/json"

)


func init() {
    
    http.HandleFunc("/", index)
	http.HandleFunc("/mainpage", mainPage)
	http.HandleFunc("/getcm", getCM)

	http.HandleFunc("/getinfo", getInfo)
	http.HandleFunc("/inputdatapts", inputPasien)
	http.HandleFunc("/getlist", listPasien)
	
	http.HandleFunc("/getiki", listIKI)
	//http.HandleFunc("/testdb", testdb)

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
   TanggalFinal    string
}

func ubahBulanIni(d int) time.Time{
   y, m, _ := time.Now().Date()
   bulan := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
   return bulan

}

/*
func countIKI(tgl time.Time, counter int, iki string)(iki1, iki2, final string){
   datatgl := tgl
   hariIni := ubahBulanIni(i)
   
   if hariIni != datatgl {
      Counter++
	     if iki == "1" {
		    Iki1 = 0
			Iki1++
		 }else{
		    if iki == "2" {
			Iki2 = 0
			Iki2++
			}
		 }
	if iki
   }
   
   
   

} */

func ubahTanggal(tgl time.Time, shift string) string{
   
   ubah := tgl   
   jam := ubah.Hour()

   if jam < 12 && shift == "3"{
	     ubah = tgl.AddDate(0,0,-1)
         }
	final := ubah.Format("02-01-2006")
	return final
}

func renderPasien(w http.ResponseWriter, data interface{}, tmp string ){
      tmpl, err := template.New("tempPasien").Parse(tmp)
	  if err != nil {
	  fmt.Fprint(w, "Error Parsing: %v", err)
	  }
	  tmpl.Execute(w, data)
}

func CreateTime() time.Time{
   t := time.Now()
   zone, err := time.LoadLocation("Asia/Makassar")
   if err != nil{
      fmt.Println("Err: ", err.Error())
   }
   jam :=t.In(zone)
   return jam  
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
   
   data := &DataPasien{
      NamaPasien: r.FormValue("namapts"),
   }
   
   kun := &KunjunganPasien{
	  Diagnosis: r.FormValue("diag"),
	  GolIKI: r.FormValue("iki"),
	  ATS: r.FormValue("ats"),
	  ShiftJaga: r.FormValue("shift"),
	  JamDatang: CreateTime(),
	  Dokter: doc,
	  LinkID: pasienKey.Encode(),
   }

  
   if PasienAda == false {
       if _, err := datastore.Put(ctx, parentKey, data);err != nil{
            fmt.Fprint(w, "Error Database: %v", err)
		    return
	     }
       if _, err := datastore.Put(ctx, pasienKey, kun); err != nil {
           fmt.Fprint(w, "Error Database: %v", err)
	       return
         }
      }else{
	   if _, err := datastore.Put(ctx, pasienKey, kun); err != nil {
           fmt.Fprint(w, "Error Database: %v", err)
	       return
         }
	 
   }
 
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

//Fungsi untuk mendapatkan list pasien bulan ini
func listPasien(w http.ResponseWriter, r *http.Request){
   ctx := appengine.NewContext(r)
   
   u := user.Current(ctx)
   email := u.Email
   item := `<tr>
            
		 	<td class="text-right"><div class="checkbox">
		 		<label><input type="checkbox" name="itemkey" id="" value="{{.LinkID}}"></label>
		 	</div></td>
		 	<td class="text-right">{{.TanggalFinal}}</td>
		 	<td class="text-right">{{.NomorCM}}</td>
		 	<td class="text-left">{{.NamaPasien}}</td>
		 	<td class="text-left">{{.Diagnosis}}</td>
		 	{{template "iki"}}
		 </tr>`
   q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Order("-JamDatang").Limit(30)
   
   t := q.Run(ctx)
   for {
      var r ListPasien

	  k, err := t.Next(&r)
	  if err == datastore.Done{
	     break
	  }
	  if err != nil{
	     fmt.Fprint(w, "Error Fetching Data: ", err)
		 break
	  }
	  
	  jam := ubahTanggal(r.JamDatang, r.ShiftJaga)
	  r.TanggalFinal = jam
	  nocm := k.Parent()
	  r.NomorCM = nocm.StringID()
	  
	  nm:= datastore.NewQuery("DataPasien").Ancestor(nocm).Project("NamaPasien")
	  c := nm.Run(ctx)
	  var nama DataPasien
	  _, err = c.Next(&nama)
	  if err == datastore.Done{
		 break
	  }
	  
	  if err != nil{
	     fmt.Fprint(w, "Error Cannot Resolve Query :", err)
		 break
	  }
	  
	  r.NamaPasien = nama.NamaPasien

	  tmpl, err := template.New("tempPasien").Parse(item)
	  if err != nil {
	     fmt.Fprint(w, "Error Parsing: %v", err)
	     }
	  var tmpliki string
	  if r.GolIKI == "1"{
	  	tmpliki = `
        {{define "iki"}}
			<td class="text-center">&#x2714;</td>
			<td class="text-center"></td>
			{{end}}
			`
	  }else{
	    tmpliki = `
		{{define "iki"}}
		   	<td class="text-center"></td>
			<td class="text-center">&#x2714;</td>
     	{{end}}
       `	
	  }
	  
	  tmplend, err := template.Must(tmpl.Clone()).Parse(tmpliki)
	  if err != nil {
	     fmt.Fprint(w, "Error Parsing Second template :", err)
	  }
	  
	     tmplend.Execute(w, r)
	  
   }

   	
}

func listIKI(w http.ResponseWriter, r *http.Request){
   ctx := appengine.NewContext(r)
   
   u := user.Current(ctx)
   email := u.Email
   awalBulan := ubahBulanIni(1)
   
   q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Project("JamDatang", "GolIKI", "ShiftJaga").Order("JamDatang").Limit(300)
   
   t := q.Run(ctx)
   result := make(map[int]KunjunganPasien)   
   i := 1
   for{
   
      var iki KunjunganPasien
	  _, err := t.Next(&iki)
	  if err == datastore.Done{
		 break
	  }
	  if err != nil{
	     fmt.Fprint(w, "Cannot Read Data: ", err)
		 break
	  }
	  
	  jam := ubahTanggal(iki.JamDatang, iki.ShiftJaga)
	  wkt, _ := time.Parse("02-01-2006", jam)
	  iki.JamDatang = wkt
	  if wkt.Before(awalBulan) == true {continue}
	  result[i] = iki
	  i++
	}
	
	
	for j := 1;j<=len(result);j++{
	   fmt.Fprint(w, result[j].GolIKI)
	   fmt.Fprint(w, result[j].JamDatang)
	   fmt.Fprintln(w, result[j].ShiftJaga)
	
	}
}