package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	//	mux.HandleFunc("/asset/categories", app.createAssetCategory())
	mux.HandleFunc("/asset/categories", app.showAssetCategory)

	static := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", static))

	return mux
}
