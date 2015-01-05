package main

import (
  "io/ioutil"
  "net/http"
  "html/template"
)

type Page struct {
  Title string
  Body []byte
}

// Page related
func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600) // 0600 for read-write permission.
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

// http handlers
func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):] // len to slice the url and slice leading "/view/" in the request
  p, _  := loadPage(title)
  t, _ := template.ParseFiles("view.html")
  t.Execute(w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/edit/"):]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  t, _ := template.ParseFiles("edit.html")
  t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

}
func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", nil)
}
