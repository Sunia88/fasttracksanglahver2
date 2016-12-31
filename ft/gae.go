package ft

import (
   "appengine"
   "appengine/datastore"
   "appengine/user"
   "time"
   "fmt"   
)
func AppCtx(ctx appengine.Context, kind1 string, id1 string, kind2 string, id2 string) (string, *datastore.Key, *datastore.Key){

   u := user.Current(ctx)
   email := u.Email
   
   gpKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
   parKey := datastore.NewKey(ctx, kind1, id1, 0, gpKey)
   chldKey := datastore.NewKey(ctx, kind2, id2, 0, parKey)

   return email, parKey, chldKey
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

func UbahBulanIni(d int) time.Time{
   y, m, _ := time.Now().Date()
   zone, err := time.LoadLocation("Asia/Makassar")
   if err != nil{
      fmt.Println("Err: ", err.Error())
   }
   bulan := time.Date(y, m, d, 0, 0, 0, 0, zone)
   return bulan

}


func UbahTanggal(tgl time.Time, shift string) time.Time{
   
   ubah := tgl   
   jam := ubah.Hour()

   if jam < 12 && shift == "3"{
	     ubah = tgl.AddDate(0,0,-1)
         }
   return ubah
}

func DatebyInt(m, y int) time.Time{
   in := time.Month(m)
   monIn := time.Date(y, in, 1, 0, 0, 0, 0, time.UTC)
   
   return monIn
}