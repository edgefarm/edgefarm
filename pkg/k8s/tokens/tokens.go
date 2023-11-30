/*
Copyright 2020 The OpenYurt Authors.

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

package tokens

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	kubeclientset "k8s.io/client-go/kubernetes"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	"k8s.io/klog/v2"

	bootstraptokenv1 "github.com/openyurtio/openyurt/pkg/util/kubernetes/kubeadm/app/apis/bootstraptoken/v1"
	kubeadmconstants "github.com/openyurtio/openyurt/pkg/util/kubernetes/kubeadm/app/constants"
	nodetoken "github.com/openyurtio/openyurt/pkg/util/kubernetes/kubeadm/app/phases/bootstraptoken/node"
)

func GetOrCreateJoinTokenString(cliSet kubeclientset.Interface) (string, error) {
	tokenSelector := fields.SelectorFromSet(
		map[string]string{
			// TODO: We hard-code "type" here until `field_constants.go` that is
			// currently in `pkg/apis/core/` exists in the external API, i.e.
			// k8s.io/api/v1. Should be v1.SecretTypeField
			"type": string(bootstrapapi.SecretTypeBootstrapToken),
		},
	)
	listOptions := metav1.ListOptions{
		FieldSelector: tokenSelector.String(),
	}
	klog.V(1).Infoln("[token] retrieving list of bootstrap tokens")
	secrets, err := cliSet.CoreV1().Secrets(metav1.NamespaceSystem).List(context.Background(), listOptions)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, "failed to list bootstrap tokens")
	}

	for _, secret := range secrets.Items {

		// Get the BootstrapToken struct representation from the Secret object
		token, err := bootstraptokenv1.BootstrapTokenFromSecret(&secret)
		if err != nil {
			klog.Warningf("%v", err)
			continue
		}
		if !usagesAndGroupsAreValid(token) {
			continue
		}

		return token.Token.String(), nil
		// Get the human-friendly string representation for the token
	}

	tokenStr, err := bootstraputil.GenerateBootstrapToken()
	if err != nil {
		return "", fmt.Errorf("couldn't generate random token, %w", err)
	}
	token, err := bootstraptokenv1.NewBootstrapTokenString(tokenStr)
	if err != nil {
		return "", err
	}

	klog.V(1).Infoln("[token] creating token")
	if err := nodetoken.CreateNewTokens(cliSet,
		[]bootstraptokenv1.BootstrapToken{{
			Token:  token,
			Usages: kubeadmconstants.DefaultTokenUsages,
			Groups: kubeadmconstants.DefaultTokenGroups,
		}}); err != nil {
		return "", err
	}
	return tokenStr, nil
}

// usagesAndGroupsAreValid checks if the usages and groups in the given bootstrap token are valid
func usagesAndGroupsAreValid(token *bootstraptokenv1.BootstrapToken) bool {
	sliceEqual := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		sort.Strings(a)
		sort.Strings(b)
		for k, v := range b {
			if a[k] != v {
				return false
			}
		}
		return true
	}

	return sliceEqual(token.Usages, kubeadmconstants.DefaultTokenUsages) && sliceEqual(token.Groups, kubeadmconstants.DefaultTokenGroups)
}
