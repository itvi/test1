package main

import (
	"ams/pkg/models"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type templateData struct {
	CurrentYear     int
	Asset           *models.Asset
	AssetCategory   *models.AssetCategory
	AssetCategories []*models.AssetCategory
}

var functions = template.FuncMap{
	"safe": func(s string) template.HTMLAttr {
		return template.HTMLAttr(s)
	},
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"ownedRoles": func(role string, roles []string) string {
		for _, r := range roles {
			if r == role {
				return "checked"
			}
		}
		return ""
	},
	"GBs": func(size string) string {
		if s, err := strconv.ParseFloat(size, 32); err == nil {
			return fmt.Sprintf("%.2f", s/1024/1024/1024) + "GB"
		}
		return ""
	},
	"GB": func(size int64) string {
		return fmt.Sprintf("%.2f", float64(size)/1024/1024/1024) + "GB"
	},
	"datetime": func(dt string) string {
		// 2019-09-17 13:17:58+08:00
		d := strings.Split(dt, "+")
		if len(d) == 2 {
			return d[0]
		}
		return ""
	},
	"fd": func(s string) string { // 20140319102030.000000+480
		str := strings.Split(s, ".")[0]
		date := str[0:4] + "-" + str[4:6] + "-" + str[6:8]
		time := str[8:10] + ":" + str[10:12] + ":" + str[12:14]
		return date + " " + time // 2014-03-19 10:20:30
	},
	// serial number
	"sn": func(i int) int {
		return i + 1
	},
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the sub directories. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*/*.html"))
	if err != nil {
		return nil, err
	}

	fmt.Println("pages:", pages)

	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file pat
		// and assign it to the name variable.
		path := strings.Replace(page, dir, "", -1) //  "asset/index.html"
		// fmt.Println("path:",path)
		name := filepath.Base(page)

		// Parse the page template file in to a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		//ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}
		// Add the template set to the cache, using the name with path of the page
		// (like 'home/index.html') as the key.
		cache[path] = ts
	}
	// Return the map.
	return cache, nil
}
