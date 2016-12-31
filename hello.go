package main

import (
    "net/http"
	"time"
	"strconv"
	"appengine"
	"appengine/datastore"
	"appengine/user"
    "fmt"
	"ft"
)


func init() {
    
    http.HandleFunc("/", index)
	http.HandleFunc("/mainpage", mainPage)
	http.HandleFunc("/getcm", getCM)

	//http.HandleFunc("/getinfo", getInfo)
	http.HandleFunc("/inputdatapts", inPts)

	
	//http.HandleFunc("/getiki", listIKI)
	http.HandleFunc("/testdb", testdb)
	
	http.HandleFunc("/entri/edit/", ft.EditEntri)
	http.HandleFunc("/entri/update", ft.UpdateEntri)
	http.HandleFunc("/entri/del/", ft.DeleteEntri)
	http.HandleFunc("/entri/delete", ft.ConfirmDeleteEntri)
	
	//http.HandleFunc("/getlaporan", listLaporan)
	http.HandleFunc("/getlaporan/", buatBCP)
	http.HandleFunc("/admin/", adminPage)

}

var PasienAda bool = false

func inPts(w http.ResponseWriter, r *http.Request){
   ft.InputPasien(w, r, PasienAda)
}

func getCM(w http.ResponseWriter, r *http.Request){
   ft.GetCM(w, r, PasienAda)
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

type ListPasien struct {
   DataPasien
   KunjunganPasien
   TanggalFinal    string
   IKI1,IKI2 string
}

type Kursor struct {
   Point    string
}

type sumIKI struct {
   Tanggal    string
   IKI1, IKI2 int
}

type WebObject struct {
   IKI     []ft.SumIKI
   List    []ft.ListPasien
   Kur     []string
   Email   string
   Logout  string
}

type Staff struct {
   Email       string
   NamaLengkap string
}

func adminPage(w http.ResponseWriter, r *http.Request){
   ctx := appengine.NewContext(r)
   //email, key, _ := appCtx(ctx,"Staff","staff","","")
   
   if u := user.Current(ctx); !u.Admin {
      fmt.Fprintln(w, "Admin login only", http.StatusUnauthorized)
	  time.Sleep(2000 * time.Millisecond)
	  http.Redirect(w, r, "/mainpage", http.StatusSeeOther)
	   
   }
      
   doc := []Staff{}
   q := datastore.NewQuery("Staff")
   _, err := q.GetAll(ctx, &doc)
   if err != nil && err != datastore.ErrNoSuchEntity{
      fmt.Fprintln(w, "Error Fetching Data Staff :", err)
   }
   ft.RenderTemplate(w, r, doc, "admin")
}
func buatBCP(w http.ResponseWriter, r *http.Request){
   y, _ := strconv.Atoi(r.URL.Path[12:16])
   m, _ := strconv.Atoi(r.URL.Path[17:19])

   var web WebObject
   web.Kur = ft.ListLaporan(w,r)
   web.List = ft.GetListByCursor(w, r, m, y)
   ft.RenderTemplate(w, r, web, "laporan")
   
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
   
   if u := user.Current(ctx); !u.Admin {
      http.Redirect(w, r, "/mainpage", http.StatusSeeOther)
	   
   } 
   if u := user.Current(ctx); u.Admin {
      http.Redirect(w, r, "/admin/", http.StatusSeeOther)
   }
   login, _ = user.LoginURL(ctx, "/mainpage")
   http.Redirect(w, r, login, http.StatusSeeOther)
}

func mainPage(w http.ResponseWriter, r *http.Request){

   ctx := appengine.NewContext(r)
   hariini := ft.CreateTime()
   bul := hariini.Format("1")
   m, _ := strconv.Atoi(bul)
   y := hariini.Year()
   ft.CreateKursor(w,ctx)
   email, _, _ := ft.AppCtx(ctx, "", "", "", "")
   
   web := WebObject{}
   web.IKI = ft.ListIKI(w, r, m, y)
   web.List = ft.GetListPasien(w, r, m, y)
   web.Kur = ft.ListLaporan(w,r)
   web.Email = email 
   logout, _ := user.LogoutURL(ctx, "/")
   web.Logout = logout
   ft.RenderTemplate(w, r, web, "main")
}

//Fungsi Misc

func testdb(w http.ResponseWriter, r *http.Request){
   ctx := appengine.NewContext(r)
   hariini := ft.CreateTime()
   bul := hariini.Format("1")
   m, _ := strconv.Atoi(bul)
   y := hariini.Year()
   ft.CreateKursor(w,ctx)
   //email, _, _ := ft.AppCtx(ctx, "", "", "", "")
   
   j := listing(w, r, m, y)

   fmt.Fprintln(w, j)
}

func listing(w http.ResponseWriter, r *http.Request, m, y int) []sumIKI{
   list := ft.GetListPasien(w, r, m, y)
   
   bl := ft.UbahBulanIni(0).Day()

	var ikiBulan []sumIKI

	for h:= bl; h > 0; h--{
	   dataIKI := sumIKI{}	
	   q := ft.UbahBulanIni(h).Format("2-01-2006")
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
	//listIKI := haha(ikiBulan)
	return ikiBulan
}
