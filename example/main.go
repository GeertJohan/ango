package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"net/http"
	"os"
)

func main() {
	httpFiles, err := rice.FindBox("http-files")
	if err != nil {
		fmt.Printf("Error opening http filex box: %s\n", err)
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(httpFiles.HTTPBox()))

	err = http.ListenAndServe(":8123", nil)
	if err != nil {
		fmt.Printf("Error listenAndServe: %s\n", err)
		os.Exit(1)
	}
}
