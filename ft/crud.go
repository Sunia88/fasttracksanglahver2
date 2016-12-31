package ft

import (
   "net/http"
   "time"
   "appengine"
   "appengine/datastore"
   "fmt"
)

type Kursor struct {
   Point    string
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
type WebObject struct {
   IKI     []SumIKI
   List    []ListPasien
   Kur     []string
   Email   string
   Logout  string
}


func CreateKursor(w http.ResponseWriter, ctx appengine.Context){

   email, _, _ := AppCtx(ctx, "", "", "", "")
   wkt := CreateTime()
   bul := wkt.Month()
   th := wkt.Year()

   tgl := wkt.Format("2006/01")
   
   _, _, kurKey := AppCtx(ctx, "Dokter", email, "Kursor", tgl)
   q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Order("-JamDatang")
   
   kur := Kursor{}
   listkur := []Kursor{}
   
   var kun KunjunganPasien

   err := datastore.Get(ctx, kurKey, &kur)
   if err != nil && err != datastore.ErrNoSuchEntity{
      fmt.Fprintln(w, "Error Fetching Data Kursor: ", err)
   }
   if err != nil && err == datastore.ErrNoSuchEntity {
      t := q.Run(ctx)

      mon := time.Date(th, bul, 1, 0, 0, 0, 0, time.UTC)
      for {
         _, err := t.Next(&kun)
	     if err == datastore.Done{break}
	     if err != nil{
	        fmt.Fprintln(w, "Error Fetching Data: ", err)
	     }
	     
		 jamEdit := UbahTanggal(kun.JamDatang, kun.ShiftJaga)

	     if jamEdit.After(mon) != true {
	        cursor, _ := t.Cursor()
		    kur.Point = cursor.String()
			listkur = append(listkur, kur)
			mon = mon.AddDate(0, -1, 0)
			bln := mon.Format("2006/01")
			_, _, keyKur := AppCtx(ctx, "Dokter", email, "Kursor", bln)
			if _, err := datastore.Put(ctx, keyKur, &kur); err != nil {
			   fmt.Fprint(w, "Error Writing to Database: ", err)
			}
	     }
      }
   }   
}

func ConfirmDeleteEntri(w http.ResponseWriter, r *http.Request){
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

   err = datastore.Delete(ctx, keyKun)

   http.Redirect(w, r, "/mainpage", http.StatusSeeOther)   
}

func DeleteEntri(w http.ResponseWriter, r *http.Request){
   keyitem := r.URL.Path[11:]
   web := WebObject{}
   web.List = append(web.List, GetDatabyKey(keyitem, w, r))
   web.Kur = ListLaporan(w,r)
   RenderTemplate(w, r, web, "delete")
}

func EditEntri(w http.ResponseWriter, r *http.Request){
   keyitem := r.URL.Path[12:]
   web := WebObject{}
   web.List = append(web.List, GetDatabyKey(keyitem, w, r))
   web.Kur = ListLaporan(w,r)
   RenderTemplate(w, r, web, "edit")  
}

func InputPasien(w http.ResponseWriter, r *http.Request, PasienAda bool){
   if r.Method != "POST" {
      http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	  return
   }

   ctx := appengine.NewContext(r)
   nocm := r.FormValue("nocm")
   doc, _, _ := AppCtx(ctx, "", "", "", "")   
   _, parentKey, pasienKey := AppCtx(ctx, "DataPasien", nocm, "KunjunganPasien", "")
   
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
   http.Redirect(w, r, "/mainpage", http.StatusSeeOther)
 
}


func UpdateEntri(w http.ResponseWriter, r *http.Request){
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