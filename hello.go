package main

import (
	"html/template"
    "net/http"

	"time"
//	"strconv"

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
	
	http.HandleFunc("/entri/edit/", editEntri)
	http.HandleFunc("/entri/update", updateEntri)
	http.HandleFunc("/entri/del/", deleteEntri)
	http.HandleFunc("/entri/delete", confirmDeleteEntri)

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

func ubahTanggal(tgl time.Time, shift string) time.Time{
   
   ubah := tgl   
   jam := ubah.Hour()

   if jam < 12 && shift == "3"{
	     ubah = tgl.AddDate(0,0,-1)
         }
   return ubah
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
	  TglDaftar: CreateTime(),
   }
   
   kun := &KunjunganPasien{
	  Diagnosis: r.FormValue("diag"),
	  GolIKI: r.FormValue("iki"),
	  ATS: r.FormValue("ats"),
	  ShiftJaga: r.FormValue("shift"),
	  JamDatang: CreateTime(),
	  Dokter: doc,
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

   tmpl, err := template.New("adaPasien").ParseFiles("templates/adapasien.html")
   if err != nil {
      fmt.Fprint(w, "Error Parsing Template: ", err)
   }
   
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
         tmpl.Execute(w, pts)
		    *adakah = false
	  }else{
   for _, pts := range pasien {
            tmpl.Execute(w, pts)
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
            
		 	<td class="text-right" id="number"></td>
		 	<td class="text-right">{{.TanggalFinal}}</td>
		 	<td class="text-right">{{.NomorCM}}</td>
		 	<td class="text-left text-capitalize">{{.NamaPasien}}</td>
		 	<td class="text-left text-capitalize">{{.Diagnosis}}</td>
		 	{{template "iki"}}
			<td class="text-center">
		       <div class="btn-group btn-group-xs">
			      <a href="/entri/edit/{{.LinkID}}" class="btn btn-info" role="button">Edit</a>
				  <a href="/entri/del/{{.LinkID}}" class="btn btn-danger" role="button">Delete</a>
			   </div>
			   <span id="btnval"></span>
			</td>
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
	  jamfinal := jam.Format("02-01-2006")
	  r.TanggalFinal = jamfinal
	  nocm := k.Parent()
	  r.LinkID = k.Encode()
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
   result := make(map[int]ListPasien)   
   i := 1
   for{
   
      var iki ListPasien
	  _, err := t.Next(&iki)
	  if err == datastore.Done{
		 break
	  }
	  if err != nil{
	     fmt.Fprint(w, "Cannot Read Data: ", err)
		 break
	  }
	  
	  jam := ubahTanggal(iki.JamDatang, iki.ShiftJaga)
	  if jam.Before(awalBulan) == true {continue}
	  iki.JamDatang = jam
	  iki.TanggalFinal = iki.JamDatang.Format("2-01-2006")
	  result[i] = iki
	  i++
	}
	m := ubahBulanIni(0).Day()
	

	type sumIKI struct {
	   Tanggal    string
	   IKI1, IKI2 int
	}
	
	ikiBulan := []sumIKI{}
	dataIKI := sumIKI{}
	for h:= 1; h <=m;h++{
	
	   q := ubahBulanIni(h).Format("2-01-2006")
	   var u1, u2 int
	   for _, v :=  range result{
	      if v.TanggalFinal != q {continue}
		  if v.GolIKI == "1" {
		     u1++
		  }else{
		     u2++
		  }
	   }
	   
	   if u1 == 0 && u2 == 0{continue}
	   dataIKI.Tanggal = q
	   dataIKI.IKI1 = u1
	   dataIKI.IKI2 = u2
	   ikiBulan = append(ikiBulan, dataIKI)
	}
	
	for a, b := range ikiBulan{
	
	fmt.Fprint(w, "<tr><td>")
	fmt.Fprint(w, a+1)
	fmt.Fprint(w, "</td><td>"+b.Tanggal+"</td><td>")
	fmt.Fprint(w, b.IKI1)
	fmt.Fprint(w, "</td><td>")
	fmt.Fprint(w, b.IKI2)
	fmt.Fprint(w, "</td></tr>")
    }
}

func getDatabyKey(item string, w http.ResponseWriter, r *http.Request) ListPasien {
   
   ctx := appengine.NewContext(r)
   dataKun := ListPasien{}
   keyKun, err := datastore.DecodeKey(item)
   
   if err != nil {
         fmt.Fprintln(w, "Error Decoding Key: ", err)
      }
   
   keyPts := keyKun.Parent()
   
   err = datastore.Get(ctx, keyKun, &dataKun)
   if err != nil {
      fmt.Fprintln(w, "Error Fetching Data Kunjungan: ", err)
      }
   
   err = datastore.Get(ctx, keyPts, &dataKun)
   if err != nil {
      fmt.Fprintln(w, "Error Fetching Data Pasien: ", err)
      }

   dataKun.LinkID = item

   return dataKun
}



func editEntri(w http.ResponseWriter, r *http.Request){
   keyitem := r.URL.Path[12:]
   editData := getDatabyKey(keyitem, w, r)
   renderTemplate(w, "edit", editData)
   
}

func updateEntri(w http.ResponseWriter, r *http.Request){
   if r.Method != "POST" {
      http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
	  return
   }
   
   ctx := appengine.NewContext(r)

   kun := &KunjunganPasien{}   
   pts := &DataPasien{}

   kun.LinkID = r.FormValue("entri")
   keyKun, err := datastore.DecodeKey(kun.LinkID)
   if err != nil {
      fmt.Fprintln(w, "Error Generating Key: ", err)
   }
   keyPts := keyKun.Parent()   
   
   err = datastore.Get(ctx, keyKun, kun)
   if err != nil {
      fmt.Fprintln(w, "Error Fetching Data: ", err)
	  return
   }
   kun.Diagnosis = r.FormValue("diagnosis")
   kun.ATS = r.FormValue("ats")
   kun.GolIKI = r.FormValue("iki")
   kun.ShiftJaga = r.FormValue("shift")
   
   err = datastore.Get(ctx, keyPts, pts)
   if err != nil {
      fmt.Fprintln(w, "Error Fetching Data: ", err)
	  return
   }
   pts.NamaPasien = r.FormValue("namapasien")

   
   if _, err := datastore.Put(ctx, keyKun, kun); err != nil {
      fmt.Fprint(w, "Error Putting Data Kunjungan: ", err)
	  return
   }
   
   if _, err := datastore.Put(ctx, keyPts, pts); err != nil {
      fmt.Fprint(w, "Error Putting Data Pasien: ", err)
      return
   }
   
   http.Redirect(w, r, "/mainpage", http.StatusSeeOther)
}

func deleteEntri(w http.ResponseWriter, r *http.Request){
   keyitem := r.URL.Path[11:]
   editData := getDatabyKey(keyitem, w, r)
   renderTemplate(w, "delete", editData)
}

func confirmDeleteEntri(w http.ResponseWriter, r *http.Request){
   if r.Method != "POST" {
      http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
	  return
   }
   
   ctx := appengine.NewContext(r)

   keyKun, err := datastore.DecodeKey(r.FormValue("entri"))
   if err != nil {
      fmt.Fprintln(w, "Error Generating Key: ", err)
   }
   
   var pts KunjunganPasien
   err = datastore.Get(ctx, keyKun, &pts)
   if err != nil {
     fmt.Fprintln(w, "Error Fetching Data: ", err)
	 
   }

   //fmt.Fprintln(w, r.FormValue("entri"))
   
   err = datastore.Delete(ctx, keyKun)

   http.Redirect(w, r, "/mainpage", http.StatusSeeOther)   
}