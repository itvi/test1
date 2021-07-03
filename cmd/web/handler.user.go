package main

import (
	"ams/pkg/models"
	"ams/pkg/models/forms"
	"net/http"
)

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	page := "user.signup.html"
	if r.Method == "GET" {
		app.render(w, r, page, &templateData{Form: forms.New(nil)})
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("sn", "name", "password")
		form.MatchesPattern("email", forms.EmailReg)
		form.MinLength("password", 3)
		form.MaxLength("sn", 8)

		if !form.Valid() {
			app.render(w, r, "user.signup.html", &templateData{Form: form})
			return
		}
		var user = &models.User{
			SN:    form.Get("sn"),
			Name:  form.Get("name"),
			Email: form.Get("email"),
		}
		password := form.Get("password")
		err := app.users.Insert(user, password)
		if err == models.ErrDuplicate {
			form.Errors.Add("sn", "用户已存在")
			app.render(w, r, page, &templateData{Form: form})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flash", "注册成功，请登录。")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	page := "user.login.html"
	if r.Method == "GET" {
		app.render(w, r, page, &templateData{Form: forms.New(nil)})
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("sn", "password")

		if !form.Valid() {
			app.render(w, r, page, &templateData{Form: form})
			return
		}

		sn := form.Get("sn")
		password := form.Get("password")
		user, err := app.users.Authenticate(sn, password)
		if err == models.ErrInvalidCredentials {
			form.Errors.Add("generic", "用户名或密码不正确！")
			app.render(w, r, page, &templateData{Form: form})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "userID", user.ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "你已成功退出。")
	http.Redirect(w, r, "/", 303)
}
