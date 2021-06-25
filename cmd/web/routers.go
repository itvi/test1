package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//	mux := http.NewServeMux()

	// mux.HandleFunc("/", app.home)
	// //	mux.HandleFunc("/asset/categories", app.createAssetCategory())
	// mux.HandleFunc("/asset/categories", app.showAssetCategory)

	// static := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", static))

	// return standardMiddleware.Then(mux)

	mux := bone.New()
	static := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", static))

	mux.GetFunc("/", app.home)

	return standardMiddleware.Then(mux)
}
