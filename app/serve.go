package app

import (
	"html/template"
	"net/http"

	"github.com/spf13/viper"

	"github.com/BNPrashanth/h24-assignment/models"
)

func HandleServeIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indexHtml := models.IndexHTML{
		BaseURL: viper.GetString("BASR_URL"),
	}

	err = tmpl.Execute(w, indexHtml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
