package netbird

import (
	"errors"

	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	"k8s.io/klog/v2"
)

func DisableVPN(uninstall, groupDel, routeDel, setypKeysDel bool) error {
	state, err := state.GetState(shared.StatePath)
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
	err = Cleanup(state, args.NetbirdToken, groupDel, routeDel, setypKeysDel)
	if err != nil {
		return err
	}
	state.Clear()

	return nil
}

func Preconfigure(identifier string) (string, error) {
	state, err := state.GetState(shared.StatePath)
	if err != nil {
		return "", err
	}
	state.SetNetbirdToken(args.NetbirdToken)
	key, err := CreateSetupKey(identifier, state, args.NetbirdToken)
	if err != nil {
		return "", err
	}
	args.NetbirdSetupKey = key.Key
	err = state.Export()
	if err != nil {
		return "", err
	}

	return key.Key, nil
}

func UnPreconfigure() error {
	state, err := state.GetState(shared.StatePath)
	if err != nil {
		return err
	}
	return Cleanup(state, args.NetbirdToken, false, false, true)
}

func EnableVPN(identifier string) error {
	klog.Info("Preconfiguring netbird")
	if _, err := Preconfigure(identifier); err != nil {
		return err
	}

	state, err := state.GetState(shared.StatePath)
	if err != nil {
		return err
	}

	if args.NetbirdToken != "" {
		klog.Infof("Deploy cluster bootstrap VPN packages")
		if err := packages.Install(shared.KubeConfigRestConfig, packages.VPN); err != nil {
			return err
		}
	}

	klog.Infof("Configuring netbird")
	err = AddRoute(identifier, state, args.NetbirdToken)
	if err != nil {
		return err
	}
	state.SetFullyConfigured()
	return nil
}
