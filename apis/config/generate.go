//go:build generate
// +build generate

/*
Copyright © 2024 EdgeFarm Authors

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

// Package v1alpha1 implements the v1alpha1 apiVersion of edgefarm's cluster configuration

// +genclient
// +k8s:deepcopy-gen=package
//go:generate go run k8s.io/code-generator/cmd/deepcopy-gen -i ./v1alpha1 -h ../../hack/boilerplate.txt --output-base . --output-package v1alpha1 -O zz_deepcopy

package config

import (
	_ "k8s.io/code-generator/cmd/deepcopy-gen" //nolint:typecheck
)
