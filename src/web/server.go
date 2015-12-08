package web

// server.go - web server for API

import (
	"fmt"
	"log"
	"net/http"
	"steamtest/src/util"
)

// Start listening for and responding to HTTP requests via the web server
func Start() {
	cfg, err := util.ReadConfig()
	if err != nil {
		log.Fatalf("Unable to read configuration to start web server for API: %s",
			err)
	}

	r := newRouter(cfg)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.WebConfig.APIWebPort), r))
}