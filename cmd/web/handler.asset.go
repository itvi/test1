package main

import (
	"ams/pkg/models"
	"ams/pkg/models/forms"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-zoo/bone"
)

func (app *application) indexAsset(w http.ResponseWriter, r *http.Request) {
	statuses, err := app.assets.GetAssets()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "asset.index.html", &templateData{
		Assets: statuses,
	})
}

func (app *application) addAsset(w http.ResponseWriter, r *http.Request) {
	page := "asset.add.html"
	if r.Method == "GET" {
		app.render(w, r, page, &templateData{
			Form: forms.New(nil),
		})
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)

		form.Required("asset_number", "category_code")
		form.MaxLength("asset_number", 9)

		if !form.Valid() {
			app.render(w, r, page, &templateData{
				Form: form,
			})
			return
		}

		number := form.Get("asset_number")
		categoryCode := form.Get("category_code")
		supplier := form.Get("supplier")
		mdl := form.Get("model")
		sn := form.Get("sn")
		warranty := form.Get("warranty")
		remark := form.Get("remark")

		warrantyDate, err := time.Parse("2006-01-02", warranty)
		if err != nil {
			app.serverError(w, err)
			return
		}

		category, err := app.assetCategory.GetCategoryByCode(categoryCode)
		if err != nil {
			app.serverError(w, err)
			return
		}
		asset := &models.Asset{
			Number:   number,
			Category: *category,
			Supplier: supplier,
			Model:    mdl,
			SN:       sn,
			Warranty: warrantyDate,
			Remark:   remark,
		}

		err = app.assets.Add(asset)
		if err == models.ErrDuplicate {
			form.Errors.Add("asset_number", "编号已存在！")
			form.Add("category_name", category.Name)
			app.render(w, r, page, &templateData{Form: form})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flash", "添加成功！")
		http.Redirect(w, r, "/assets", http.StatusSeeOther)
	}
}

func (app *application) editAsset(w http.ResponseWriter, r *http.Request) {
	page := "asset.edit.html"
	id, err := strconv.Atoi(bone.GetValue(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	if r.Method == "GET" {
		asset, err := app.assets.GetAssetByID(id)
		if err == models.ErrNoRecord {
			app.notFound(w)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		form := forms.New(nil)
		app.render(w, r, page, &templateData{
			Form:  form,
			Asset: asset,
		})
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("asset_number", "category_code")

		number := form.Get("asset_number")
		categoryCode := form.Get("category_code")
		supplier := form.Get("supplier")
		mdl := form.Get("model")
		sn := form.Get("sn")
		warranty := form.Get("warranty")
		remark := form.Get("remark")

		warrantyDate, err := time.Parse("2006-01-02", warranty)
		if err != nil {
			app.serverError(w, err)
			return
		}

		asset := &models.Asset{
			ID:       id,
			Number:   number,
			Category: models.AssetCategory{Code: categoryCode},
			Supplier: supplier,
			Model:    mdl,
			SN:       sn,
			Warranty: warrantyDate,
			Remark:   remark,
		}

		if !form.Valid() {
			app.render(w, r, page, &templateData{
				Form:  form,
				Asset: asset,
			})
			return
		}

		if err = app.assets.Edit(asset); err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flash", "更新成功!")
		http.Redirect(w, r, "/assets", http.StatusSeeOther)
	}
}

func (app *application) deleteAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(bone.GetValue(r, "id"))
	if err != nil {
		app.notFound(w)
		return
	}
	if err = app.assets.Del(id); err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "删除成功!")
}

// getDevices return json for client table
func (app *application) getAssets(w http.ResponseWriter, r *http.Request) {
	data, err := app.assets.GetAssets()
	if err != nil {
		log.Println(err)
		return
	}

	j, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (app *application) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		app.render(w, r, "asset.upload.html", nil)
	}
	if r.Method == "POST" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			fmt.Println(err)
			return
		}
		files := r.MultipartForm.File["files"]
		var msg string
		if err := app.assets.Upload(files); err != nil {
			fmt.Println(err)
			msg = err.Error()
		}

		j, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}

func (app *application) assetDropdown(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	term := r.Form.Get("q")
	list, err := app.assets.GetAssetsByNumber(term)
	if err != nil {
		app.serverError(w, err)
		return
	}

	j, err := json.Marshal(list)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
