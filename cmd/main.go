package main

import (
	Asciiart "ascii-art-web/cmd/asciiart"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var tmpl, errTmpl *template.Template

type PageData struct {
	Title string
	Item  string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		data := PageData{
			Title: "Error 404",
			Item:  "Page not found",
		}
		errTmpl.Execute(w, data)
		return
	}
	out, err := os.ReadFile("cmd/output.txt")
	if err != nil {
		os.Create("cmd/output.txt")
	}

	data := PageData{
		Title: "Ascii Art Web",
		Item:  string(out),
	}
	tmpl.Execute(w, data)
}

func postArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		w.WriteHeader(405)
		data := PageData{
			Title: "Error 405",
			Item:  "Method not Allowed",
		}
		errTmpl.Execute(w, data)
		return
	}

	r.ParseForm()
	input, ok := r.Form["input"]
	if !ok {
		log.Print("not ok")
	}
	banner, ok := r.Form["banner"]
	if !ok {
		log.Print("not ok")
	}

	if Asciiart.NonAsciiCheck(input) {
		w.WriteHeader(400)
		data := PageData{
			Title: "Error 400. Bad request.",
			Item:  "Exeptional characters. Really.",
		}
		errTmpl.Execute(w, data)
		return
	}

	if input[0] == "" {
		data := PageData{
			Title: "Ascii Art Web",
			Item:  "",
		}
		tmpl.Execute(w, data)
		return
	}

	if !Asciiart.Asciiart(input[0], banner[0]) {
		w.WriteHeader(400)
		data := PageData{
			Title: "Error 400. Bad request.",
			Item:  "",
		}
		errTmpl.Execute(w, data)
		return
	}

	out, err := os.ReadFile("cmd/output.txt")
	if err != nil {
		os.Create("cmd/output.txt")
	}

	data := PageData{
		Title: "Ascii Art Web",
		Item:  string(out),
	}
	if r.Form["action"][0] == "download" {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(string(out))))
		w.Header().Set("Content-Disposition", "name='cmd/output.txt'; filename='code.txt'")
		w.Write(out)
		return
	} else if r.Form["action"][0] == "submit" {
		os.Remove("cmd/output.txt")
		tmpl.Execute(w, data)
		return
	}

	w.WriteHeader(500)
	data = PageData{
		Title: "Error 500.",
		Item:  "Internal Server Error",
	}
	errTmpl.Execute(w, data)
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))

	tmpl = template.Must(template.ParseFiles("templates/a.htm"))
	errTmpl = template.Must(template.ParseFiles("templates/error.htm"))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/ascii-art", postArt)

	// fs := http.FileServer(http.Dir("./templates"))
	// mux.Handle("/templates/", http.StripPrefix("/templates/", fs))

	fmt.Println("Listening on http://localhost:8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
