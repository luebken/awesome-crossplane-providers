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

	result := "Repository;URL;Description;Stars;Subscribers;Open Issues;Last Update;Last Release;Docs\n"
	for _, oR := range util.ReadReposFromFile() {
		fmt.Print(".")
		repo, _, err := client.Repositories.Get(ctx, oR.Owner, oR.Repo)
		if err != nil {
			fmt.Println("err ", err)
		}
		fmt.Print("|")
		release, _, err := client.Repositories.GetLatestRelease(ctx, oR.Owner, oR.Repo)
		last_release := ""
		docs_url := ""
		if err != nil {
			//fmt.Println("err ", err)
		} else {
			last_release = release.CreatedAt.Time.Format("2006-01-02")
			docs_url = "https://doc.crds.dev/github.com/" + *repo.GetOwner().Login + "/" + *repo.Name + "@" + *release.TagName
		}
		if !*repo.Archived {
			desc := ""
			if repo.Description != nil {
				desc = *repo.Description
			}
			result = result + fmt.Sprintf("%s;%s;%s;%d;%d;%d;%v;%v;%v\n",
				*repo.FullName,
				*repo.HTMLURL,
				desc,
				*repo.StargazersCount,
				*repo.SubscribersCount,
				*repo.OpenIssuesCount,
				repo.UpdatedAt.Time.Format("2006-01-02"),
				last_release,
				docs_url)
		}
	}
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
