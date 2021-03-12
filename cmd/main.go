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

package main

import (
	"fmt"

	arg "github.com/alexflint/go-arg"
	"github.com/gw-tester/ip-discover/pkg/discover"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type args struct {
	Log     logLevel `arg:"env:LOG_LEVEL" default:"warn" help:"Defines the level of logging for this program."`
	Network string   `help:"Network CIDR to monitor the creation and existence of IP addresses."`
}

type logLevel struct {
	Level log.Level
}

func (n *logLevel) UnmarshalText(b []byte) error {
	s := string(b)

	logLevel, err := log.ParseLevel(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse the log level")
	}

	n.Level = logLevel

	return nil
}

func (args) Version() string {
	return "ip-discovery 0.0.1"
}

func (args) Description() string {
	return "this program discovers IP address of a given Network range"
}

func main() {
	var args args

	arg.MustParse(&args)
	log.SetLevel(args.Log.Level)

	if ip, err := discover.GetIPFromNetwork(args.Network); err == nil {
		//nolint:forbidigo
		fmt.Println(ip)
	}
}
