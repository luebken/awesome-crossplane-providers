package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v42/github"
	"github.com/luebken/awesome-crossplane-providers/util"
	"golang.org/x/oauth2"
)

func main() {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	result := "Repository;URL;Description;Stars;Subscribers;Open Issues;Last Update;Last Release;Docs;CRDs Total;CRDs Alpha;CRDs Beta;CRDs V1\n"

	providersTotal := 0
	providersAlpha := 0
	providersBeta := 0
	providersV1 := 0
	crdsTotalAlpha := 0
	crdsTotalBeta := 0
	crdsTotalV1 := 0
	crdsTotalTotal := 0

	for _, oR := range util.ReadReposFromFile() {
		fmt.Print(".") // progress indicator
		repo, _, err := client.Repositories.Get(ctx, oR.Owner, oR.Repo)
		if err != nil {
			fmt.Println("err ", err)
		}
		fmt.Print("|") // progress indicator
		release, _, err := client.Repositories.GetLatestRelease(ctx, oR.Owner, oR.Repo)
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
			if repo.Description != nil {
				desc = *repo.Description
			}
			result += fmt.Sprintf("%s;%s;%s;%d;%d;%d;%v;%v;%v;%d;%d;%d;%d\n",
				*repo.FullName,
				*repo.HTMLURL,
				desc,
				*repo.StargazersCount,
				*repo.SubscribersCount,
				*repo.OpenIssuesCount,
				repo.UpdatedAt.Time.Format("2006-01-02"),
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

// for {
// 	newIssues, resp, err := client.Issues.ListByRepo(ctx, "crossplane", "crossplane", opt)
// 	if err != nil {
// 		fmt.Println("err ", err)
// 	}
// 	issues = append(issues, newIssues...)
// 	if resp.NextPage == 0 {
// 		break
// 	}
// 	opt.Page = resp.NextPage
// }
