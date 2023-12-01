/*
Copyright © 2023 EdgeFarm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This code is based on https://github.com/acim/github-latest

package selfupdate

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/blang/semver"
	"github.com/edgefarm/edgefarm/cmd/local-up/cmd"
	"github.com/fatih/color"
	"github.com/google/go-github/v30/github"
	version "github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
)

const (
	GithubOrga = "edgefarm"
	GithubRepo = "edgefarm"
)

func httpClient() *http.Client {
	tok := os.Getenv("GITHUB_ACCESS_TOKEN")
	if tok == "" {
		return http.DefaultClient
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tok}, //nolint:exhaustruct
	)

	return oauth2.NewClient(ctx, ts)
}

func CheckNewVersion() (string, string, bool, error) {
	// Only check if the version is not a dev version
	if cmd.Version != cmd.DevVersion {
		client := github.NewClient(httpClient())

		rels, res, err := client.Repositories.ListReleases(context.TODO(), GithubOrga, GithubRepo, nil)
		if err != nil {
			if res.StatusCode == http.StatusNotFound {
				fmt.Printf("Repository %s/%s not found\n", GithubOrga, GithubRepo) //nolint:forbidigo
			}

			fmt.Printf("Error: %v", err) //nolint:forbidigo
		}

		versions := make([]*version.Version, 0, len(rels))

		for _, rel := range rels {
			ver, err := version.NewVersion(*rel.TagName)
			if err != nil {
				continue
			}

			if ver.Prerelease() != "" {
				continue
			}

			versions = append(versions, ver)
		}

		sort.Sort(version.Collection(versions))

		if len(versions) == 0 {
			fmt.Println("No releases found") //nolint:forbidigo
		} else {
			v, err := semver.Parse(cmd.Version)
			if err != nil {
				return "", "", false, fmt.Errorf("CheckNewVersion: error occurred while parsing version: %s, maybe this is a dev version?", cmd.Version)
			}
			detected, err := semver.Parse(versions[len(versions)-1].String())
			if err != nil {
				return "", "", false, err
			}
			if v.LT(detected) {
				return cmd.Version, detected.String(), true, nil
			}
			return cmd.Version, detected.String(), false, nil
		}
	}
	return "", "", false, nil
}

func InformAboutNewVersion(current, latest string) {
	c := color.New(color.FgHiGreen)
	c.Printf("Current version is '%s'. A new version is available (%s). Please consider updating to the latest version.\n", current, latest)
	c.Printf("Go to https://github.com/edgefarm/edgefarm/releases/latest to download the latest version.\n\n")
	c.Printf("You can also use the following command to update to the latest version:\n")
	c.Printf("curl -sfL https://raw.githubusercontent.com/edgefarm/edgefarm/main/install.sh | sh -s -- -b ~/bin\n\n")
}
