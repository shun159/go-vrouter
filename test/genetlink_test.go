// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter_test

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/shun159/go-vrouter/vrouter"
	"github.com/shun159/vr"
	"golang.org/x/sys/unix"
)

func TestLookupGenlFamily(t *testing.T) {
	var err error
	vr_msg, _ := vrouter.NewVrMessage()

	if ret, err := vr_msg.UpdateVRouter(
		vrouter.Perfs(1),
		vrouter.Perfr1(1),
	); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ret, err)
	}

	if ret, err := vr_msg.ResetVRouter(); err != nil {
		fmt.Printf("error in reseting vrouter: %+v\n", err)
	} else {
		fmt.Println(ret, err)
	}

	if ret, err := vr_msg.DumpVif(
		vrouter.VifRid(0),
		vrouter.VifCore(0),
		vrouter.VifMarker(-1),
	); err != nil {
		fmt.Printf("error in dump vr_interface: %+v\n", err)
	} else {
		fmt.Printf("dump: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.GetVif(
		vrouter.VifRid(0),
		vrouter.VifIdx(4353),
	); err != nil {
		fmt.Printf("error in get vr_interface: %+v\n", err)
	} else {
		fmt.Printf("get: %+v %+v\n", ret, err)
	}

	pkt0 := &vrouter.Pkt0Device{TxBufferCount: 100}
	if err = pkt0.Init(); err != nil {
		fmt.Printf("err in create control interface: %+v", err)
	} else {
		fmt.Printf("pkt0: %+v", pkt0)
	}

	vif_pkt0_mac := []int8{}
	for _, o := range pkt0.HardwareAddr {
		vif_pkt0_mac = append(vif_pkt0_mac, int8(o))
	}

	if ret, err := vr_msg.AddVif(
		vrouter.VifName("pkt0"),
		vrouter.VifRid(0),
		vrouter.VifType(vr.VIF_TYPE_AGENT),
		vrouter.VifFlags(vr.VIF_FLAG_L3_ENABLED),
		vrouter.VifTransport(vr.VIF_TRANSPORT_SOCKET),
		vrouter.VifOsIdx(int32(pkt0.Index)),
		vrouter.VifMac(vif_pkt0_mac),
	); err != nil {
		fmt.Printf("error in add vr_interface: %+v\n", err)
	} else {
		fmt.Printf("get: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.DumpNexthop(
		vrouter.NhMarker(-1),
	); err != nil {
		fmt.Printf("error in dump vr_nexthop: %+v\n", err)
	} else {
		fmt.Printf("dump: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.AddNexthop(
		vrouter.NhID(1),
		vrouter.NhType(vr.NH_TYPE_ENCAP),
		vrouter.NhFlags(vr.NH_FLAG_VALID|vr.NH_FLAG_MCAST),
		vrouter.NhFamily(unix.AF_INET),
		vrouter.NhRid(0),
		vrouter.NhVrf(0),
		vrouter.NhEncapOifID([]int32{4353}),
		vrouter.NhEncapFamily(0x0806),
		vrouter.NhEncap([]int8{
			0x1, 0x2, 0x3, 0x4, 0x5, 0x6, // Dst MAC address
			0xa, 0xb, 0xc, 0xd, 0xe, 0xf, // Src MAC address
			0x08, 0x00, // ether type
		}),
		vrouter.NhTunSip(0),
		vrouter.NhTunDip(0),
	); err != nil {
		fmt.Printf("error in add vr_nexthop: %+v\n", err)
	} else {
		fmt.Printf("add: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.GetNexthop(
		vrouter.NhID(1),
	); err != nil {
		fmt.Printf("error in get vr_nexthop: %+v\n", err)
	} else {
		fmt.Printf("get: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.AddVxlan(
		vrouter.VxlanVnid(1),
		vrouter.VxlanNhid(1),
	); err != nil {
		fmt.Printf("error in add vr_vxlan: %+v\n", err)
	} else {
		fmt.Printf("add vxlan: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.GetVxlan(
		vrouter.VxlanVnid(1),
	); err != nil {
		fmt.Printf("error in get vr_vxlan: %+v\n", err)
	} else {
		fmt.Printf("get vxlan: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.DumpVxlan(vrouter.VxlanVnid(-1)); err != nil {
		fmt.Printf("error in dump vr_vxlan: %+v\n", err)
	} else {
		fmt.Printf("dump vxlan: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.DelVxlan(
		vrouter.VxlanVnid(1),
	); err != nil {
		fmt.Printf("error in del vr_vxlan: %+v\n", err)
	} else {
		fmt.Printf("del_vxlan: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.DelNexthop(
		vrouter.NhID(1),
	); err != nil {
		fmt.Printf("error in del vr_nexthop: %+v\n", err)
	} else {
		fmt.Printf("del: %+v %+v\n", ret, err)
	}

	if ret, err := vr_msg.ResetStatsVif(); err != nil {
		fmt.Printf("error in reset vr_interface: %+v\n", err)
	} else {
		fmt.Printf("reset: %+v %+v\n", ret, err)
	}

	sk, err := vrouter.OpenNetlinkSocket(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Errorf("generic netlink socket open: %s", err)
	}

	if _, err := sk.LookupGenlFamily("vrouter"); err != nil {
		t.Errorf("lookup genl family: %s", err)
	}
}
