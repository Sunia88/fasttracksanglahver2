package ft

import (
   "net/http"
   "html/template"
   "appengine"
   "appengine/datastore"
   "fmt"
)

func GetCM(w http.ResponseWriter, r *http.Request, PasienAda bool){

   tmpl, err := template.New("adaPasien").ParseFiles("templates/adapasien.html")   //Parsing template untuk request ajax
   if err != nil {
      fmt.Fprint(w, "Error Parsing Template: ", err)
   }
   
   adakah := &PasienAda 
   ctx := appengine.NewContext(r)
   
   nocm := r.FormValue("nocm");
   
   parentKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   pasienKey := datastore.NewKey(ctx, "DataPasien", nocm, 0, parentKey)
   
   var pts DataPasien
   err = datastore.Get(ctx, pasienKey, &pts)
   if err != nil && err == datastore.ErrNoSuchEntity {
         tmpl.Execute(w, pts)
		    *adakah = false      
   }else{
   tmpl.Execute(w, pts)
   *adakah = true
   }  

}

func GetDatabyKey(item string, w http.ResponseWriter, r *http.Request) ListPasien {
   
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

func GetKursor(w http.ResponseWriter, ctx appengine.Context, tgl string) *datastore.Query{
   email, _, _ := AppCtx(ctx, "", "", "", "")
   _, _, kurK := AppCtx(ctx, "Dokter", email, "Kursor", tgl)

   kur := Kursor{}
   
   err := datastore.Get(ctx, kurK, &kur)
   if err != nil {
      fmt.Fprintln(w, "Error Fetching Database Kursor :", err)
      }
   kursor, err := datastore.DecodeCursor(kur.Point)
   if err != nil {
      fmt.Fprintln(w, "Error Decoding Cursor :", err)
      }
   q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Order("-JamDatang")
   q = q.Start(kursor)

   return q   
}

func GetListByCursor(w http.ResponseWriter, r *http.Request, m, y int)[]ListPasien{
   ctx := appengine.NewContext(r)
   monIn := DatebyInt(m,y)
   tgl := monIn.Format("2006/01")
   q := GetKursor(w, ctx, tgl)
   list := IterateList(ctx, w, q, monIn)
   return list
}

func GetListPasien(w http.ResponseWriter, r *http.Request, m,y int) []ListPasien{
   ctx := appengine.NewContext(r)
   email, _, _ := AppCtx(ctx, "", "", "", "")
   monIn := DatebyInt(m,y)
   q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Order("-JamDatang")
   list := IterateList(ctx, w, q, monIn)
   return list
}
