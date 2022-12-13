// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import (
	"bytes"
	"fmt"

	"github.com/shun159/vr/vr"
)

type NexthopOption func(*vr.VrNexthopReq)

func NhMarker(marker int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrMarker = marker
	}
}

func NhRid(rid int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrRid = rid
	}
}

func NhID(nh_id int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrID = nh_id
	}
}

func NhFlags(flags int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrFlags = flags
	}
}

func NhType(nh_type int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrType = nh_type
	}
}

func NhFamily(family int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrFamily = family
	}
}

func NhVrf(vrf int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrVrf = vrf
	}
}

func NhEncapOifID(if_id []int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrEncapOifID = if_id
	}
}

func NhEncapFamily(family int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrEncapFamily = family
	}
}

func NhEncap(encap []int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrEncap = encap
	}
}

func NhTunSip(ip int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrTunSip = ip
	}
}

func NhTunDip(ip int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrTunDip = ip
	}
}

func NhTunSip6(ip []int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrTunSip6 = ip
	}
}

func NhTunDip6(ip []int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrTunDip6 = ip
	}
}

func NhTunSport(port int16) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrTunSport = port
	}
}

func NhTunDport(port int16) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrTunDport = port
	}
}

func NhEncapValid(valid_list []int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrEncapValid = valid_list
	}
}

func NhCryptTraffic(crypt int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrCryptTraffic = crypt
	}
}

func NhCryptPathAvailable(crypt_path_available int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrCryptPathAvailable = crypt_path_available
	}
}

func NhTransportLabel(label int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrTransportLabel = label
	}
}

func NhRwDstMac(mac []int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrRwDstMac = mac
	}
}

func NhPbbMac(mac []int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrPbbMac = mac
	}
}

func NhEcmpConfigHash(hash int8) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrEcmpConfigHash = hash
	}
}

func NhNhList(nh []int32) NexthopOption {
	return func(args *vr.VrNexthopReq) {
		args.NhrNhList = nh
	}
}

func (vr_msg VrMessage) DumpNexthop(setters ...NexthopOption) ([]vr.VrNexthopReq, error) {
	r := vr.NewVrNexthopReq()
	r.HOp = vr.SandeshOp_DUMP
	r.NhrMarker = -1

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	nh_list := []vr.VrNexthopReq{}
	vr_resp, multipart, err := vr_msg.syncMultipart(r)
	if err != nil {
		return nh_list, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to dump nexthop. non-zero resp-code: %v", resp_code)
		return nh_list, errmsg
	}

	for _, m := range multipart {
		buf := bytes.NewBuffer(m.data)
		vr_msg.sandesh.transport.Buffer = buf
		for vr_msg.sandesh.transport.Buffer.Len() > 8 {
			nh := vr.NewVrNexthopReq()
			if err := nh.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
				fmt.Printf("failed to parse vr_nexthop: %v", err)
				break
			}
			nh_list = append(nh_list, *nh)
		}
	}

	return nh_list, nil

}

func (vr_msg *VrMessage) GetNexthop(setters ...NexthopOption) (*vr.VrNexthopReq, error) {
	r := vr.NewVrNexthopReq()
	r.HOp = vr.SandeshOp_GET

	defer vr_msg.sandesh.protocol.ReadI32(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vr_resp, err := vr_msg.sync(r)
	if err != nil {
		return nil, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to get nexthop. non-zero resp-code: %v", resp_code)
		return nil, errmsg
	}

	nh := vr.NewVrNexthopReq()
	if err := nh.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse binary into vr_nexthop_req: %s", err)
		return nil, errmsg
	}

	return nh, nil
}

func (vr_msg *VrMessage) AddNexthop(setters ...NexthopOption) (int32, error) {
	r := vr.NewVrNexthopReq()
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
		errmsg := fmt.Errorf("failed to create nexthop with non-zero resp-code: %v", resp_code)
		return resp_code, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg *VrMessage) DelNexthop(setters ...NexthopOption) (int32, error) {
	r := vr.NewVrNexthopReq()
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
		errmsg := fmt.Errorf("failed to delete nexthop with non-zero resp-code: %v", resp_code)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}
