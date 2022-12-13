// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import (
	"fmt"

	"github.com/shun159/vr/vr"
	vr_raw "github.com/shun159/vr/vr"
)

type VRouterOption func(*vr.VrouterOps)

func LogLevel(level int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoLogLevel = level
	}
}

func LogTypeEnable(types []int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoLogTypeEnable = types
	}
}

func LogTypeDisable(types []int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoLogTypeDisable = types
	}
}

// Turn on/off packet dumps
func PacketDump(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPacketDump = flag
	}
}

// Turn on/off GRO.
func Perfr(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfr = flag
	}
}

// Turn on/off TSO.
func Perfs(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfs = flag
	}
}

// Turn on|off TCP MSS on packets from VM
func FromVMMssAdj(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoFromVMMssAdj = flag
	}
}

// Turn on|off TCP MSS on packets to VM
func ToVMMssAdj(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoToVMMssAdj = flag
	}
}

// RPS after pulling inner hdr
func Perfr1(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfr1 = flag
	}
}

// RPS after GRO on pkt1
func Perfr2(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfr2 = flag
	}
}

// RPS from phys rx handler
func Perfr3(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfr3 = flag
	}
}

// Pull inner hdr (faster version)
func Perfp(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfp = flag
	}
}

// CPU to send pkts to if perfr1 set
func Perfq1(cpu int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfq1 = cpu
	}
}

// CPU to send pkts to if perfr2 set
func Perfq2(cpu int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfq2 = cpu
	}
}

// CPU to send pkts to if perfr3 set
func Perfq3(cpu int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPerfq3 = cpu
	}
}

// NIC cksum offload for outer UDP hdr
func UDPCoff(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoUDPCoff = flag
	}
}

// Turn on|off flow hold limit
func FlowHoldLimit(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoFlowHoldLimit = flag
	}
}

// Turn on|off MPLS over UDP globally
func Mudp(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoMudp = flag
	}
}

// total burst tokens
func BurstTokens(tokens int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoBurstTokens = tokens
	}
}

// timer interval of burst tokens in ms
func BurstInterval(inv int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoBurstInterval = inv
	}
}

// burst tokens to add at every interval
func BurstStep(step int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoBurstStep = step
	}
}

// priority tagging on the NIC
func PriorityTagging(flag int32) VRouterOption {
	return func(args *vr.VrouterOps) {
		args.VoPriorityTagging = flag
	}
}

func (vr_msg *VrMessage) UpdateVRouter(setters ...VRouterOption) (int32, error) {
	args := &vr_raw.VrouterOps{}
	args.HOp = vr.SandeshOp_ADD
	args.VoLogLevel = -1
	args.VoLogTypeEnable = []int32{}
	args.VoLogTypeDisable = []int32{}
	args.VoPerfr = -1
	args.VoPerfs = -1
	args.VoFromVMMssAdj = -1
	args.VoToVMMssAdj = -1
	args.VoPerfr1 = -1
	args.VoPerfr2 = -1
	args.VoPerfr3 = -1
	args.VoPerfp = -1
	args.VoPerfq1 = -1
	args.VoPerfq2 = -1
	args.VoPerfq3 = -1
	args.VoUDPCoff = -1
	args.VoFlowHoldLimit = -1
	args.VoMudp = -1
	args.VoBurstTokens = -1
	args.VoBurstInterval = -1
	args.VoBurstStep = -1
	args.VoPriorityTagging = -1
	args.VoPacketDump = -1

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(args)
	}

	vr_resp, err := vr_msg.sync(args)
	if err != nil {
		return -1, err
	}

	if vr_resp.RespCode != 0 {
		errmsg := fmt.Errorf("failed to get vrouter params. non-zero resp-code: %v", vr_resp.RespCode)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg *VrMessage) GetVRouter() (*vr.VrouterOps, error) {
	vr_req := &vr_raw.VrouterOps{}
	vr_req.HOp = vr.SandeshOp_GET

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	vr_resp, err := vr_msg.sync(vr_req)
	if err != nil {
		return nil, err
	}

	if vr_resp.RespCode != 0 {
		errmsg := fmt.Errorf("failed to get vrouter params. non-zero resp-code: %v", vr_resp.RespCode)
		return nil, errmsg
	}

	vro := vr_raw.NewVrouterOps()
	if err := vro.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to get vrouter params. parse error: %v", err)
		return nil, errmsg
	}

	return vro, nil
}

func (vr_msg *VrMessage) ResetVRouter() (int32, error) {
	vr_req := &vr_raw.VrouterOps{}
	vr_req.HOp = vr.SandeshOp_RESET

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	vr_resp, err := vr_msg.sync(vr_req)
	if err != nil {
		return -1, err
	}

	if vr_resp.RespCode != 0 {
		errmsg := fmt.Errorf("failed to get vrouter params. non-zero resp-code: %v", vr_resp.RespCode)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}

type HugePageOption func(*vr.VrHugepageConfig)

/*
 * Array of virtual memory addresses of the user space process
 * whose huge pages need to be mapped by Vrouter
 */
func HugepagesMem(hpages []int64) HugePageOption {
	return func(args *vr.VrHugepageConfig) {
		args.VhpMem = hpages
	}
}

// Number of addresses present in the hpages argument
func HugepagesPsize(n_hpages []int32) HugePageOption {
	return func(args *vr.VrHugepageConfig) {
		args.VhpPsize = n_hpages
	}
}

// Array of size of each page that is present in hpages
func HugepagesMemSz(hpage_size []int32) HugePageOption {
	return func(args *vr.VrHugepageConfig) {
		args.VhpMemSz = hpage_size
	}
}

// Buffer containing huge page file paths
func HugepagesFilePaths(hpage_file_paths []int8) HugePageOption {
	return func(args *vr.VrHugepageConfig) {
		args.VhpFilePaths = hpage_file_paths
	}
}

// Number of huge page file paths
func HugepagesFilePathSz(hpage_file_pathsz []int32) HugePageOption {
	return func(args *vr.VrHugepageConfig) {
		args.VhpFilePathSz = hpage_file_pathsz
	}
}

func (vr_msg *VrMessage) HugePageConfig(setters ...HugePageOption) (*vr.VrHugepageConfig, error) {
	args := vr.NewVrHugepageConfig()
	args.VhpOp = vr.SandeshOp_ADD
	args.VhpMem = []int64{}
	args.VhpPsize = []int32{}
	args.VhpMemSz = []int32{}
	args.VhpFilePaths = []int8{}
	args.VhpFilePathSz = []int32{}

	for _, setter := range setters {
		setter(args)
	}

	vr_resp, err := vr_msg.sync(args)
	if err != nil {
		return nil, err
	}

	if vr_resp.RespCode != 0 {
		errmsg := fmt.Errorf("failed to get vrouter params. non-zero resp-code: %v", vr_resp.RespCode)
		return nil, errmsg
	}

	vhp := vr_raw.NewVrHugepageConfig()
	if err := vhp.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to get vrouter params. parse error: %v", err)
		return nil, errmsg
	}

	return vhp, nil
}
