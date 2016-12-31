package ft

import (
   "appengine"
   "appengine/datastore"
   "net/http"
   "fmt"
   "time"

)

type ListPasien struct {
   DataPasien
   KunjunganPasien
   TanggalFinal    string
   IKI1,IKI2 string
}

type SumIKI struct {
   Tanggal    string
   IKI1, IKI2 int
}

func IterateList(ctx appengine.Context, w http.ResponseWriter, q *datastore.Query, mon time.Time) []ListPasien{
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
      
      jam := UbahTanggal(daf.JamDatang, daf.ShiftJaga)
      if jam.Before(mon) == true{break}
      
      tar.TanggalFinal = jam.Format("2-01-2006")
      
      nocm := k.Parent()
      tar.NomorCM = nocm.StringID()

      err = datastore.Get(ctx, nocm, &pts)
      if err != nil {
            continue
			fmt.Fprintln(w, "Error Fetching Data Pasien: ", err)
         }
   
      tar.NamaPasien = ProperTitle(pts.NamaPasien)
	  tar.Diagnosis = ProperTitle(daf.Diagnosis)
	  
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

func ListLaporan(w http.ResponseWriter, r *http.Request)[]string{
   ctx := appengine.NewContext(r)
   email, _, _ := AppCtx(ctx, "", "", "", "")
   _, key, _ := AppCtx(ctx, "Dokter", email, "Kursor", "")
   kur := []Kursor{}
   q := datastore.NewQuery("Kursor").Ancestor(key)
   keys, err := q.GetAll(ctx, &kur)
   if err != nil{
      fmt.Fprintln(w, "Error Fetching Kursor :", err)
   }

   var list []string
   for _, v := range keys {
       m := v.StringID()
	   list = append(list, m)
   }
   for i, j := 0, len(list)-1 ; i < j ; i, j = i+1, j-1{
      list[i], list[j] = list[j], list[i]
   }   
   return list
}

func ListIKI(w http.ResponseWriter, r *http.Request, m, y int)[]SumIKI{
   list := GetListPasien(w, r, m, y)
   
   bl := UbahBulanIni(0).Day()

	var ikiBulan []SumIKI
	ikiBulan = append(ikiBulan, SumIKI{})
	for h:= bl; h > 0; h--{
	   dataIKI := SumIKI{}	
	   q := UbahBulanIni(h).Format("2-01-2006")
	   var u1, u2 int
	   for _, v :=  range list{
	      if v.TanggalFinal != q {continue}
		  if v.IKI1 == "1" {
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
    
   for i, j := 0, len(ikiBulan)-1 ; i < j ; i, j = i+1, j-1{
      ikiBulan[i], ikiBulan[j] = ikiBulan[j], ikiBulan[i]
   }
	return ikiBulan
}
