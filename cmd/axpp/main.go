package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/luebken/awesome-crossplane-providers/providers"
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
	Fullname       string
	HTMLURL        string
	Description    string
	Stargazers     int
	Subscribers    int
	OpenIssues     int
	UpdatedAt      time.Time
	CreatedAt      time.Time
	LastReleaseAt  time.Time
	LastReleaseTag string
	DocsURL        string
	CRDsTotal      int
	CRDsBeta       int
	CRDsAlpha      int
	CRDsV1         int
}
type ByUpdatedAt []ProviderStat

func (ps ByUpdatedAt) Len() int           { return len(ps) }
func (ps ByUpdatedAt) Less(i, j int) bool { return ps[i].UpdatedAt.After(ps[j].UpdatedAt) }
func (ps ByUpdatedAt) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }

func main() {
	fmt.Println("Start")

	// TODO don't do this every run but selectively with manual PR
	providers.UpdateProviderNamesToFile(client, ctx)

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

	repos := providers.ProviderRepos(client, ctx)
	for _, repo := range repos {
		go func(repo *github.Repository) {
			ps := ProviderStat{
				Fullname:   *repo.Owner.Login + " / " + *repo.Name, // Extra whitepace for readability
				HTMLURL:    *repo.HTMLURL,
				Stargazers: *repo.StargazersCount,
				OpenIssues: *repo.OpenIssuesCount,
				UpdatedAt:  repo.UpdatedAt.Time,
				CreatedAt:  repo.CreatedAt.Time,
			}
			if repo.Description != nil {
				// remove ","" for csv export
				ps.Description = strings.Replace(*repo.Description, ",", "", -1)
			}

			release, _, err := client.Repositories.GetLatestRelease(ctx, *repo.Owner.Login, *repo.Name)
			if err != nil {
				//fmt.Println("err ", err)
			} else {
				ps.LastReleaseAt = release.CreatedAt.Time
				ps.DocsURL = "https://doc.crds.dev/github.com/" + *repo.GetOwner().Login + "/" + *repo.Name + "@" + *release.TagName
				crds := util.GetNumberOfCRDs(ps.DocsURL)
				ps.CRDsTotal = crds.Total
				ps.CRDsAlpha = crds.Alpha
				ps.CRDsBeta = crds.Beta
				ps.CRDsV1 = crds.V1
				ps.LastReleaseAt = release.CreatedAt.Time
				ps.LastReleaseTag = *release.TagName
			}

			//TODO doesn't work
			if repo.SubscribersCount != nil {
				ps.Subscribers = *repo.SubscribersCount
			}

			// summary
			providersTotal += 1
			if ps.CRDsAlpha > 0 {
				providersAlpha += 1
			}
			if ps.CRDsBeta > 0 {
				providersBeta += 1
			}
			if ps.CRDsV1 > 0 {
				providersV1 += 1
			}
			crdsTotalAlpha += ps.CRDsAlpha
			crdsTotalBeta += ps.CRDsBeta
			crdsTotalV1 += ps.CRDsV1
			crdsTotalTotal += ps.CRDsV1

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
			ps.LastReleaseAt.Format("2006-01-02"),
			ps.DocsURL,
			ps.CRDsTotal,
			ps.CRDsAlpha,
			ps.CRDsBeta,
			ps.CRDsV1)

	}
	util.WriteToFile(statsString, fmt.Sprintf("reports/repo-stats-%s.csv", time.Now().Format("2006-01-02")))
	util.WriteToFile(statsString, "reports/repo-stats-latest.csv")

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
	util.WriteToFile(summary, fmt.Sprintf("reports/repo-stats-summary-%s.csv", time.Now().Format("2006-01-02")))
	util.WriteToFile(summary, "reports/repo-stats-summary.csv")

	sort.Sort(ByUpdatedAt(stats))

	//released-providers.md
	readme := "# Released Providers:\n\n"
	readme += "||Updated|CRDs:|Alpha|Beta|V1|\n"
	readme += "|---|---|---|---|---|---|\n"
	for _, ps := range stats {
		if ps.LastReleaseTag != "" {
			readme += fmt.Sprintf("|[%s](%s) - [docs](%s)|%s||%d|%d|%d|\n",
				ps.Fullname, ps.HTMLURL, ps.DocsURL, ps.UpdatedAt.Format("2006-01-02"), ps.CRDsAlpha, ps.CRDsBeta, ps.CRDsV1)
		}
	}
	readme += "\nGenerated at: " + time.Now().Format("2006-01-02")
	util.WriteToFile(readme, "released-providers.md")

	//data.js
	datajs := "const columns = [\n"
	datajs += "{ title: 'Provider', field: 'name', filtering:false, headerStyle: {minWidth: 400}, cellStyle: {minWidth: 400}, render: rowData => <a href={rowData.url}>{rowData.name}</a> },\n"
	datajs += "{ title: '', field: 'docs', filtering:false, render: rowData => { if (rowData.docsURL !== '') return <a href={rowData.docsURL}>Docs</a> } },\n"
	datajs += "{ title: 'Updated', field: 'updated', filtering:false},\n"
	datajs += "{\n"
	datajs += "  title: 'CRDs maturity', field: 'crdsMaturity',\n"
	datajs += "  lookup: { Unreleased: 'Unreleased', Alpha: 'Alpha', Beta: 'Beta', V1: 'V1' },\n"
	datajs += "  defaultFilter: ['Alpha', 'Beta', 'V1']\n"
	datajs += "},\n"
	datajs += "{ title: 'CRDs', field: 'crdsTotal', filtering:false, type: 'numeric' },\n"
	datajs += "{ title: '', field: 'description', hidden:true, searchable:true, sorting:false, filtering: false	},\n"
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

		//TODO create proper json
		datajs += fmt.Sprintf("  {name:'%s', description:'%s', url: '%s', docsURL: '%s','updated': '%s', 'created': '%s',  'lastReleaseDate': '%s', 'lastReleaseTag': '%s', 'crdsMaturity': '%s', 'crdsAlpha': '%d','crdsBeta': '%d','crdsV1': '%d','crdsTotal': '%d','stargazers': '%d','subscribers': '%d','openIssues': '%d',},\n",
			ps.Fullname,
			ps.Description,
			ps.HTMLURL,
			ps.DocsURL,
			util.DiffToTimeAsHumanReadable(ps.UpdatedAt),
			util.DiffToTimeAsHumanReadable(ps.CreatedAt),
			util.DiffToTimeAsHumanReadable(ps.LastReleaseAt),
			ps.LastReleaseTag,
			crdsMaturity,
			ps.CRDsAlpha,
			ps.CRDsBeta,
			ps.CRDsV1,
			ps.CRDsTotal,
			ps.Stargazers,
			ps.Subscribers,
			ps.OpenIssues,
		)
	}
	datajs += "];\n"
	datajs += fmt.Sprintf("const exported = { data: data, columns: columns, date: '%s' }\n", time.Now().Format("2006-01-02"))
	datajs += "export default exported;\n"
	util.WriteToFile(datajs, "site/src/data.js")

	fmt.Println("End")
}
