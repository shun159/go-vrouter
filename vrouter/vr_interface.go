// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/shun159/vr/vr"
	"golang.org/x/sys/unix"
)

type VifOption func(*vr.VrInterfaceReq)

func VifRid(routerId int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrRid = routerId
	}
}

func VifIdx(vif_idx int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrIdx = vif_idx
	}
}

func VifOsIdx(vif_os_idx int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrOsIdx = vif_os_idx
	}
}

func VifMarker(mark int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrMarker = mark
	}
}

func VifCore(core int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrCore = core
	}
}

func VifFlags(flags int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFlags = flags
	}
}

func VifType(vif_type int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrType = vif_type
	}
}

func VifVlanId(vlan_id int16) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrVlanID = vlan_id
	}
}

func VifOVlanId(vlan_id int16) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrOvlanID = vlan_id
	}
}

func VifParentVifIndex(vif_idx int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrParentVifIdx = vif_idx
	}
}

func VifLoopbackIP(ip int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrLoopbackIP = ip
	}
}

func VifCrossConnectIdx(indexes []int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrCrossConnectIdx = indexes
	}
}

func VifMac(mac []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrMac = mac
	}
}

func VifSrcMac(mac []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrSrcMac = mac
	}
}

func VifFatFlowProtocolPort(ports []int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowProtocolPort = ports
	}
}

func VifFatFlowSrcPrefixH(prefixes []int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowSrcPrefixH = prefixes
	}
}

func VifFatFlowSrcPrefixL(prefixes []int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowSrcPrefixL = prefixes
	}
}

func VifFatFlowSrcPrefixMask(masks []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowSrcPrefixMask = masks
	}
}

func VifFatFlowSrcAggregatePlen(prefix_len []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowSrcAggregatePlen = prefix_len
	}
}

func VifFatFlowDstPrefixH(prefixes []int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowDstPrefixH = prefixes
	}
}

func VifFatFlowDstPrefixL(prefixes []int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowDstPrefixL = prefixes
	}
}

func VifFatFlowDstPrefixMask(masks []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowDstPrefixMask = masks
	}
}

func VifFatFlowDstAggregatePlen(prefix_len []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowDstAggregatePlen = prefix_len
	}
}

func VifFatFlowExcludeIPList(ip_list []int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowExcludeIPList = ip_list
	}
}

func VifFatFlowExcludeIp6LList(ip_list []int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowExcludeIp6LList = ip_list
	}
}

func VifFatFlowExcludeIp6UList(ip_list []int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowExcludeIp6UList = ip_list
	}
}

func VifFatFlowExcludeIp6PlenList(prefix_len []int16) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrFatFlowExcludeIp6PlenList = prefix_len
	}
}

func VifPbbMac(mac []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrPbbMac = mac
	}
}

func VifIsid(isid int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrIsid = isid
	}
}

func VifIp6L(ip int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrIp6L = ip
	}
}

func VifIp6U(ip int64) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrIp6U = ip
	}
}

func VifVhostuserMode(mode int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrVhostuserMode = mode
	}
}

func VifName(name string) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrName = name
	}
}

func VifHwQueues(queues []int16) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrHwQueues = queues
	}
}

func VifMirID(mir_id int16) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrMirID = mir_id
	}
}

func VifInMirrorMd(mirror_md []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrInMirrorMd = mirror_md
	}
}

func VifOutMirrorMd(mirror_md []int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrOutMirrorMd = mirror_md
	}
}

func VifTransport(transport_type int8) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrTransport = transport_type
	}
}

func VifVrf(vrf int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrVrf = vrf
	}
}

func VifMcastVrf(vrf int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrMcastVrf = vrf
	}
}

func VifMtu(mtu int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrMtu = mtu
	}
}

func VifIP(ip int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrIP = ip
	}
}

func VifNhID(nh_id int32) VifOption {
	return func(args *vr.VrInterfaceReq) {
		args.VifrNhID = nh_id
	}
}

func (vr_msg *VrMessage) DumpVif(setters ...VifOption) ([]vr.VrInterfaceReq, error) {
	r := vr.NewVrInterfaceReq()
	r.HOp = vr.SandeshOp_DUMP
	r.VifrIdx = -1

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vifs := []vr.VrInterfaceReq{}
	vr_resp, multipart, err := vr_msg.syncMultipart(r)
	if err != nil {
		return vifs, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to dump interfaces. non-zero resp-code: %v", resp_code)
		return vifs, errmsg
	}

	for _, m := range multipart {
		buf := bytes.NewBuffer(m.data)
		vr_msg.sandesh.transport.Buffer = buf
		for vr_msg.sandesh.transport.Buffer.Len() > 8 {
			vif := vr.NewVrInterfaceReq()
			if err := vif.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
				fmt.Printf("failed to parse vr_interface: %v", err)
				break
			}
			vifs = append(vifs, *vif)
		}
	}

	return vifs, nil
}

func (vr_msg *VrMessage) GetVif(setters ...VifOption) (*vr.VrInterfaceReq, error) {
	r := vr.NewVrInterfaceReq()
	r.HOp = vr.SandeshOp_GET

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vr_resp, err := vr_msg.sync(r)
	if err != nil {
		return nil, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to get interface. non-zero resp-code: %v", resp_code)
		return nil, errmsg
	}

	vif := vr.NewVrInterfaceReq()
	if err := vif.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse binary into vr_interface_req: %s", err)
		return nil, errmsg
	}

	return vif, nil
}

