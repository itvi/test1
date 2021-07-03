package main

import (
	"ams/pkg/models"
	"ams/pkg/models/forms"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
)

func (app *application) assetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.assetCategory.GetCategories()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "assetCategory.index.html", &templateData{
		AssetCategories: categories,
	})
}

func (app *application) createAssetCategory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			app.render(rw, r, "assetCategory.add.html", &templateData{
				Form: forms.New(nil),
			})
		}

		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				app.clientError(rw, http.StatusBadRequest)
				return
			}

			form := forms.New(r.PostForm)
			form.Required("code", "name")
			form.MaxLength("code", 4)

			if !form.Valid() {
				app.render(rw, r, "assetCategory.add.html", &templateData{
					Form: form,
				})
				return
			}
			code := form.Get("code")
			name := form.Get("name")
			category := &models.AssetCategory{Code: code, Name: name}

			err := app.assetCategory.Create(category)
			if err == models.ErrDuplicate {
				form.Errors.Add("generic", "编号或名称已存在！")
				app.render(rw, r, "asset.category.add.html", &templateData{
					Form: form,
				})
				return
			} else if err != nil {
				app.serverError(rw, err)
				return
			}

			app.session.Put(r, "flash", "添加成功！")
			http.Redirect(rw, r, "/asset/categories", http.StatusSeeOther)
		}
	}
}

func (app *application) editAssetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(bone.GetValue(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	if r.Method == "GET" {
		category, err := app.assetCategory.GetCategoryByID(id)
		if err == models.ErrNoRecord {
			app.notFound(w)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		form := forms.New(nil)
		app.render(w, r, "assetCategory.edit.html", &templateData{
			Form:          form,
			AssetCategory: category,
		})
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("code", "name")
		form.MaxLength("code", 4)

		code := form.Get("code")
		name := form.Get("name")
		category := &models.AssetCategory{ID: id, Code: code, Name: name}

		if !form.Valid() {
			app.render(w, r, "assetCategory.edit.html", &templateData{
				Form:          form,
				AssetCategory: category,
			})
			return
		}

		if err = app.assetCategory.Edit(category); err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flash", "更新成功!")
		http.Redirect(w, r, "/asset/categories", http.StatusSeeOther)
	}
}

func (app *application) deleteAssetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(bone.GetValue(r, "id"))
	if err != nil {
		app.notFound(w)
		return
	}
	if err = app.assetCategory.Delete(id); err != nil {
		app.serverError(w, err)
		return
	}
	//w.Write([]byte("删除成功！"))
	app.session.Put(r, "flash", "删除成功!")
}

func (app *application) assetCategoryDropdown(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	term := r.Form.Get("q")
	list, err := app.assetCategory.GetCategoriesByName(term)
	if err != nil {
		app.serverError(w, err)
		return
	}
	json, err := json.Marshal(list)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
