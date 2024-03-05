package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/BNPrashanth/h24-assignment/app"
)

func main() {
	viper.AutomaticEnv()
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running!!")
	}).Methods("GET")

	router.HandleFunc("/", app.HandleServeIndex).Methods("GET")

	router.HandleFunc("/analyse", app.HandleAnalyseWebPage).Methods("POST")

	http.ListenAndServe(":"+viper.GetString("PORT"), router)
}
