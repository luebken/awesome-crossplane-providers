package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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

const file_prefix = ""

//const file_prefix = "/repo/"

func WriteToFile(s string, filename string) {
	err := os.WriteFile(file_prefix+filename, []byte(s), 0644)
	if err != nil {
		fmt.Println("Err WriteFile ", err)
	} else {
		fmt.Println("Wrote to " + file_prefix + filename)
	}
}

func ReadFromFile(filename string) ([]string, error) {
	bytes, err := os.ReadFile(file_prefix + filename)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("Read file " + file_prefix + filename)
	}
	lines := strings.Split(string(bytes), "\n")
	return lines, nil
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

func GetNumberOfCRDsFromCRDsDev(docsURL string) (cRDsStats *CRDsStats, err error) {
	res, err := http.Get(docsURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	s := string(buf)
	r, err := regexp.Compile(`JSON\.parse\(.(.*).\)`)
	if err != nil {
		return nil, err
	}
	match := r.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil, errors.New("FindStringSubmatch")
	}
	rawMatch := match[1]
	var orgData orgData
	err = json.Unmarshal([]byte(rawMatch), &orgData)
	if err != nil {
		return nil, err
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
	return cRDsStats, nil
}

func DiffToTimeAsHumanReadable(t2 time.Time) string {
	t1 := time.Now()
	diffInDays := int(math.Round(t1.Sub(t2).Hours() / 24))
	diffInWeeks := int(math.Round(t1.Sub(t2).Hours() / 24 / 7))
	diffInMonths := int(math.Round(t1.Sub(t2).Hours() / 24 / 31))
	diffString := ""
	unit := ""

	if diffInDays < 1 {
		diffString = " yesterday"
	} else if diffInDays < 7 {
		if diffInDays < 2 {
			unit = "day"
		} else {
			unit = "days"
		}
		diffString = fmt.Sprintf(" %d %s ago", diffInDays, unit)
	} else if diffInWeeks < 5 {
		if diffInWeeks < 2 {
			unit = "week"
		} else {
			unit = "weeks"
		}
		diffString = fmt.Sprintf(" %d %s ago", diffInWeeks, unit)
	} else if diffInMonths > 3000 { //super hacky checking that this is not a valid date
		diffString = "-"
	} else {
		if diffInMonths < 2 {
			unit = "month"
		} else {
			unit = "months"
		}
		diffString = fmt.Sprintf(" %d %s ago", diffInMonths, unit)
	}

	return diffString

}
