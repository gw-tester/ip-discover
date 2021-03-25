/*
Copyright 2021
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

package discover

import (
	"errors"
	"net"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// ErrInvalidNetwork indicates that an invalid network was passed to discover.
var ErrInvalidNetwork = errors.New("invalid network")

func getLinksByNetwork(network *net.IPNet) ([]netlink.Link, bool) {
	log.WithFields(log.Fields{
		"network": network,
	}).Debug("Finding local routes")

	routes := []netlink.Route{}

	// ip route show table local
	if localRoutes, err := netlink.RouteListFiltered(netlink.FAMILY_V4,
		&netlink.Route{Table: 254}, netlink.RT_FILTER_TABLE); err == nil {
		log.Debug("Adding local routes")

		routes = append(routes, localRoutes...)
	}

	// ip route show table main
	if mainRoutes, err := netlink.RouteListFiltered(netlink.FAMILY_V4,
		&netlink.Route{Table: 255}, netlink.RT_FILTER_TABLE); err == nil {
		log.Debug("Adding main routes")

		routes = append(routes, mainRoutes...)
	}

	log.WithFields(log.Fields{
		"routes": routes,
	}).Debug("Routes found")

	links := []netlink.Link{}

	// Filtering local routes by destination
	for _, route := range routes {
		if route.Dst != nil && route.Dst.IP != nil && route.Dst.IP.String() == network.IP.String() {
			log.WithFields(log.Fields{
				"route": route,
			}).Debug("Route matched")

			link, _ := netlink.LinkByIndex(route.LinkIndex)
			links = append(links, link)
		}
	}

	return links, len(links) > 0
}

func waitForNetworkCreation(network *net.IPNet) []netlink.Link {
	routeUpdates := make(chan netlink.RouteUpdate)
	done := make(chan struct{})

	defer close(done)

	if err := netlink.RouteSubscribe(routeUpdates, done); err != nil {
		log.WithError(err).Panic("Failed to susbscribe to local route change event")
	}

	for {
		update, ok := <-routeUpdates
		if !ok {
			panic("route event closed for some unknown reason, re-subscribe")
		}

		if update.Type == syscall.RTM_NEWROUTE {
			log.WithFields(log.Fields{
				"destination": update.Route.Dst,
				"gateway":     update.Route.Gw,
				"network":     network,
			}).Debug("Route add event received") // sudo ip route add 192.168.1.0/24 via 192.168.0.1 dev eno1

			if links, ok := getLinksByNetwork(network); ok {
				log.Infof("%s network was created", network)

				return links
			}
		}
	}
}

func getFirstIPAddress(links []netlink.Link) *net.IPNet {
	if links == nil || len(links) < 1 {
		return nil
	}

	log.WithFields(log.Fields{
		"links": links,
	}).Debug("Getting first IP address")

	addresses, err := netlink.AddrList(links[0], netlink.FAMILY_V4)
	if err != nil {
		log.WithError(err).Panic("Error getting the IPv4 addresses")
	}

	log.WithFields(log.Fields{
		"addresses": addresses,
		"device":    links[0].Attrs().Name,
	}).Debug("IP addresses retrieved")

	if addresses == nil || len(addresses) < 1 {
		return nil
	}

	return addresses[0].IPNet
}

// GetIPFromNetwork waits until a specific network is created and returns its first IP address.
func GetIPFromNetwork(network string) (*net.IPNet, error) {
	log.Infof("Getting first IP address from %s network", network)

	_, parsedNetwork, err := net.ParseCIDR(network)
	if err != nil {
		log.WithError(err).Errorf("Failed to parse %s network", network)

		return nil, ErrInvalidNetwork
	}

	links, ok := getLinksByNetwork(parsedNetwork)
	if !ok {
		log.Warn("Waiting for creation of the network...")

		links = waitForNetworkCreation(parsedNetwork)
	}

	return getFirstIPAddress(links), nil
}
