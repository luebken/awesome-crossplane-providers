#!/bin/sh
touch repos.txt
curl -s "https://api.github.com/search/repositories?q=in%3Areadme+sort%3Aupdated+%22is+a+crossplane+provider%22&per_page=100" | jq -r .items[].full_name >> repos.txt
curl -s "https://api.github.com/search/repositories?q=in%3Areadme+sort%3Aupdated+%22is+a+minimal+crossplane+provider%22&per_page=100" | jq -r .items[].full_name >> repos.txt
curl -s "https://api.github.com/search/repositories?q=in%3Areadme+sort%3Aupdated+%22is+an+experimental+crossplane+provider%22&per_page=100" | jq -r .items[].full_name >> repos.txt
curl -s "https://api.github.com/search/repositories?q=in%3Areadme+sort%3Aupdated+%22crossplane+infrastructure+provider%22&per_page=100" | jq -r .items[].full_name >> repos.txt
curl -s https://api.github.com/orgs/crossplane-contrib/repos?per_page=100 | jq -r .[].full_name | grep "/provider-" >> repos.txt
curl -s https://api.github.com/orgs/crossplane/repos?per_page=100 | jq -r .[].full_name | grep "/provider-" >> repos.txt
# black list
sed -i "" 's/terrytangyuan.awesome-argo//g' repos.txt
sort -u repos.txt -o repos.txt