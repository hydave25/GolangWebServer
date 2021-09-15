// https://tutorialedge.net/golang/creating-simple-web-server-with-golang/

package main

import (
	//	"fmt"
	"fmt"
	"strconv"
	"text/template"

	//	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	//	"github.com/go-delve/delve/pkg/proc/test"
	//	"go.opencensus.io/stats/view"
	//	"golang.org/x/net/http/httpproxy"
)

var counter int
var mutex = &sync.Mutex{}

const defaultAddr = ":8080"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

type templateData struct {
	Message string
}

var (
	data templateData
	tmpl *template.Template
)

// home responds to requests by rendering an HTML page.
func home(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello from Cloud Code! Received request: %s %s", r.Method, r.URL.Path)
	if err := tmpl.Execute(w, data); err != nil {
		msg := http.StatusText(http.StatusInternalServerError)
		log.Printf("template.Execute: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprint(w, strconv.Itoa(counter))
	mutex.Unlock()
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {

	http.HandleFunc("/", echoString)
	//	http.HandleFunc("/", home)
	//	http,HandleFunc("/", func viewHandler())
	//http.HandleFunc("/", http,http.ResponseWriter, *http.NewRequestWithContext()
	http.HandleFunc("/increment", incrementCounter)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
