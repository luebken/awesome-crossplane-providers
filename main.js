const https = require('https')

const fs = require('fs');
const readline = require('readline');

// cat README.md | ggrep -oP '\[(.*\/.*)\]' > repos.txt
async function processLineByLine() {
    const fileStream = fs.createReadStream('repos.txt');

    const rl = readline.createInterface({
        input: fileStream,
        crlfDelay: Infinity
    });
    let repos = []
    for await (const line of rl) {
        repos.push(line)
    }
    return repos
}

let options = {
    hostname: 'api.github.com',
    port: 443,
    method: 'GET',
    headers: {
        "User-Agent": "luebken-awesome-operators",
        "Authorization": "token " + process.env.GITHUB_TOKEN
    }
}


async function queryAll() {
    repos = await processLineByLine();
    console.log(repos)

    let stats = []
    for (const reponame of repos) {
        let stat = await queryRepoStats(reponame)
        if (stat != "archived") {
            let last_release = await queryRepoRelease(reponame)
            stat.last_release = last_release
            stats.push(stat)
        }
    }
    // sort by stargazers
    stats = stats.sort((a, b) => b.stargazers_count - a.stargazers_count);

    console.log("\n\n\n---")
    console.log("Providers with at least one release:")
    console.log("| Github | Description | License | Stargazers | Last Update | Last Release |")
    console.log("|--------|-------------|---------|------------|-------------|--------------|")
    stats.forEach(stat => {
        if (stat.last_release.published_at) {

            let s = "| [" + stat.full_name + "](https://github.com/" + stat.full_name + ")"
                + " | " + stat.description
                + " | " + stat.license
                + " | " + stat.stargazers_count
                + " | " + stat.updated_at.split('T')[0]

            if (stat.last_release.published_at) {
                s += " | " + stat.last_release.name + " " + stat.last_release.published_at.split('T')[0]
            } else {
                s += " | No release yet"
            }
            s += " |"
            console.log(s)
        }
    });
}

async function queryRepoStats(reponame, cb) {
    return new Promise((resolve, reject) => {
        options.path = '/repos/' + reponame
        process.stdout.write("-");
        const req = https.request(options, res => {
            let body = "";
            let status = res.statusCode
            let stats = {}
            res.on("data", (chunk) => {
                body += chunk;
            });

            res.on("end", () => {
                try {
                    let json = JSON.parse(body);
                    if (status == 200) {
                        if (!json.archived) {
                            stats.full_name = json.full_name
                            stats.description = json.description
                            stats.stargazers_count = json.stargazers_count
                            if (json.license)
                                stats.license = json.license.spdx_id
                            stats.updated_at = json.updated_at

                            resolve(stats)
                        } else {
                            //console.log(reponame + " archived")
                            resolve("archived")
                        }
                    } else {
                        console.log(reponame + " status " + status)
                        resolve({})
                    }
                } catch (error) {
                    console.error(error.message);
                };

            });
        })

        req.on('error', error => {
            console.error(error)
        })

        req.end()

    })

}

async function queryRepoRelease(reponame, cb) {
    return new Promise((resolve, reject) => {
        options.path = '/repos/' + reponame + "/releases"
        process.stdout.write(".");
        const req = https.request(options, res => {
            let body = "";
            let status = res.statusCode
            let last_release = {}
            res.on("data", (chunk) => {
                body += chunk;
            });

            res.on("end", () => {
                try {
                    let json = JSON.parse(body);
                    if (status == 200) {
                        if (!json.archived) {
                            if (json.length != null && json.length > 0) {
                                var json0 = json[0]
                                last_release.name = json0.name
                                last_release.html_url = json0.html_url
                                last_release.published_at = json0.published_at
                                resolve(last_release)
                            }
                            resolve({})

                        } else {
                            resolve({})
                        }
                    } else {
                        console.log(reponame + " status " + status)
                        resolve({})
                    }
                } catch (error) {
                    console.error(error.message);
                };

            });
        })

        req.on('error', error => {
            console.error(error)
        })

        req.end()

    })

}

queryAll()