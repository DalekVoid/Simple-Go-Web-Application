package main

import (
  "fmt"
  "io/ioutil"
)

type Page struct {
  Title string
  Body []byte
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600) // 0600 for read-write permission.
}

func loadPage(title string) *Page {
  filename := title + ".txt"
  body, _  := ioutil.ReadFile(filename) // wildcard character(is it how it is called in go?(oh it is called the blank identifier)) throwaway the unwanted error value.
  return &Page{Title: title, Body: body}
}


