package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/luebken/awesome-crossplane-providers/query"
	"github.com/luebken/awesome-crossplane-providers/util"
	"golang.org/x/oauth2"
)

var (
	ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("MY_GITHUB_TOKEN")},
	)
	ctx    = context.Background()
	tc     = oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
)

type ProviderStat struct {
	Fullname    string
	HTMLURL     string
	Description string
	Stargazers  int
	Subscribers int
	OpenIssues  int
	UpdatedAt   time.Time
	CreatedAt   time.Time
	LastRelease string
	DocsURL     string
	CRDsTotal   int
	CRDsBeta    int
	CRDsAlpha   int
	CRDsV1      int
}
type ByUpdatedAt []ProviderStat

func (ps ByUpdatedAt) Len() int           { return len(ps) }
func (ps ByUpdatedAt) Less(i, j int) bool { return ps[i].UpdatedAt.After(ps[j].UpdatedAt) }
func (ps ByUpdatedAt) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }

func main() {
	fmt.Println("Start")

	providersTotal := 0
	providersAlpha := 0
	providersBeta := 0
	providersV1 := 0
	crdsTotalAlpha := 0
	crdsTotalBeta := 0
	crdsTotalV1 := 0
	crdsTotalTotal := 0

	ch := make(chan ProviderStat)
	defer func() {
		close(ch)
	}()

	repos := query.QueryPotentialProviderRepos(client, ctx)
	for _, repo := range repos {
		go func(repo *github.Repository) {
			release, _, err := client.Repositories.GetLatestRelease(ctx, *repo.Owner.Login, *repo.Name)
			last_release := ""
			docs_url := ""
			crdsTotal := 0
			crdsAlpha := 0
			crdsBeta := 0
			crdsV1 := 0
			if err != nil {
				//fmt.Println("err ", err)
			} else {
				last_release = release.CreatedAt.Time.Format("2006-01-02")
				docs_url = "https://doc.crds.dev/github.com/" + *repo.GetOwner().Login + "/" + *repo.Name + "@" + *release.TagName
				crdsTotal = util.GetNumberOfCRDs(docs_url).Total
				crdsAlpha = util.GetNumberOfCRDs(docs_url).Alpha
				crdsBeta = util.GetNumberOfCRDs(docs_url).Beta
				crdsV1 = util.GetNumberOfCRDs(docs_url).V1
			}
			ps := ProviderStat{
				Fullname: *repo.FullName,
				HTMLURL:  *repo.HTMLURL,
				//Description:
				Stargazers: *repo.StargazersCount,
				// Subscribers:
				OpenIssues:  *repo.OpenIssuesCount,
				UpdatedAt:   repo.UpdatedAt.Time,
				CreatedAt:   repo.CreatedAt.Time,
				LastRelease: last_release,
				DocsURL:     docs_url,
				CRDsTotal:   crdsTotal,
				CRDsAlpha:   crdsAlpha,
				CRDsBeta:    crdsBeta,
				CRDsV1:      crdsV1,
			}
			if repo.SubscribersCount != nil {
				ps.Subscribers = *repo.SubscribersCount
			}
			if repo.Description != nil {
				ps.Description = strings.Replace(*repo.Description, ",", "", -1)
			}

			providersTotal += 1
			if crdsAlpha > 0 {
				providersAlpha += 1
			}
			if crdsBeta > 0 {
				providersBeta += 1
			}
			if crdsV1 > 0 {
				providersV1 += 1
			}
			crdsTotalAlpha += crdsAlpha
			crdsTotalBeta += crdsBeta
			crdsTotalV1 += crdsTotalV1
			crdsTotalTotal += crdsTotal

			ch <- ps
		}(repo)
	}

	stats := []ProviderStat{}
	for i := 0; i < len(repos); i++ {
		fmt.Print(".") // progress indicator
		ps := <-ch
		stats = append(stats, ps)
	}
	sort.Sort(ByUpdatedAt(stats))

	// repo-stats-%s.csv
	statsString := "Repository,URL,Description,Stars,Subscribers,Open Issues,Last Update,Created,Last Release,Docs,CRDs Total,CRDs Alpha,CRDs Beta,CRDs V1\n"
	for _, ps := range stats {
		statsString += fmt.Sprintf("%s,%s,%s,%d,%d,%d,%v,%v,%v,%v,%d,%d,%d,%d\n",
			ps.Fullname,
			ps.HTMLURL,
			ps.Description,
			ps.Stargazers,
			ps.Subscribers,
			ps.OpenIssues,
			ps.UpdatedAt.Format("2006-01-02"),
			ps.CreatedAt.Format("2006-01-02"),
			ps.LastRelease,
			ps.DocsURL,
			ps.CRDsTotal,
			ps.CRDsAlpha,
			ps.CRDsBeta,
			ps.CRDsV1)

	}

	util.WriteToFile(statsString, fmt.Sprintf("/repo/reports/repo-stats-%s.csv", time.Now().Format("2006-01-02")))
	util.WriteToFile(statsString, "/repo/reports/repo-stats-latest.csv")

	//repo-stats-summary-%s.csv
	summary := fmt.Sprintf("\nProviders Total:,%d\nProviders Alpha:,%d\nProviders Beta:,%d\nProviders V1:,%d\nCRDs Total:,%d\nCRDs Alpha:,%d\nCRDs Beta:,%d\nCRDs V1:,%d\n",
		providersTotal,
		providersAlpha,
		providersBeta,
		providersV1,
		crdsTotalTotal,
		crdsTotalAlpha,
		crdsTotalBeta,
		crdsTotalV1,
	)
	util.WriteToFile(summary, fmt.Sprintf("/repo/reports/repo-stats-summary-%s.csv", time.Now().Format("2006-01-02")))
	util.WriteToFile(summary, "/repo/reports/repo-stats-summary.csv")

	//Readme
	sort.Sort(ByUpdatedAt(stats))

	//readme.md
	readme := "# Released Providers:\n\n"
	readme += "||Updated|CRDs:|Alpha|Beta|V1|\n"
	readme += "|---|---|---|---|---|---|\n"
	for _, ps := range stats {
		if ps.LastRelease != "" {
			readme += fmt.Sprintf("|[%s](%s) - [docs](%s)|%s||%d|%d|%d|\n",
				ps.Fullname, ps.HTMLURL, ps.DocsURL, ps.UpdatedAt.Format("2006-01-02"), ps.CRDsAlpha, ps.CRDsBeta, ps.CRDsV1)
		}
	}
	readme += "\nGenerated at: " + time.Now().Format("2006-01-02")
	util.WriteToFile(readme, "/repo/readme.md")

	//data.js
	datajs := "const columns = [\n"
	datajs += "{ title: 'Provider', field: 'name', filtering:false, render: rowData => <a href={rowData.url}>{rowData.name}</a> },\n"
	datajs += "{ title: '', field: 'docs', filtering:false, render: rowData => { if (rowData.docsURL != '') return <a href={rowData.docsURL}>Docs</a> } },\n"
	datajs += "{ title: 'Updated', field: 'updated', filtering:false},\n"
	datajs += "{\n"
	datajs += "  title: 'CRDs maturity', field: 'crdsMaturity',\n"
	datajs += "  lookup: { Unreleased: 'Unreleased', Alpha: 'Alpha', Beta: 'Beta', V1: 'V1' },\n"
	datajs += "  defaultFilter: ['Alpha', 'Beta', 'V1']\n"
	datajs += "},\n"
	datajs += "{ title: 'CRDs', field: 'crds', filtering:false, type: 'numeric' },\n"
	datajs += "];\n"

	datajs += "const data = [\n"
	for _, ps := range stats {
		crdsMaturity := "Unreleased"
		if ps.CRDsAlpha > 0 {
			crdsMaturity = "Alpha"
		}
		if ps.CRDsBeta > 0 {
			crdsMaturity = "Beta"
		}
		if ps.CRDsV1 > 0 {
			crdsMaturity = "V1"
		}
		datajs += fmt.Sprintf("  {name:'%s', url: '%s', docsURL: '%s','updated': '%s', 'crdsMaturity': '%s', 'crds': '%d',},\n",
			ps.Fullname, ps.HTMLURL, ps.DocsURL, ps.UpdatedAt.Format("2006-01-02"), crdsMaturity, ps.CRDsTotal)
	}
	datajs += "];\n"
	datajs += "const exported = { data: data, columns: columns }\n"
	datajs += "export default exported;\n"

	util.WriteToFile(datajs, "/repo/site/src/data.js")

	fmt.Println("End")
}
