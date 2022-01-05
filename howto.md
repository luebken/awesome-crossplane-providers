# Sources
curl -s https://api.github.com/orgs/crossplane-contrib/repos?per_page=100 | jq -r .[].full_name | grep "/provider-" | pbcopy
curl -s https://api.github.com/orgs/crossplane/repos?per_page=100 | jq -r .[].full_name | grep "/provider-" | pbcopy
jianh619/provider-instana
jianh619/provider-ssh
jianh619/cloudpak-provider