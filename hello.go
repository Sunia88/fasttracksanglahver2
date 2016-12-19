package main

import (
	"html/template"
    "net/http"
//	"sort"
	"time"
	"strconv"
	"strings"

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

	
	http.HandleFunc("/getiki", listIKI)
	//http.HandleFunc("/testdb", testdb)
	
	http.HandleFunc("/entri/edit/", editEntri)
	http.HandleFunc("/entri/update", updateEntri)
	http.HandleFunc("/entri/del/", deleteEntri)
	http.HandleFunc("/entri/delete", confirmDeleteEntri)
	
	http.HandleFunc("/getlaporan", listLaporan)
	http.HandleFunc("/getlaporan/", buatBCP)

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
   IKI1,IKI2 string
}

func ubahBulanIni(d int) time.Time{
   y, m, _ := time.Now().Date()
   bulan := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
   return bulan

}


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


func listLaporan(w http.ResponseWriter, r *http.Request){
   ctx := appengine.NewContext(r)
   u := user.Current(ctx)
   email := u.Email
   
   q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Project("JamDatang").Order("JamDatang")
   t := q.Run(ctx)
   var iki ListPasien
   _, err := t.Next(&iki)
   if err !=nil {
       fmt.Fprintln(w, "Error Fetching Data: ", err)
   }

   start := iki.JamDatang
   awal := start.Format("2006/01")
   dd := []string{awal}

   tgl := CreateTime().Format("2006/01")

   for awal != tgl {
      awal = start.AddDate(0, 1, 0).Format("2006/01")
      dd = append(dd, awal)
   }
   for _, v := range dd {
   fmt.Fprintln(w, "<li class=\"text-center\"><a href=\"/getlaporan/"+v+"\">"+v+"</a></li>")
   }
}

   
func renderTemplate(w http.ResponseWriter, r *http.Request, p interface{}, tmpls ...string){
   tmp, _ := template.ParseFiles("templates/base.html")

   
   for _, v := range tmpls{
      tmp, _ = template.Must(tmp.Clone()).ParseFiles("templates/"+v+".html")
   }
 
   tmp.Execute(w, p)
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

   hariini := CreateTime()
   bul := hariini.Format("1")
   m, _ := strconv.Atoi(bul)
   y := hariini.Year()
   
   list := getListPasien(w, r, m, y)
   //sort.Reverse(sort.Interface(list))
   
   renderTemplate(w, r, list, "main", "listpts")
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
   renderTemplate(w, r, editData, "edit")
   
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
   renderTemplate(w, r, editData, "delete")
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

func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " dan atau dr. "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

func getListPasien(w http.ResponseWriter, r *http.Request, m,y int) []ListPasien{
   ctx := appengine.NewContext(r)
   
   u := user.Current(ctx)
   email := u.Email
   
   zone, err := time.LoadLocation("Asia/Makassar")
   if err != nil{
      fmt.Println("Err: ", err.Error())
   }
   in := time.Month(m)
   monIn := time.Date(y, in, 1, 0, 0, 0, 0, zone)
   monOut := monIn.AddDate(0, 1, 0)
   
   q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Order("JamDatang")

   t := q.Run(ctx)
   
   var daf KunjunganPasien
   var tar ListPasien
   var pts DataPasien
   var list []ListPasien
   list = append(list, ListPasien{})
   for {
      k, err := t.Next(&daf)
      if err == datastore.Done{break}
      if err != nil{
         fmt.Fprintln(w, "Error Fetching Data: ", err)
      }
      
      jam := ubahTanggal(daf.JamDatang, daf.ShiftJaga)
      if jam.Before(monIn) == true{continue}
      if jam.After(monOut) == true{continue}
      
      tar.TanggalFinal = jam.Format("2-01-2006")
      
      nocm := k.Parent()
      tar.NomorCM = nocm.StringID()

      err = datastore.Get(ctx, nocm, &pts)
      if err != nil {
            continue
			fmt.Fprintln(w, "Error Fetching Data Pasien: ", err)
         }
   
      tar.NamaPasien = properTitle(pts.NamaPasien)
	  tar.Diagnosis = properTitle(daf.Diagnosis)
	  
	  tar.LinkID = k.Encode()
      
      if daf.GolIKI == "1"{
         tar.IKI1 = "1"
         tar.IKI2 = ""
         }else{
         tar.IKI1 = ""
         tar.IKI2 = "1"
		 }
		 
	  list = append(list, tar)
   }
   return list
} 

func buatBCP(w http.ResponseWriter, r *http.Request){
   y, _ := strconv.Atoi(r.URL.Path[12:16])
   m, _ := strconv.Atoi(r.URL.Path[17:19])

   list := getListPasien(w, r, m, y)
   renderTemplate(w, r, list, "laporan")
   
} 