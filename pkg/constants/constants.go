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

package constants

import "time"

const (
	DefaultKubeConfigPath    = "~/.edgefarm-local-up/kubeconfig"
	BootstrapTokenDefaultTTL = time.Hour * 24
	OpenYurtVersion          = "v1.4.0"
)

var (
	YurtHubImageFormat     = "ghcr.io/openyurtio/openyurt/yurthub:%s"
	YurtManagerImageFormat = "ghcr.io/openyurtio/openyurt/yurt-manager:%s"
	NodeServantImageFormat = "ghcr.io/openyurtio/openyurt/node-servant:%s"
)
