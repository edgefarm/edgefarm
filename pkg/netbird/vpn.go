package netbird

import (
	"errors"

	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	"k8s.io/klog/v2"
)

func DisableVPN(uninstall, groupDel, routeDel, peerDel, setypKeysDel bool) error {
	state, err := state.GetState()
	if err != nil {
		return err
	}
	args.NetbirdSetupKey = state.GetNetbirdSetupKey()
	args.NetbirdToken = state.GetNetbirdToken()
	if args.NetbirdToken == "" {
		return errors.New("cluster is not VPN enabled. Please run 'local-up vpn enable' first")
	}

	if uninstall {
		if args.NetbirdToken != "" {
			klog.Infof("Uninstall VPN packages")
			if err := packages.Uninstall(shared.KubeConfigRestConfig, packages.VPN); err != nil {
				return err
			}
		}
	}

	klog.Infof("netbird.io: Cleanup")
	err = Cleanup(state, args.NetbirdToken, groupDel, routeDel, peerDel, setypKeysDel)
	if err != nil {
		return err
	}
	state.Clear()

	return nil
}

func EnableVPN() error {
	state, err := state.GetState()
	if err != nil {
		return err
	}
	state.SetNetbirdToken(args.NetbirdToken)
	klog.Infof("Prepare edge nodes")
	klog.Info("Start to prepare netbird")
	key, err := CreateSetupKey(state, args.NetbirdToken)
	if err != nil {
		return err
	}
	args.NetbirdSetupKey = key.Key

	if args.NetbirdToken != "" {
		klog.Infof("Deploy cluster bootstrap VPN packages")
		if err := packages.Install(shared.KubeConfigRestConfig, packages.VPN); err != nil {
			return err
		}
	}

	klog.Infof("Configuring netbird")
	err = AddRoute(state, args.NetbirdToken)
	if err != nil {
		return err
	}
	state.SetFullyConfigured()
	return nil
}
