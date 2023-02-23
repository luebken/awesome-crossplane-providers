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
	if len(os.Args) != 2 {
		fmt.Println("Usage axxpp [provider-names|provider-stats]")
		return
	}
	if !(os.Args[1] == "provider-names" || os.Args[1] == "provider-stats") {
		fmt.Println("Usage axxpp [provider-names|provider-stats]")
		return
	}

	if os.Args[1] == "provider-names" {
		providers.UpdateProviderNamesToFile(client, ctx)
		return
	}

	//	os.Args[1] == "provider-stats"
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
		time.Sleep(20 * time.Millisecond)
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
				crds, err := util.GetNumberOfCRDsFromCRDsDev(ps.DocsURL)
				if err == nil {
					ps.CRDsTotal = crds.Total
					ps.CRDsAlpha = crds.Alpha
					ps.CRDsBeta = crds.Beta
					ps.CRDsV1 = crds.V1
				}
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
	fmt.Println("End")
}
