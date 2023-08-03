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

package selfupdate

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/edgefarm/edgefarm/cmd/local-up/cmd"
	"github.com/fatih/color"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const (
	// GithubOrgaRepo is the github orga and repo name
	GithubOrgaRepo = "edgefarm/edgefarm"
)

func CheckNewVersion() (string, string, bool, error) {
	latest, found, err := selfupdate.DetectLatest(GithubOrgaRepo)
	if err != nil {
		return "", "", false, fmt.Errorf("CheckNewVersion: error occurred while detecting version: %s", err.Error())
	}
	// Only check if the version is not a dev version
	if cmd.Version != cmd.DevVersion {
		v, err := semver.Parse(cmd.Version)
		if err != nil {
			return "", "", false, fmt.Errorf("CheckNewVersion: error occurred while parsing version: %s, maybe this is a dev version?", cmd.Version)
		}
		if !found || latest.Version.LTE(v) {
			return "", "", false, nil
		}
		return cmd.Version, latest.Version.String(), true, nil
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
