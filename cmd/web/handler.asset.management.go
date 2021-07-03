package main

import (
	"ams/pkg/models"
	"ams/pkg/models/forms"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (app *application) addAssetManagement(w http.ResponseWriter, r *http.Request) {
	page := "assetManagement.add.html"
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

		form.Required("asset_number")
		form.MaxLength("asset_number", 9)

		if !form.Valid() {
			app.render(w, r, page, &templateData{
				Form: form,
			})
			return
		}

		number := form.Get("asset_number")
		mvt := form.Get("mvt")
		docDateString := form.Get("doc_date")
		fromEmployee := form.Get("from_employee")
		fromLoc := form.Get("from_loc")
		toEmployee := form.Get("to_employee")
		toLoc := form.Get("to_loc")
		qtyString := form.Get("qty")
		remark := form.Get("remark")

		qty, err := strconv.Atoi(qtyString)
		if err != nil {
			app.serverError(w, err)
			return
		}
		docDate, err := time.Parse("2006-01-02", docDateString)
		if err != nil {
			app.serverError(w, err)
			return
		}

		assetMov := &models.AssetManagement{
			Asset:        models.Asset{Number: number},
			Mvt:          mvt,
			Qty:          qty,
			FromLoc:      fromLoc,
			ToLoc:        toLoc,
			FromEmployee: fromEmployee,
			ToEmployee:   toEmployee,
			DocumentDate: docDate,
			Remark:       remark,
		}

		err = app.assetManagement.Add(assetMov)
		if err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flash", "添加成功！")
		http.Redirect(w, r, "/asset/management", http.StatusSeeOther)
	}
}

func (app *application) addAssetMovAndConfig(w http.ResponseWriter, r *http.Request) {
	page := "assetManagement.add.html"
	if r.Method == "GET" {
		app.render(w, r, page, &templateData{
			Form: forms.New(nil),
		})
	}

	if r.Method == "POST" {
		r.Header.Set("Content-Type", "application/json")
		type Info struct {
			Status  string // success|danger|warning
			Message string
		}

		var info Info
		var movAndConfig models.AssetMovAndConfig

		jsonInfo := r.PostFormValue("info")
		err := json.Unmarshal([]byte(jsonInfo), &movAndConfig)
		if err != nil {
			info.Status = "danger"
			info.Message = err.Error()
			json.NewEncoder(w).Encode(&info)
			return
		}

		err = app.assetManagement.AddMovAndConfig(&movAndConfig)
		if err != nil {
			info.Status = "danger"
			info.Message = err.Error()
		} else {
			info.Status = "success"
			info.Message = "操作成功!"

			app.session.Put(r, "flash", info.Message)
		}
		json.NewEncoder(w).Encode(info)
	}
}
