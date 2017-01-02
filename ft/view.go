package ft

import (
   "net/http"
   "html/template"
   "strings"
   "fmt"
)
func RenderPasien(w http.ResponseWriter, data interface{}, tmp string ){
      tmpl, err := template.New("tempPasien").Parse(tmp)
	  if err != nil {
	  fmt.Fprint(w, "Error Parsing: %v", err)
	  }
	  tmpl.Execute(w, data)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, p interface{}, tmpls ...string){
   tmp := template.New("base")
   tmp = tmp.Funcs(template.FuncMap{"strS": StringShift})   
   tmp, _ = template.ParseFiles("templates/base.html")
   for _, v := range tmpls{
      tmp, _ = template.Must(tmp.Clone()).ParseFiles("templates/"+v+".html")
   }
   tmp.Execute(w, p)
}

func ProperTitle(input string) string {
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

func StringShift(n string)string{
   var m string
   switch n {
     case "1" : m = "Pagi"
	 case "2" : m = "Sore"
	 case "3" : m = "Malam"
   }
   return m
}