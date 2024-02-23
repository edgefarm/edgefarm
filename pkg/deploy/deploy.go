package deploy

import (
	"fmt"
	"os"
	"time"

	configv1 "github.com/edgefarm/edgefarm/pkg/config/v1alpha1"
	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/openyurt"
	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

func AddFlagsForDeploy(flagset *pflag.FlagSet) {
	flagset.BoolVar(&shared.Args.Skip.EdgeFarmCore, "skip-openyurt", false, "Skip installaing of openyurt components. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Skip.EdgeFarmApplications, "skip-applications", false, "Skip installing edgefarm.applications.")
	flagset.BoolVar(&shared.Args.Skip.EdgeFarmNetwork, "skip-network", false, "Skip installing edgefarm.network.")
	flagset.BoolVar(&shared.Args.Skip.EdgeFarmMonitor, "skip-monitor", false, "Skip installing edgefarm.monitor.")
	flagset.BoolVar(&shared.Args.Skip.EdgeFarmCore, "skip-core", false, "Skip installing edgefarm.core.")
	flagset.BoolVar(&shared.Args.Only.EdgeFarmApplications, "only-applications", false, "Only install edgefarm.applications.")
	flagset.BoolVar(&shared.Args.Only.EdgeFarmNetwork, "only-network", false, "Only install edgefarm.network.")
	flagset.BoolVar(&shared.Args.Only.EdgeFarmMonitor, "only-monitor", false, "Only install edgefarm.monitor.")
	flagset.BoolVar(&shared.Args.Only.EdgeFarmCore, "only-core", false, "Only install edgefarm.core.")

	if os.Getenv("LOCAL_UP_EXPERIMENTAL") == "true" {
		flagset.BoolVar(&shared.Args.Skip.Crossplane, "skip-crossplane", false, "Skip installing crossplane. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Skip.Kyverno, "skip-kyverno", false, "Skip installing kyverno. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Skip.Metacontroller, "skip-metacontroller", false, "Skip installing metacontroller. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Skip.VaultOperator, "skip-vault-operator", false, "Skip installing vault-operator. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Skip.Vault, "skip-vault", false, "Skip installing vault. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Skip.Ingress, "skip-ingress", false, "Skip installing ingress-nginx. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Skip.CertManager, "skip-cert-manager", false, "Skip installing cert-manager. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Only.Crossplane, "only-crossplane", false, "Only install crossplane. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Only.Kyverno, "only-kyverno", false, "Only install kyverno. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Only.Metacontroller, "only-metacontroller", false, "Only install metacontroller. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Only.VaultOperator, "only-vault-operator", false, "Only install vault-operator. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Only.Vault, "only-vault", false, "Only install vault. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Only.Ingress, "only-ingress", false, "Only install ingress-nginx. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
		flagset.BoolVar(&shared.Args.Only.CertManager, "only-cert-manager", false, "Only install cert-manager. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")

	}
}

func Deploy(t configv1.ConfigType, config *rest.Config) error {
	if !shared.Args.Skip.Ingress {
		klog.Infoln("Deploy ingress packages")
		if err := packages.Install(config, packages.Ingress); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.CertManager {
		klog.Infoln("Deploy cert-manager packages")
		if err := packages.Install(config, packages.CertManager); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.Kyverno {
		klog.Infoln("Deploy kyverno packages")
		if err := packages.Install(config, packages.Kyverno); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.Crossplane {
		klog.Infoln("Deploy crossplane packages")
		if err := packages.Install(config, packages.Crossplane); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.Metacontroller {
		klog.Infoln("Deploy metacontroller packages")
		if err := packages.Install(config, packages.Metacontroller); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.VaultOperator {
		klog.Infoln("Deploy vault-operator packages")
		if err := packages.Install(config, packages.VaultOperator); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.Vault {
		klog.Infoln("Deploy vault packages")
		if err := packages.Install(config, packages.Vault); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.EdgeFarmCore {
		klog.Info("Start to deploy OpenYurt components")
		openyurtDeployer := &openyurt.DeployOpenYurt{
			YurthubHealthCheckTimeout: 2 * time.Minute,
			YurthubImage:              fmt.Sprintf(constants.YurtHubImageFormat, constants.OpenYurtVersion),
			YurtManagerImage:          fmt.Sprintf(constants.YurtManagerImageFormat, constants.OpenYurtVersion),
			NodeServantImage:          fmt.Sprintf(constants.NodeServantImageFormat, constants.OpenYurtVersion),
			EnableDummyIf:             true,
		}
		if err := openyurtDeployer.Run(t, config); err != nil {
			klog.Errorf("errors occurred when deploying openyurt components")
			return err
		}
	}

	if !shared.Args.Skip.EdgeFarmApplications {
		klog.Infof("Deploy edgefarm applications packages")
		if err := packages.Install(config, packages.Applications); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.EdgeFarmNetwork {
		klog.Infof("Deploy edgefarm network packages")
		if err := packages.Install(config, packages.Network); err != nil {
			return err
		}
	}

	if !shared.Args.Skip.EdgeFarmMonitor {

		klog.Infof("Deploy edgefarm applications packages")
		if err := packages.Install(config, packages.Monitor); err != nil {
			return err
		}
	}

	green := color.New(color.FgHiGreen)
	yellow := color.New(color.FgHiYellow)
	green.Printf("The local EdgeFarm cluster is ready to use! Have fun exploring EdgeFarm.\n")
	green.Println("To access the cluster use 'kubectl', e.g.")
	yellow.Println("  $ kubectl get nodes")

	return nil
}
