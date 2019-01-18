package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Category struct {
	Name    string
	Mapping map[string]string
}

type Meetings struct {
	Team       string
	Categories []Category
}

type IndexHandler struct {
	meetings Meetings
}

func NewIndexHandler(meetings Meetings) *IndexHandler {
	return &IndexHandler{meetings}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	tmpl := template.Must(template.ParseFiles("handlers/index.html"))
	tmpl.Execute(w, h.meetings)
}

type RedirectHandler struct {
	mapping map[string]string
}

func NewRedirectHandler(meetings Meetings) *RedirectHandler {
	mapping := map[string]string{}
	for _, category := range meetings.Categories {
		for name, id := range category.Mapping {
			mapping[name] = id
		}
	}

	return &RedirectHandler{mapping}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["meeting"]

	if id, ok := h.mapping[name]; ok {
		zoomURL := fmt.Sprintf("https://pivotal.zoom.us/j/%s", id)
		http.Redirect(w, r, zoomURL, http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Meeting with name \"%s\" not found.\n", name)
}
