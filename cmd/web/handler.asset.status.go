package main

import (
	"ams/pkg/models"
	"ams/pkg/models/forms"
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
)

func (app *application) indexAssetStatus(w http.ResponseWriter, r *http.Request) {
	statuses, err := app.assetStatus.GetStatuses()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "assetStatus.index.html", &templateData{
		AssetStatuses: statuses,
	})
}

func (app *application) addAssetStatus(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		app.render(rw, r, "assetStatus.add.html", &templateData{
			Form: forms.New(nil),
		})
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.clientError(rw, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("name")
		form.MaxLength("name", 2)

		if !form.Valid() {
			app.render(rw, r, "assetStatus.add.html", &templateData{
				Form: form,
			})
			return
		}

		name := form.Get("name")
		status := &models.AssetStatus{Name: name}

		err := app.assetStatus.Create(status)
		if err == models.ErrDuplicate {
			form.Errors.Add("name", "已存在！")
			app.render(rw, r, "assetStatus.add.html", &templateData{Form: form})
			return
		} else if err != nil {
			app.serverError(rw, err)
			return
		}

		app.session.Put(r, "flash", "添加成功！")
		http.Redirect(rw, r, "/asset/statuses", http.StatusSeeOther)
	}
}

func (app *application) editAssetStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(bone.GetValue(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	if r.Method == "GET" {
		status, err := app.assetStatus.GetStatusByID(id)
		if err == models.ErrNoRecord {
			app.notFound(w)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		form := forms.New(nil)
		app.render(w, r, "assetStatus.edit.html", &templateData{
			Form:        form,
			AssetStatus: status,
		})
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("name")
		form.MaxLength("name", 2)

		name := form.Get("name")
		status := &models.AssetStatus{ID: id, Name: name}

		if !form.Valid() {
			app.render(w, r, "assetStatus.edit.html", &templateData{
				Form:        form,
				AssetStatus: status,
			})
			return
		}

		if err = app.assetStatus.Edit(status); err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flash", "更新成功!")
		http.Redirect(w, r, "/asset/statuses", http.StatusSeeOther)
	}
}

func (app *application) deleteAssetStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(bone.GetValue(r, "id"))
	if err != nil {
		app.notFound(w)
		return
	}
	if err = app.assetStatus.Delete(id); err != nil {
		app.serverError(w, err)
		return
	}
	//w.Write([]byte("删除成功！"))
	app.session.Put(r, "flash", "删除成功!")
}
