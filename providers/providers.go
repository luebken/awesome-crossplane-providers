package providers

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/luebken/awesome-crossplane-providers/util"
)

type ProviderImplementationType string

type FullRepoName string
type SortedFullRepoName []FullRepoName

func (a SortedFullRepoName) Len() int           { return len(a) }
func (a SortedFullRepoName) Less(i, j int) bool { return a[i] < a[j] }
func (a SortedFullRepoName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (frn FullRepoName) GetOwner() string {
	return strings.Split(string(frn), "/")[0]
}
func (frn FullRepoName) GetRepo() string {
	return strings.Split(string(frn), "/")[1]
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

// Get latest repo info for "providerNames"
func GetRepositories(client *github.Client, ctx context.Context, providerNames []FullRepoName) []*github.Repository {
	repos := []*github.Repository{}

	ch := make(chan *github.Repository)
	defer func() {
		close(ch)
	}()

	for _, pr := range providerNames {
		// workaround for socket: too many open files error
		time.Sleep(20 * time.Millisecond)
		go func(pr FullRepoName) {
			fmt.Print("r")
			repo, _, err := client.Repositories.Get(ctx, pr.GetOwner(), pr.GetRepo())
			if err != nil {
				fmt.Printf("\nError querying 'https://github.com/%v/%v\n", pr.GetOwner(), pr.GetRepo())
				fmt.Printf("\n%v\n", err)
				ch <- nil
			} else {
				ch <- repo
			}
		}(pr)
	}

	for i := 0; i < len(providerNames); i++ {
		repo := <-ch
		if repo != nil {
			repos = append(repos, repo)
		}
	}

	fmt.Printf("\nQueried %d repos.\n", len(repos))
	return repos
}

// All provider names from "providers.txt"
func getProviderNamesFromFile() []FullRepoName {
	lines, err := util.ReadFromFile("providers.txt")
	if err != nil {
		panic(err)
	}
	providerNames := []FullRepoName{}
	for _, l := range lines {
		if !strings.HasPrefix(l, "#") && len(l) > 1 {
			providerNames = append(providerNames, FullRepoName(l))
		}
	}
	fmt.Printf("Found %d providers in providers.txt\n", len(providerNames))
	return providerNames
}

// All provider names from "providers.txt" excluding "providers-ignored.txt"
func ReadProviderNamesFromFileWithoutIgnored() []FullRepoName {
	lines, err := util.ReadFromFile("providers.txt")
	if err != nil {
		panic(err)
	}
	linesIgnored, err := util.ReadFromFile("providers-ignored.txt")
	if err != nil {
		panic(err)
	}
	ignoreList := strings.Join(linesIgnored, " ")
	providerNames := []FullRepoName{}
	for _, l := range lines {
		if !strings.HasPrefix(l, "#") && !strings.Contains(ignoreList, l) {
			providerNames = append(providerNames, FullRepoName(l))
		}
	}
	fmt.Printf("Found %d providers in providers.txt excluding providers-ignored.txt\n", len(providerNames))
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

	reposNames := []FullRepoName{}
	for _, r := range repos {
		reposNames = append(reposNames, FullRepoName(*r.FullName))
	}

	oldRepoNames := getProviderNamesFromFile()
	reposNames = append(reposNames, oldRepoNames...)

	reposNames = removeDuplicateRepos(reposNames)
	sort.Sort(SortedFullRepoName(reposNames))
	fmt.Printf("\nFound %d provider repos.\n", len(reposNames))
	s := ""
	for _, rn := range reposNames {
		s = s + string(rn) + "\n"
	}
	util.WriteToFile(s, "providers.txt")
}

func removeDuplicateRepos(repos []FullRepoName) []FullRepoName {
	keys := make(map[FullRepoName]bool)
	list := []FullRepoName{}
	for _, repo := range repos {
		if _, value := keys[repo]; !value {
			keys[repo] = true
			list = append(list, repo)
		}
	}
	return list
}
