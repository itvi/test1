package main

import (
	"ams/pkg/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	categories, err := app.assetCategory.GetCategories()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, category := range categories {
		fmt.Fprintf(w, "%v\n", category)
	}
}
func (app *application) createAssetCategory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			rw.Write([]byte("GET"))
		}

		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
				return
			}

			category := &models.AssetCategory{Code: "H001", Name: "Hardware"}

			err := app.assetCategory.Create(category)
			if err != nil {
				app.serverError(rw, err)
				return
			}

			http.Redirect(rw, r, "/asset/categories", http.StatusSeeOther)
		}
	}
}

func (app *application) showAssetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	category, err := app.assetCategory.GetCategoryByID(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Fprintf(w, "%v", category)
}
