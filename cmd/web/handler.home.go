package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "index.html", nil)
}

func (app *application) inventory(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "inventory.html", nil)
}
