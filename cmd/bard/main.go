// (c) 2020 Emir Erbasan (humanova)
// MIT License, see LICENSE for more details

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
	"log"
	"net/http"
)

type Config struct {
	PostDirectory    string
	ListenPrefixPath string
	Port             string
}

func main() {
	var config Config
	err := gonfig.GetConf("./config.json", &config)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	cmsRouter := r.PathPrefix(config.ListenPrefixPath).Subrouter()

	cmsRouter.HandleFunc("/create_post", config.createPostHandler).Methods("POST")
	cmsRouter.HandleFunc("/update_post", config.updatePostHandler).Methods("POST")
	cmsRouter.HandleFunc("/delete_post", config.deletePostHandler).Methods("POST")
	cmsRouter.HandleFunc("/cms", config.cmsHandler).Methods("GET")

	addr := fmt.Sprintf(":%s", config.Port)
	fmt.Printf("Starting the server on %s\n", fmt.Sprintf("%s%s", addr, config.ListenPrefixPath))

	err = http.ListenAndServe(addr, cmsRouter)
	log.Fatal(err)
}


