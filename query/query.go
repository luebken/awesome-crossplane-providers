package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v42/github"
)

var (
	providerSearchQueries = []string{
		`archived:false in:readme,name,description "is a crossplane provider"`,
		`archived:false in:readme,name,description "is a minimal crossplane provider"`,
		`archived:false in:readme,name,description "is an experimental crossplane provider"`,
		`archived:false in:readme,name,description "crossplane infrastructure provider"`,
	}
	providerOrgs = []string{
		"crossplane-contrib", "crossplane",
	}
	blacklist = "https://github.com/terrytangyuan/awesome-argo " + " "
)

// Query potential Provider repos by
// a) search for common phrases in readme, name, description
// b) look for "provider" repos in orgs
func QueryPotentialProviderRepos(client *github.Client, ctx context.Context) (repos []*github.Repository) {
	fmt.Println("Searching for potential Crossplane provider repos.")

	for _, q := range providerSearchQueries {
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

	for _, o := range providerOrgs {
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
	repos = removeReposWhichAreOnTheBlacklist(repos)
	fmt.Printf("\n Found %d provider repos.\n", len(repos))
	return repos
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

func removeReposWhichAreOnTheBlacklist(repos []*github.Repository) []*github.Repository {
	list := []*github.Repository{}
	for _, repo := range repos {
		if !strings.Contains(blacklist, *repo.HTMLURL) {
			list = append(list, repo)
		}
	}
	return list
}
