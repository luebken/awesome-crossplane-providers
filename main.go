package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/luebken/awesome-crossplane-providers/util"
	"golang.org/x/oauth2"
)

var (
	ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	ctx    = context.Background()
	tc     = oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
)

func main() {

	result := "Repository;URL;Description;Stars;Subscribers;Open Issues;Last Update;Created;Last Release;Docs;CRDs Total;CRDs Alpha;CRDs Beta;CRDs V1\n"

	providersTotal := 0
	providersAlpha := 0
	providersBeta := 0
	providersV1 := 0
	crdsTotalAlpha := 0
	crdsTotalBeta := 0
	crdsTotalV1 := 0
	crdsTotalTotal := 0

	for _, repo := range queryRepos() {
		fmt.Print(".") // progress indicator
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
		if !*repo.Archived {
			desc := ""
			subscribersCount := 0
			if repo.Description != nil {
				desc = *repo.Description
			}
			if repo.SubscribersCount != nil {
				subscribersCount = *repo.SubscribersCount
			}
			result += fmt.Sprintf("%s;%s;%s;%d;%d;%d;%v;%v;%v;%v;%d;%d;%d;%d\n",
				*repo.FullName,
				*repo.HTMLURL,
				desc,
				*repo.StargazersCount,
				subscribersCount,
				*repo.OpenIssuesCount,
				repo.UpdatedAt.Time.Format("2006-01-02"),
				repo.CreatedAt.Time.Format("2006-01-02"),
				last_release,
				docs_url,
				crdsTotal,
				crdsAlpha,
				crdsBeta,
				crdsV1)
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
		}
	}

	summary := fmt.Sprintf("\nProviders Total:;%d\nProviders Alpha:;%d\nProviders Beta:;%d\nProviders V1:;%d\nCRDs Total:;%d\nCRDs Alpha:;%d\nCRDs Beta:;%d\nCRDs V1:;%d\n",
		providersTotal,
		providersAlpha,
		providersBeta,
		providersV1,
		crdsTotalTotal,
		crdsTotalAlpha,
		crdsTotalBeta,
		crdsTotalV1,
	)

	result += summary
	util.WriteToFile(result)
}

func queryRepos() (repos []*github.Repository) {
	fmt.Println("Searching for potential Crossplane provider repos.")

	// search for common phrases in readme, name, description
	searchQueries := []string{
		`archived:false in:readme,name,description "is a crossplane provider"`,
		`archived:false in:readme,name,description "is a minimal crossplane provider"`,
		`archived:false in:readme,name,description "is an experimental crossplane provider"`,
		`archived:false in:readme,name,description "crossplane infrastructure provider"`,
	}
	for _, q := range searchQueries {
		opt := &github.SearchOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		for {
			fmt.Print(".")
			result, resp, err := client.Search.Repositories(ctx, q, opt)
			if err != nil {
				fmt.Println("err ", err)
			}
			repos = append(repos, result.Repositories...)

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}

	// crossplane orgs repos with "provider"
	crossplaneOrgs := []string{
		"crossplane-contrib", "crossplane",
	}
	for _, o := range crossplaneOrgs {
		for {
			opt := &github.RepositoryListOptions{
				ListOptions: github.ListOptions{PerPage: 100},
			}
			fmt.Print(".")
			newRepos, resp, err := client.Repositories.List(ctx, o, opt)
			if err != nil {
				fmt.Println("err ", err)
			}

			for _, r := range newRepos {
				if strings.Contains(*r.Name, "provider") {
					repos = append(repos, r)
				}
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}

	repos = removeDuplicateRepos(repos)
	repos = removeBlacklist(repos)
	fmt.Printf("\nAnalyzing %d provider repos.\n", len(repos))

	return repos
}

func removeBlacklist(repos []*github.Repository) []*github.Repository {
	list := []*github.Repository{}
	blacklist := "https://github.com/terrytangyuan/awesome-argo "
	blacklist += " "
	for _, repo := range repos {
		if !strings.Contains(blacklist, *repo.HTMLURL) {
			list = append(list, repo)
		}
	}
	return list
}

func removeDuplicateRepos(repos []*github.Repository) []*github.Repository {
	keys := make(map[string]bool)
	list := []*github.Repository{}
	for _, repo := range repos {
		if _, value := keys[*repo.HTMLURL]; !value {
			keys[*repo.HTMLURL] = true
			list = append(list, repo)
		}
	}
	return list
}
