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

package route

import (
	"fmt"
	"net"
)

type Route struct {
	Interface string
	IP        string
}

func GetRoute(i string) (*Route, error) {
	iface := &net.Interface{}
	var err error
	if i == "" {
		iface, err = getDefaultRouteInterface()
		if err != nil {
			return nil, err
		}
	} else {
		iface, err = net.InterfaceByName(i)
		if err != nil {
			return nil, err
		}
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	ip, ok := addrs[0].(*net.IPNet)
	if !ok {
		return nil, err
	}

	return &Route{
		Interface: iface.Name,
		IP:        ip.IP.String(),
	}, nil

}

func getDefaultRouteInterface() (*net.Interface, error) {
	routes, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.Flags&net.FlagUp != 0 && route.Flags&net.FlagLoopback == 0 {
			addrs, err := route.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				ip, ok := addr.(*net.IPNet)
				if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
					return &route, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("default route interface not found")
}