func (vr_msg *VrMessage) AddVif(setters ...VifOption) (int32, error) {
	r := vr.NewVrInterfaceReq()
	r.HOp = vr.SandeshOp_ADD

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vr_resp, err := vr_msg.sync(r)
	if err != nil {
		return -1, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to create interface with non-zero resp-code: %v", resp_code)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg *VrMessage) DelVif(setters ...VifOption) (int32, error) {
	r := vr.NewVrInterfaceReq()
	r.HOp = vr.SandeshOp_DEL

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vr_resp, err := vr_msg.sync(r)
	if err != nil {
		return -1, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to delete interface with non-zero resp-code: %v", resp_code)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg *VrMessage) ResetStatsVif(setters ...VifOption) (*vr.VrInterfaceReq, error) {
	r := vr.NewVrInterfaceReq()
	r.HOp = vr.SandeshOp_RESET
	r.VifrIdx = -1

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vr_resp, err := vr_msg.sync(r)
	if err != nil {
		return nil, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to delete interface with non-zero resp-code: %v", resp_code)
		return nil, errmsg
	}

	vif := vr.NewVrInterfaceReq()
	if err := vif.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse binary into vr_interface_req: %s", err)
		return nil, errmsg
	}

	return vif, nil
}

/*
 * vrouter and vrouter-agent uses tun/tap device
 * so that reception and transmission packet in the userspace
 */
const TUN_INTF_CLONE_DEV = "/dev/net/tun"
const PKT0_IFNAME = "pkt0"

// pkt0 interface
type Pkt0Device struct {
	Index         int
	HardwareAddr  []uint8
	TapFD         int
	TxBufferCount int
}

func (pkt0 *Pkt0Device) Init() error {
	tap_fd, err := unix.Open(
		TUN_INTF_CLONE_DEV,
		os.O_RDWR,
		0,
	)
	if err != nil {
		errmsg := fmt.Errorf("packet tap error: %s openning tap-device", err)
		return errmsg
	}

	ifr, err := unix.NewIfreq(PKT0_IFNAME)
	ifr.SetUint16(unix.IFF_TUN_EXCL | unix.IFF_TAP | unix.IFF_NO_PI)
	if err != nil {
		return err
	}

	if err := unix.IoctlIfreq(
		tap_fd,
		unix.TUNSETIFF,
		ifr,
	); err != nil {
		errmsg := fmt.Errorf("packet tap error: %s creating tap-device", err)
		return errmsg
	}

	if _, err := unix.FcntlInt(
		uintptr(tap_fd),
		unix.F_SETFD,
		unix.FD_CLOEXEC,
	); err != nil {
		errmsg := fmt.Errorf("packet tap error: %s setting fcntl on %s", err, PKT0_IFNAME)
		return errmsg
	}

	ifr.SetUint16(1)
	if err := unix.IoctlIfreq(
		tap_fd,
		unix.TUNSETPERSIST,
		ifr,
	); err != nil {
		errmsg := fmt.Errorf("packet tap error: %s making tap-interface non-persisitent", err)
		return errmsg
	}

	iface, err := net.InterfaceByName(PKT0_IFNAME)
	if err != nil {
		errmsg := fmt.Errorf("packet tap error: %s retriving interface of the tap device", err)
		return errmsg
	}

	raw_fd, err := unix.Socket(
		unix.AF_PACKET,
		unix.SOCK_RAW,
		unix.ETH_P_ALL,
	)
	if err != nil {
		errmsg := fmt.Errorf("packet tap error: %s creating raw socket", err)
		return errmsg
	}
	defer unix.Close(raw_fd)

	sll := &unix.SockaddrLinklayer{
		Protocol: unix.ETH_P_ALL,
		Ifindex:  iface.Index,
	}

	if err := unix.Bind(raw_fd, sll); err != nil {
		errmsg := fmt.Errorf("packet tap error: %s binding the socket to the tap interface", err)
		return errmsg
	}

	ifr.SetUint32(uint32(pkt0.TxBufferCount))
	if err := unix.IoctlIfreq(
		raw_fd,
		unix.SIOCGIFTXQLEN,
		ifr,
	); err != nil {
		errmsg := fmt.Errorf("packet tap error: %s setting tx-buffer size", err)
		return errmsg
	}

	ifr.SetUint16(0)
	if err := unix.IoctlIfreq(
		raw_fd,
		unix.SIOCGIFFLAGS,
		ifr,
	); err != nil {
		errmsg := fmt.Errorf("packet tap error: %s getting socket flags", err)
		return errmsg
	}

	ifr.SetUint16(ifr.Uint16() | unix.IFF_UP)
	if err := unix.IoctlIfreq(
		raw_fd,
		unix.SIOCSIFFLAGS,
		ifr,
	); err != nil {
		errmsg := fmt.Errorf("packet tap error: %s setting socket flags", err)
		return errmsg
	}

	pkt0.Index = iface.Index
	pkt0.HardwareAddr = iface.HardwareAddr
	pkt0.TapFD = tap_fd

	return nil
}

func (pkt0 *Pkt0Device) Close() error {
	err := exec.Command("ip", "link", "del", PKT0_IFNAME).Run()
	if err != nil {
		return err
	}

	return nil
}
