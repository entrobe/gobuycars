package main

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("new.html", "view.html"))

type Item struct {
	Title       string
	Description string
	ImageURL    string
}

func (i *Item) String() string {
	return i.Title
}

func loadItem(key string) (*Item, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}

	reply, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	if len(reply) == 0 {
		return nil, errors.New("No Item Found")
	}

	item := &Item{}
	err = redis.ScanStruct(reply, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, i *Item) {
	err := templates.ExecuteTemplate(w, tmpl+".html", i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Path[len("/view/"):]
	i, err := loadItem(item)
	if err != nil {
		log.Print(err)
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "view", i)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new", nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	desc := r.FormValue("description")
	title := r.FormValue("title")
	imageURL := r.FormValue("imageurl")
	i := &Item{Title: title, Description: desc, ImageURL: imageURL}
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Do("HMSET", redis.Args{}.Add(i.Title).AddFlat(i)...)

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
