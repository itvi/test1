package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	m0 := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	m1 := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := bone.New()

	static := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", static))

	mux.Get("/", m1.ThenFunc(app.home))

	// asset
	mux.Get("/assets/dropdown", m1.ThenFunc(app.assetDropdown)) // asset number dropdown list
	mux.Get("/assets/list", m1.ThenFunc(app.getAssets))
	mux.Get("/assets", m1.ThenFunc(app.indexAsset))
	mux.Get("/asset", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.addAsset))
	mux.Post("/assets", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.addAsset))
	mux.Get("/assets/:id", m1.ThenFunc(app.editAsset))
	mux.Post("/assets/:id", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.editAsset))
	mux.Delete("/assets/:id", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteAsset))
	mux.Get("/upload", m1.ThenFunc(app.Upload))
	mux.Post("/upload", m1.ThenFunc(app.Upload))

	// asset category
	mux.Get("/asset/categories/dropdown", http.HandlerFunc(app.assetCategoryDropdown))
	mux.Get("/asset/categories", m1.ThenFunc(app.assetCategories))
	mux.Get("/asset/category", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.createAssetCategory()))
	mux.Post("/asset/categories", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.createAssetCategory()))
	mux.Get("/asset/categories/:id", m1.ThenFunc(app.editAssetCategory))
	mux.Post("/asset/categories/:id", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.editAssetCategory))
	mux.Delete("/asset/categories/:id", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteAssetCategory))

	// asset status
	mux.Get("/asset/statuses", m1.ThenFunc(app.indexAssetStatus))
	mux.Get("/asset/status", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.addAssetStatus))
	mux.Post("/asset/statuses", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.addAssetStatus))
	mux.Get("/asset/statuses/:id", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.editAssetStatus))
	mux.Post("/asset/statuses/:id", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.editAssetStatus))
	mux.Delete("/asset/statuses/:id", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteAssetStatus))

	// asset management
	mux.Get("/asset/management/mov", m1.ThenFunc(app.addAssetMovAndConfig))
	mux.Post("/asset/management/mov", m1.ThenFunc(app.addAssetMovAndConfig))

	// computer config
	mux.Get("/computer/config/init", m1.ThenFunc(app.InitializeComputerConfig))
	mux.Post("/computer/config/init", m1.ThenFunc(app.InitializeComputerConfig))
	mux.Get("/computer/config/search", m1.ThenFunc(searchByIP))
	mux.Get("/searchByIPs", m1.ThenFunc(searchByIPs))
	mux.Get("/inventory", m1.ThenFunc(app.inventory))

	// user
	mux.Get("/user/signup", m1.ThenFunc(app.signup))
	mux.Post("/user/signup", m1.ThenFunc(app.signup))
	mux.Get("/user/login", m1.ThenFunc(app.login))
	mux.Post("/user/login", m1.ThenFunc(app.login))
	mux.Post("/user/logout", m1.Append(app.requireAuthenticatedUser).ThenFunc(app.logout))

	// upload
	mux.Get("/upload", m1.ThenFunc(app.Upload))
	mux.Post("/upload", m1.ThenFunc(app.Upload))

	return m0.Then(mux)
}
