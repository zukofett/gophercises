package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	urlshortner "github.com/zukofett/gophercises/url_shortner/handler"
)

func main() {
	YamlFile := flag.String("yaml", "url_paths.yaml", "a file with path to url mappings")
	JsonFile := flag.String("json", "url_paths.json", "a file with path to url mappings")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshortner.MapHandler(pathsToUrls, mux)

	yaml_file, err := os.Open(*YamlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer yaml_file.Close()

	yaml, err := io.ReadAll(yaml_file)
	if err != nil {
		log.Fatal(err)
	}

	json_file, err := os.Open(*JsonFile)
	if err != nil {
		log.Fatal(err)
	}
	defer json_file.Close()

	json, err := io.ReadAll(json_file)
	if err != nil {
		log.Fatal(err)
	}

	jsonHandler, err := urlshortner.JSONHandler([]byte(json), mapHandler)
    if err != nil {
        panic(err)
    }

	yamlHandler, err := urlshortner.YAMLHandler([]byte(yaml), jsonHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}
