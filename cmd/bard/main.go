// (c) 2020 Emir Erbasan (humanova)
// MIT License, see LICENSE for more details

package main

import (
	"bard/config"
	"fmt"
	"github.com/gorilla/mux"
	"bard/internal/handler"
	"log"
	"net/http"
)

func main() {
	var conf config.Config

	err := config.GetConfig("./config.json", &conf)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	cmsRouter := r.PathPrefix(conf.ListenPrefixPath).Subrouter()

	cmsRouter.HandleFunc("/create_post", handler.CreatePostHandler(conf)).Methods("POST")
	cmsRouter.HandleFunc("/update_post", handler.UpdatePostHandler(conf)).Methods("POST")
	cmsRouter.HandleFunc("/delete_post", handler.DeletePostHandler(conf)).Methods("POST")
	cmsRouter.HandleFunc("/cms", handler.CmsHandler(conf)).Methods("GET")

	addr := fmt.Sprintf(":%s", conf.Port)
	fmt.Printf("Starting the server on %s\n", fmt.Sprintf("%s%s", addr, conf.ListenPrefixPath))

	err = http.ListenAndServe(addr, cmsRouter)
	log.Fatal(err)
}


