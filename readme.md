# Awesome Crossplane Providers

This project queries Github to find Awesome Crossplane Providers and generates different stats. For production ready providers please see the [Upbound Marketplace](https://marketplace.upbound.io/providers).

## e.g. [repo-stats-latest.csv](./reports/repo-stats-latest.csv)

### How?

This project consists 2 automated steps: 

#### 1) Generate the list of providers

This is done via the command "`axpp provider-names`" which runs in the Github action [provider-names.yml](.github/workflows/provider-names.yml). It queries Github with a set of pre-defined queries and patterns (see [providers.go](/providers/providers.go)) to generate an alphabetical orderd list of providers and saves them to [provider.txt](provider.txt). The queries are somewhat fuzzy and can include false hits. Therefor we ignore all repositories listed in [providers-ignored.txt](providers-ignored.txt).

#### 2) Update provider statistics

This is done via the command "`axpp provider-stats`" which runs in the Github action [provider-stats.yml](.github/workflows/provider-stats.yml). It reads provider.txt and queries Github for current repository information and release information and http://doc.crds.dev for information about the Providers CRDs. This command generates all artefacts apart from the site.