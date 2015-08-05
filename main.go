package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var result Result

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "ae2c7e54221371785fd119d5208655119cc4c3d9"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated"}
	repos, _, err := client.Repositories.List("mediamath", opt)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	languages := make(map[string]Language)

	for _, repo := range repos {

		fmt.Println(*repo.Name)
		langs, _, err := client.Repositories.ListLanguages("Mediamath", *repo.Name)

		if err != nil {
			fmt.Errorf(err.Error())
		}

		for langStr := range langs {
			if lang, ok := languages[langStr]; ok {
				lang.Repos = append(lang.Repos, *repo.Name)
				languages[langStr] = lang
			} else {
				languages[langStr] = Language{Name: langStr, Repos: make([]string, 1, 100)}
				languages[langStr].Repos[0] = *repo.Name
			}
		}
	}

	langArray := make([]Language, len(languages))
	i := 0
	for _, lang := range languages {
		langArray[i] = lang
		i++
	}

	result = Result{langArray}

	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(result)
	if err == nil {
		w.Write(data)
	} else {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
	}
}

type Language struct {
	Name  string   `json:"language"`
	Repos []string `json:"repos"`
}

type Result struct {
	Languages []Language `json:"data"`
}
