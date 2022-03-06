package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type OwnerRepo struct {
	Owner string
	Repo  string
}

// writes file to /data
func WriteToFile(s string) {
	filename := fmt.Sprintf("repo-stats-%s.csv", time.Now().Format("2006-01-02"))
	full_filename := "/data/" + filename
	fmt.Println("\nWant to write to " + full_filename)
	err := os.WriteFile(full_filename, []byte(s), 0644)
	if err != nil {
		fmt.Println("Err WriteFile ", err)
	} else {
		fmt.Println("Wrote to " + full_filename)
	}
}

type repoCRD struct {
	Path     string
	Filename string
	Group    string
	Version  string
	Kind     string
	CRD      []byte
}

type orgData struct {
	Repo  string
	Tag   string
	At    string
	Tags  []string
	CRDs  map[string]repoCRD
	Total int
}

type CRDsStats struct {
	Total int
	Alpha int
	Beta  int
	V1    int
}

func GetNumberOfCRDs(url string) (cRDsStats *CRDsStats) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	s := string(buf)
	r, _ := regexp.Compile(`JSON\.parse\(.(.*).\)`)
	rawMatch := r.FindStringSubmatch(s)[1]

	var orgData orgData
	err = json.Unmarshal([]byte(rawMatch), &orgData)
	if err != nil {
		log.Fatal(err)
	}

	cRDsStats = new(CRDsStats)
	cRDsStats.Total = orgData.Total

	for _, crd := range orgData.CRDs {
		if strings.Contains(crd.Version, "alpha") {
			cRDsStats.Alpha += 1
		}
		if strings.Contains(crd.Version, "beta") {
			cRDsStats.Beta += 1
		}
		if crd.Version == "v1" {
			cRDsStats.V1 += 1
		}
	}
	//TODO calculate alpha / beta resources
	return cRDsStats
}
