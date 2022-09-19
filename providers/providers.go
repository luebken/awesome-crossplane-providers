package providers

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/luebken/awesome-crossplane-providers/util"
)

type ProviderName struct {
	Owner string
	Repo  string
}

func (pn ProviderName) Fullname() string {
	return pn.Owner + "/" + pn.Repo
}

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
)

// Read providers.txt and get latest repo info
func ProviderRepos(client *github.Client, ctx context.Context) []*github.Repository {
	repos := []*github.Repository{}

	ch := make(chan *github.Repository)
	defer func() {
		close(ch)
	}()

	providerNames := getProviderNamesFromFile()
	for _, pr := range providerNames {
		// workaround for socket: too many open files error
		time.Sleep(20 * time.Millisecond)
		go func(pr ProviderName) {
			fmt.Print(".")
			repo, _, err := client.Repositories.Get(ctx, pr.Owner, pr.Repo)
			if err != nil {
				fmt.Printf("Error querying 'https://github.com/%v/%v'. Ignoring. ", pr.Owner, pr.Repo)
				fmt.Printf("\n\nerr %v\n\n", err)
				ch <- nil
			} else {
				ch <- repo
			}
		}(pr)
	}

	for i := 0; i < len(providerNames); i++ {
		fmt.Print(".") // progress indicator
		repo := <-ch
		if repo != nil {
			repos = append(repos, repo)
		}
	}

	fmt.Printf("Queried %d repos.\n", len(repos))
	return repos
}

// All provider names from "providers.txt" excluding "providers-ignored.txt"
func getProviderNamesFromFile() []ProviderName {
	lines, err := util.ReadFromFile("providers.txt")
	if err != nil {
		panic(err)
	}
	linesIgnored, err := util.ReadFromFile("providers-ignored.txt")
	if err != nil {
		panic(err)
	}
	ignoreList := strings.Join(linesIgnored, " ")
	providerNames := []ProviderName{}
	for _, l := range lines {
		if !strings.HasPrefix(l, "#") && !strings.Contains(ignoreList, l) {
			var providerName ProviderName
			providerName.Owner = strings.Split(l, "/")[0]
			providerName.Repo = strings.Split(l, "/")[1]
			providerNames = append(providerNames, providerName)
		}
	}
	fmt.Printf("Found %d unique providers in providers.txt.\n", len(providerNames))
	return providerNames
}

// Update provider.txt
func UpdateProviderNamesToFile(client *github.Client, ctx context.Context) {
	fmt.Println("Searching for potential Crossplane provider repos.")

	repos := []*github.Repository{}

	// from search queries
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

	// from provider orgs
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

	// TODO search from topics

	repos = removeDuplicateRepos(repos)
	reposNames := []string{}
	for _, r := range repos {
		reposNames = append(reposNames, *r.FullName)
	}
	sort.Strings(reposNames)
	// TODO don't write if its on the ingore list
	fmt.Printf("\nFound %d provider repos.\n", len(reposNames))
	s := strings.Join(reposNames, "\n")
	util.WriteToFile(s, "providers.txt")
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
