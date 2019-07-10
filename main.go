package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func main() {
	build := flag.Int("b", 0, "circle build id")
	expand := flag.Bool("e", false, "print expanded links")
	repo := flag.String("r", "influxdata/telegraf", "vcs username/reponame")
	vcs := flag.String("t", "github", "vcs type")

	flag.Parse()

	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://circleci.com/api/v1.1/project/%s/%s/%d/artifacts",
			*vcs, *repo, *build), nil)
	if err != nil {
		fmt.Println("Failed to create request", err)
		return
	}

	req.SetBasicAuth(os.Getenv("CIRCLE_TOKEN"), "")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to make request", err)
		return
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	Afacts := []struct {
		URL string
	}{}

	err = dec.Decode(&Afacts)
	if err != nil {
		fmt.Println("Failed to decode", err)
		return
	}

	urls := make([]string, len(Afacts))
	for i := range Afacts {
		urls[i] = Afacts[i].URL
	}

	artifacts := map[string]map[string][]info{}
	geese := make(map[string][]info)

	for i := range urls {
		dir, file := path.Split(urls[i][strings.Index(urls[i], "/build/")+7:])
		base, err := url.PathUnescape(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		dir, file = path.Split(path.Dir(dir))
		if _, ok := artifacts[dir]; !ok {
			geese = map[string][]info{}
		}
		geese[file] = append(geese[file], info{url: urls[i], file: base})
		artifacts[dir] = geese
	}

	if *expand {
		printExpanded(artifacts)
	} else {
		printCollapsable(artifacts)
	}
}

type info struct {
	url  string
	file string
}

func printCollapsable(artifacts map[string]map[string][]info) {
	fmt.Println("<ul>")
	for goos, v := range artifacts {
		fmt.Printf("  <li><details><summary>%s</summary><ul>\n", goos)
		for arch, v := range v {
			fmt.Printf("    <li><details><summary>%s</summary><ul>\n", arch)
			for _, i := range v {
				fmt.Printf("      <li><a href=\"%s\">%s</a></li>\n", i.url, i.file)
			}
			fmt.Println("    </li></ul></details>")
		}
		fmt.Println("  </li></ul></details>")
		fmt.Println()
	}
	fmt.Println("</ul>")
}

func printExpanded(artifacts map[string]map[string][]info) {
	for goos, v := range artifacts {
		fmt.Printf("* %s\n", goos)
		for arch, v := range v {
			fmt.Printf("  - %s\n", arch)
			for _, i := range v {
				fmt.Printf("    + [%s](%s)\n", i.file, i.url)
			}
		}
	}
}
