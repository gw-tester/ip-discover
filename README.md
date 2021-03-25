# IP Discover
[![Go Report Card](https://goreportcard.com/badge/github.com/gw-tester/ip-discover)](https://goreportcard.com/report/github.com/gw-tester/ip-discover)
[![GoDoc](https://godoc.org/github.com/gw-tester/ip-discover?status.svg)](https://godoc.org/github.com/gw-tester/ip-discover)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Super-Linter](https://github.com/gw-tester/ip-discover/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)

## Summary

This project provides a Go library to retrieve the first local
IP address of a given Network range (CIDR). If there is no local IP Address
is detected the process waits until a new route is created.
