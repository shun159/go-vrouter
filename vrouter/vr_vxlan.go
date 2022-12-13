// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import (
	"bytes"
	"fmt"

	"github.com/shun159/vr/vr"
)

type VxlanOption func(*vr.VrVxlanReq)

func VxlanRid(rid int16) VxlanOption {
	return func(args *vr.VrVxlanReq) {
		args.VxlanrRid = rid
	}
}

func VxlanVnid(vnid int32) VxlanOption {
	return func(args *vr.VrVxlanReq) {
		args.VxlanrVnid = vnid
	}
}

func VxlanNhid(nhid int32) VxlanOption {
	return func(args *vr.VrVxlanReq) {
		args.VxlanrNhid = nhid
	}
}

func (vr_msg *VrMessage) AddVxlan(setters ...VxlanOption) (int16, error) {
	r := vr.NewVrVxlanReq()
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
		return int16(resp_code), errmsg
	}

	return 0, nil
}

func (vr_msg *VrMessage) GetVxlan(setters ...VxlanOption) (*vr.VrVxlanReq, error) {
	r := vr.NewVrVxlanReq()
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
		errmsg := fmt.Errorf("failed to get nexthop with non-zero resp-code: %v", resp_code)
		return nil, errmsg
	}

	vxlan := vr.NewVrVxlanReq()
	if err := vxlan.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse binary into vr_vxlan_req: %s", err)
		return nil, errmsg
	}

	return vxlan, nil
}

func (vr_msg *VrMessage) DelVxlan(setters ...VxlanOption) (int32, error) {
	r := vr.NewVrVxlanReq()
	r.HOp = vr.SandeshOp_DEL

	defer vr_msg.sandesh.protocol.ReadI32(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vr_resp, err := vr_msg.sync(r)
	if err != nil {
		return -1, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to delete vxlan with non-zero resp-code: %v", resp_code)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg VrMessage) DumpVxlan(setters ...VxlanOption) ([]vr.VrVxlanReq, error) {
	r := vr.NewVrVxlanReq()
	r.HOp = vr.SandeshOp_DUMP

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vxlanr_list := []vr.VrVxlanReq{}
	vr_resp, multipart, err := vr_msg.syncMultipart(r)
	if err != nil {
		return vxlanr_list, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to dump vxlan. non-zero resp-code: %v", resp_code)
		return vxlanr_list, errmsg
	}

	for _, m := range multipart {
		buf := bytes.NewBuffer(m.data)
		vr_msg.sandesh.transport.Buffer = buf
		for vr_msg.sandesh.transport.Buffer.Len() > 8 {
			nh := vr.NewVrVxlanReq()
			if err := nh.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
				fmt.Printf("failed to parse vr_vxlan: %v", err)
				break
			}
			vxlanr_list = append(vxlanr_list, *nh)
		}
	}

	return vxlanr_list, nil
}
