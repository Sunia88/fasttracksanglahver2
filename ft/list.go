package ft

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
)

type ListPasien struct {
	DataPasien
	KunjunganPasien
	TanggalFinal string
	IKI1, IKI2   string
}

type SumIKI struct {
	Tanggal    string
	IKI1, IKI2 int
}

func IterateList(ctx appengine.Context, w http.ResponseWriter, q *datastore.Query, mon time.Time) []ListPasien {
	t := q.Run(ctx)
	monAf := mon.AddDate(0, 1, 0)
	var daf KunjunganPasien
	var tar ListPasien
	var pts DataPasien
	var list []ListPasien
	for {
		k, err := t.Next(&daf)
		if err == datastore.Done {
			break
		}
		if err != nil {
			fmt.Fprintln(w, "Error Fetching Data: ", err)
		}
		daf.JamDatang = daf.JamDatang.Add(time.Duration(8) * time.Hour)
		jam := UbahTanggal(daf.JamDatang, daf.ShiftJaga)
		if jam.After(monAf) == true {
			continue
		}
		if jam.Before(mon) == true {
			break
		}
		if daf.Hide == true {
			continue
		}
		tar.TanggalFinal = jam.Format("02-01-2006")

		nocm := k.Parent()
		tar.NomorCM = nocm.StringID()

		err = datastore.Get(ctx, nocm, &pts)
		if err != nil {
			fmt.Fprintln(w, "Error Fetching Data Pasien: ", err)
		}

		tar.NamaPasien = ProperTitle(pts.NamaPasien)
		tar.Diagnosis = ProperTitle(daf.Diagnosis)
		tar.ShiftJaga = daf.ShiftJaga
		tar.LinkID = k.Encode()

		if daf.GolIKI == "1" {
			tar.IKI1 = "1"
			tar.IKI2 = ""
		} else {
			tar.IKI1 = ""
			tar.IKI2 = "1"
		}

		list = append(list, tar)
	}

	return list
}

func ListLaporan(w http.ResponseWriter, r *http.Request) []string {
	ctx := appengine.NewContext(r)
	email, _, _ := AppCtx(ctx, "", "", "", "")
	_, key, _ := AppCtx(ctx, "Dokter", email, "Kursor", "")
	kur := []Kursor{}
	q := datastore.NewQuery("Kursor").Ancestor(key)
	keys, err := q.GetAll(ctx, &kur)
	if err != nil {
		fmt.Fprintln(w, "Error Fetching Kursor :", err)
	}

	var list []string
	for _, v := range keys {
		m := v.StringID()
		list = append(list, m)
	}
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

func ListIKI(w http.ResponseWriter, r *http.Request, m, y int, n []ListPasien) []SumIKI {
	for i, j := 0, len(n)-1; i < j; i, j = i+1, j-1 {
		n[i], n[j] = n[j], n[i]
	}
	mo := DatebyInt(m, y)
	wkt := time.Date(mo.Year(), mo.Month(), 1, 0, 0, 0, 0, time.UTC)
	strbl := wkt.AddDate(0, 1, -1).Format("2")
	bl, _ := strconv.Atoi(strbl)
	var ikiBulan []SumIKI
	ikiBulan = append(ikiBulan, SumIKI{})
	for h := 1; h <= bl; h++ {
		dataIKI := SumIKI{}
		q := time.Date(mo.Year(), mo.Month(), h, 0, 0, 0, 0, time.UTC).Format("02-01-2006")
		var u1, u2 int
		for _, v := range n {
			if v.TanggalFinal != q {
				continue
			}
			if v.IKI1 == "1" {
				u1++
			} else {
				u2++
			}
		}

		if u1 == 0 && u2 == 0 {
			continue
		}
		dataIKI.Tanggal = q
		dataIKI.IKI1 = u1
		dataIKI.IKI2 = u2

		ikiBulan = append(ikiBulan, dataIKI)
	}

	return ikiBulan
}
