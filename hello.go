package main

import (
	"fmt"
	"ft"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func init() {

	http.HandleFunc("/", index)
	http.HandleFunc("/mainpage", mainPage)
	http.HandleFunc("/getcm", getCM)

	//http.HandleFunc("/getinfo", getInfo)
	http.HandleFunc("/inputdatapts", inPts)

	//http.HandleFunc("/getiki", listIKI)
	//http.HandleFunc("/testdb", testdb)

	http.HandleFunc("/entri/edit/", ft.EditEntri)
	http.HandleFunc("/entri/update", ft.UpdateEntri)
	http.HandleFunc("/entri/del/", ft.DeleteEntri)
	http.HandleFunc("/entri/delete", ft.ConfirmDeleteEntri)
	http.HandleFunc("/entri/editdate/", ft.EditDate)
	http.HandleFunc("/entri/updatetanggal", ft.UpdateTanggal)

	//http.HandleFunc("/getlaporan", listLaporan)
	http.HandleFunc("/getlaporan/", buatBCP)
	http.HandleFunc("/admin", adminPage)
	http.HandleFunc("/admin/addstaff", ft.AddStaff)
	http.HandleFunc("/admin/delete/", ft.DeletePage)
	http.HandleFunc("/admin/confdel/", ft.ConfDel)
	http.HandleFunc("/test", test)

}

var PasienAda bool = false

func inPts(w http.ResponseWriter, r *http.Request) {
	ft.InputPasien(w, r, PasienAda)
}

func getCM(w http.ResponseWriter, r *http.Request) {
	ft.GetCM(w, r, PasienAda)
}

type DataPasien struct {
	NamaPasien              string
	NomorCM, JenKel, Alamat string
	TglDaftar, Umur         time.Time
}

type KunjunganPasien struct {
	Diagnosis, LinkID      string
	GolIKI, ATS, ShiftJaga string
	JamDatang              time.Time
	Dokter                 string
	Hide                   bool
}

type ListPasien struct {
	DataPasien
	KunjunganPasien
	TanggalFinal string
	IKI1, IKI2   string
}

type Kursor struct {
	Point string
}

type sumIKI struct {
	Tanggal    string
	IKI1, IKI2 int
}

type WebObject struct {
	IKI    []ft.SumIKI
	List   []ft.ListPasien
	Kur    []string
	Email  string
	Logout string
}

type Staff struct {
	Email       string
	NamaLengkap string
	LinkID      string
}
type Web struct {
	Welcome, Link, LinkStr string
	Staff                  []ft.Staff
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	login, _ := user.LoginURL(ctx, "/")
	logout, _ := user.LogoutURL(ctx, "/")
	var web Web
	if u.Admin {
		web.Welcome = "Welcome Admin-san"
		web.Link = logout
		web.LinkStr = "Logout"
		tmp := template.Must(template.New("adminOk.html").ParseFiles("templates/adminOk.html"))
		n := ft.GetStaff(ctx, u.Email)
		web.Staff = n
		err := tmp.Execute(w, web)
		if err != nil {
			fmt.Fprintln(w, "Error Parsing Template :", err)
		}
	} else {
		web.Welcome = "Admin login only"
		web.Link = login
		web.LinkStr = "Login"
		tmp := template.Must(template.New("adminNil.html").ParseFiles("templates/adminNil.html"))
		err := tmp.Execute(w, web)
		if err != nil {
			fmt.Fprintln(w, "Error Parsing Template :", err)
		}
	}
}
func buatBCP(w http.ResponseWriter, r *http.Request) {
	y, _ := strconv.Atoi(r.URL.Path[12:16])
	m, _ := strconv.Atoi(r.URL.Path[17:19])

	var web WebObject
	web.Kur = ft.ListLaporan(w, r)
	x := ft.GetListByCursor(w, r, m, y)
	web.IKI = ft.ListIKI(w, r, m, y, x)
	var k []ft.ListPasien
	k = append(k, ft.ListPasien{})
	k = append(k, x...)
	web.List = k
	ft.RenderTemplate(w, r, web, "laporan")

}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	const maaf1 string = `
<html>
<head>
<title>Welcome</title>
</head>
<body>
<p>Maaf anda tidak dapat mengakses aplikasi. Silahkan hubungi admin</p><br>
<a href=`

	const maaf2 string = `
>Logout</a>
</body>
</html>
   `
	ctx := appengine.NewContext(r)
	//   var login string
	if u := user.Current(ctx); !u.Admin {
		email, _, _ := ft.AppCtx(ctx, "", "", "", "")
		_, key, _ := ft.AppCtx(ctx, "Staff", email, "", "")
		logout, _ := user.LogoutURL(ctx, "/")
		//     login, _ := user.LoginURL(ctx, "/")

		var staff Staff
		err := datastore.Get(ctx, key, &staff)
		if err != nil {
			fmt.Fprintln(w, maaf1+logout+maaf2)
		} else {
			http.Redirect(w, r, "/mainpage", http.StatusSeeOther)
		}
	}
	if u := user.Current(ctx); u.Admin {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	hariini := ft.CreateTime()
	bul := hariini.Format("1")
	m, _ := strconv.Atoi(bul)
	y := hariini.Year()

	if hariini.Day() == 1 && hariini.Hour() > 8 {
		ft.CreateKursor(w, ctx)
	}
	email, _, _ := ft.AppCtx(ctx, "", "", "", "")

	web := WebObject{}
	web.List = ft.GetListPasien(w, r, m, y)
	web.IKI = ft.ListIKI(w, r, m, y, web.List)
	web.Kur = ft.ListLaporan(w, r)
	web.Email = email
	logout, _ := user.LogoutURL(ctx, "/")
	web.Logout = logout
	ft.RenderTemplate(w, r, web, "main")
}

func test(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	hariini := ft.CreateTime()
	bul := hariini.Format("1")
	m, _ := strconv.Atoi(bul)
	y := hariini.Year()

	if hariini.Day() == 1 && hariini.Hour() > 8 {
		ft.CreateKursor(w, ctx)
	}
	email, _, _ := ft.AppCtx(ctx, "", "", "", "")

	web := WebObject{}
	web.List = ft.GetListPasien(w, r, m, y)
	web.IKI = ft.ListIKI(w, r, m, y, web.List)
	web.Kur = ft.ListLaporan(w, r)
	web.Email = email
	logout, _ := user.LogoutURL(ctx, "/")
	web.Logout = logout
	days := time.Date(2016, time.February, 1, 0, 0, 0, 0, time.UTC)
	jml := days.AddDate(0, 1, -1).Format("2")
	fmt.Fprintln(w, jml)
	fmt.Fprintln(w, days)
	for _, v := range web.List {
		fmt.Fprintln(w, v.TanggalFinal)
	}

}
