package util

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type OwnerRepo struct {
	Owner string
	Repo  string
}

func ReadReposFromFile() []OwnerRepo {
	readFile, err := os.Open("repos.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()

	var repos []OwnerRepo
	for _, line := range lines {
		o := new(OwnerRepo)
		o.Owner = strings.Split(line, "/")[0]
		o.Repo = strings.Split(line, "/")[1]
		repos = append(repos, *o)
	}
	return repos
}

func WriteToFile(s string) {
	d1 := []byte(s)
	err := os.WriteFile("repo-stats.csv", d1, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func GetNumberOfCRDs(url string) (crds string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("body > div > div.content-wrapper > div > p > b").Each(func(i int, s *goquery.Selection) {
		crds = s.Text()
	})

	//TODO calculate alpha / beta resources
	return crds
}
