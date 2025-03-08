package urlshortner

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.RequestURI()
		if url, exists := pathToUrls[path]; exists {
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
    parsedYaml, err := parseYAML(yml)
    if err != nil {
        return nil, err
    }

    pathMap := buildMap(parsedYaml)

    return MapHandler(pathMap, fallback), nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
    parsedJson, err := parseJSON(json)
    if err != nil {
        return nil, err
    }

    pathMap := buildMap(parsedJson)

    return MapHandler(pathMap, fallback), nil
}



type urlPath struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func parseYAML(yml []byte) ([]urlPath, error) {
	var urls []urlPath

	err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func parseJSON(js []byte) ([]urlPath, error) {
    var urls []urlPath

	err := json.Unmarshal(js, &urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func buildMap(url_paths []urlPath) map[string]string {
	pathMap := make(map[string]string)

	for _, url := range url_paths {
		pathMap[url.Path] = url.URL
	}

	return pathMap
}

