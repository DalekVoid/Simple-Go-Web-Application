package main

import (
  "io/ioutil"
  "net/http"
  "html/template"
  "regexp"
  "flag"
  "log"
  "net"
)

var (
  addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
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
//global variable to cache template parser
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// global variable to cache regex parser
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// helper functions
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

// wrapper function
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request){
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2]) // m[2] is the title
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err  := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
  flag.Parse()
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))

  if *addr {
      l, err := net.Listen("tcp", "127.0.0.1:0")
      if err != nil {
          log.Fatal(err)
      }
      err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
      if err != nil {
          log.Fatal(err)
      }
      s := &http.Server{}
      s.Serve(l)
      return
  }

  http.ListenAndServe(":8080", nil)
}
