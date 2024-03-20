package netbird

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	netbird "github.com/edgefarm/netbird-go"
	wait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
)

const (
	routerPeerName = "edgefarm-control-plane"
)

func CreateSetupKey(identifier string, state *state.CurrentState, token string) (*netbird.SetupKey, error) {
	client := netbird.NewClient(token)
	groups, err := client.ListGroups()
	if err != nil {
		return nil, err
	}
	createGroup := true
	group := netbird.Group{}
	for _, g := range groups {
		if g.ID == state.GetNetbirdGroupID() {
			group = g
			createGroup = false
			klog.Infoln("netbird.io: group already exists")
			break
		}
	}

	if createGroup {
		klog.Infoln("netbird.io: creating group")
		g, err := client.CreateGroup(&netbird.Group{
			Name:       fmt.Sprintf("%s-%s", identifier, GetRandomID(8)),
			PeersCount: 0,
			Peers:      []netbird.GroupPeers{},
		})
		if err != nil {
			return nil, err
		}
		group = *g
		state.SetNetbirdGroupID(g.ID)
	}

	setupKeys, err := client.ListSetupKeys()
	if err != nil {
		return nil, err
	}
	createSetupkey := true
	setupKey := netbird.SetupKey{}
	for _, k := range setupKeys {
		if k.ID == state.GetNetbirdSetupKeyID() {
			setupKey = k
			createSetupkey = false
			klog.Infoln("netbird.io: setup key already exists")
			break
		}
	}

	if createSetupkey {
		klog.Infoln("netbird.io: creating setup-key")
		s, err := client.CreateSetupKey(&netbird.SetupKey{
			Name:       fmt.Sprintf("%s-%s", identifier, GetRandomID(8)),
			ExpiresIn:  8640000,
			Type:       "reusable",
			Revoked:    false,
			AutoGroups: []string{group.ID},
			UsageLimit: 0,
			Ephemeral:  false,
		})
		if err != nil {
			return nil, err
		}
		setupKey = *s
		state.SetNetbirdSetupKey(s.Key)
		state.SetNetbirdSetupKeyID(s.ID)
	}
	return &setupKey, nil
}

func AddRoute(identifier string, state *state.CurrentState, token string) error {
	client := netbird.NewClient(token)
	group, err := client.GetGroup(state.GetNetbirdGroupID())
	if err != nil {
		return err
	}
	klog.Infoln("netbird.io: waiting for routing peer to be available")
	err = WaitForRoutingPeer(token)
	if err != nil {
		return err
	}

	createRoute := true
	routes, err := client.ListRoutes()
	if err != nil {
		return err
	}
	for _, r := range routes {
		if r.NetworkID == identifier {
			createRoute = false
			klog.Infoln("netbird.io: route already exists")
			break
		}
	}
	if createRoute {
		peerId, err := client.GetPeerIdByHostname(routerPeerName)
		if err != nil {
			return err
		}

		klog.Infoln("netbird.io: adding route")
		r, err := client.CreateRoute(&netbird.Route{
			NetworkType: "IPv4",
			Description: "edgefarm local cluster",
			NetworkID:   identifier,
			Enabled:     true,
			Peer:        peerId,
			Network:     "172.254.0.0/16",
			Metric:      9999,
			Masquerade:  true,
			Groups:      []string{group.ID},
		})
		if err != nil {
			return err
		}
		state.SetNetbirdRouteID(r.ID)
	}

	return nil
}

func WaitForRoutingPeer(token string) error {
	client := netbird.NewClient(token)
	err := wait.PollImmediateWithContext(context.Background(), 5*time.Second, 5*time.Minute,
		func(_ context.Context) (bool, error) {
			peers, err := client.ListPeers()
			if err != nil {
				return false, err
			}
			for _, p := range peers {
				if p.Name == routerPeerName {
					return true, nil
				}
			}
			return false, nil
		})

	if err != nil {
		return err
	}
	return nil
}

func Cleanup(state *state.CurrentState, token string, groupDel, routeDel, setypKeysDel bool) error {
	client := netbird.NewClient(token)
	if groupDel {
		groups, err := client.ListGroups()
		if err != nil {
			return err
		}
		for _, g := range groups {
			if g.ID == state.GetNetbirdGroupID() {
				klog.Infof("netbird.io: deleting group %s", g.ID)
				client.DeleteGroup(g.ID)
				state.SetNetbirdGroupID("")
				if err := state.Export(); err != nil {
					return err
				}
			}
		}
	}
	if routeDel {
		routes, err := client.ListRoutes()
		if err != nil {
			return err
		}
		for _, r := range routes {
			if r.ID == state.GetNetbirdRouteID() {
				klog.Infof("netbird.io: deleting route %s", r.ID)
				client.DeleteRoute(r.ID)
				state.SetNetbirdRouteID("")
				if err := state.Export(); err != nil {
					return err
				}
			}
		}
	}
	if setypKeysDel {
		peers, err := client.ListPeers()
		if err != nil {
			return err
		}
		for _, p := range peers {
			if strings.HasPrefix(p.Hostname, "edgefarm-") {
				klog.Infof("netbird.io: deleting peer %s", p.Hostname)
				client.DeletePeer(p.ID)
				if err := state.Export(); err != nil {
					return err
				}
			}
		}
	}
	if setypKeysDel {
		keys, err := client.ListSetupKeys()
		if err != nil {
			return err
		}
		for _, k := range keys {
			if k.ID == state.GetNetbirdSetupKeyID() {
				klog.Infof("netbird.io: deleting setup-key %s", k.ID)
				client.DeleteSetupKey(k.ID)
				state.SetNetbirdSetupKeyID("")
				if err := state.Export(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func RoutingPeerIP(token string) (string, error) {
	client := netbird.NewClient(token)
	peers, err := client.ListPeers()
	if err != nil {
		return "", err
	}
	for _, p := range peers {
		if p.Hostname == routerPeerName {
			return p.IP, nil
		}
	}
	return "", errors.New("routing peer not found")
}

func GetGroupPeers(token string) ([]netbird.Peer, error) {
	client := netbird.NewClient(token)
	peers, err := client.ListPeers()
	if err != nil {
		return nil, err
	}
	state, err := state.GetState(shared.StatePath)
	if err != nil {
		return nil, err
	}
	relevantPeers := []netbird.Peer{}
	groupId := state.GetNetbirdGroupID()
	for _, p := range peers {
		for _, g := range p.Groups {
			if g.ID == groupId {
				relevantPeers = append(relevantPeers, p)
				break
			}
		}
	}

	return relevantPeers, nil
}

func GetGroup(token string, id string) (*netbird.Group, error) {
	client := netbird.NewClient(token)
	return client.GetGroup(id)
}

func GetPeerByHostname(token, hostname string) (*netbird.Peer, error) {
	client := netbird.NewClient(token)
	peers, err := client.ListPeers()
	if err != nil {
		return nil, err
	}
	for _, p := range peers {
		if p.Hostname == hostname {
			return &p, nil
		}
	}
	return nil, errors.New("peer not found")
}
