go-vrouter: A Go library to control the vrouter in-kernel datapath
----

___under development___

## Background

vrouter consists two parts:

1. vrouter is a part of tungstenfabric component, and vrouter in-kernel module is a version of built for run in Linux kernel. The vrouter is controlled from userspace program via sandesh binary protocol over netlink.

2. A userspace daemon that manage the vrouter(known as vrouter-agent or vnsw), settings the flows, routing entries, VRF, nexthops and etc, and handing any misses reported by the vrouter when a packet doesn't match any flows or routing entries.

This library allows Go program to control the vrouter directly via netlink.

## Introduction

## Installation

```golang
import "github.com/shun159/go-vrouter/vrouter"
```
