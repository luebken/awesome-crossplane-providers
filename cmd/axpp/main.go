package main

import (
	"context"
	"fmt"
	"os"

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

func main() {
	fmt.Println("Start")

	result := "Repository;URL;Description;Stars;Subscribers;Open Issues;Last Update;Created;Last Release;Docs;CRDs Total;CRDs Alpha;CRDs Beta;CRDs V1\n"

	providersTotal := 0
	providersAlpha := 0
	providersBeta := 0
	providersV1 := 0
	crdsTotalAlpha := 0
	crdsTotalBeta := 0
	crdsTotalV1 := 0
	crdsTotalTotal := 0

	for _, repo := range query.QueryPotentialProviderRepos(client, ctx) {
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
		if !*repo.Archived { // double check
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

	fmt.Println("End")
}
