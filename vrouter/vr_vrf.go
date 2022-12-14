// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import (
	"bytes"
	"fmt"

	"github.com/shun159/vr/vr"
)

type VrfOption func(*vr.VrVrfReq)

func VrfRid(rid int16) VrfOption {
	return func(vvr *vr.VrVrfReq) {
		vvr.VrfRid = rid
	}
}

func VrfIdx(idx int32) VrfOption {
	return func(vvr *vr.VrVrfReq) {
		vvr.VrfIdx = idx
	}
}

func VrfFlags(flags int32) VrfOption {
	return func(vvr *vr.VrVrfReq) {
		vvr.VrfFlags = flags
	}
}

func VrfHbflVifIdx(idx int32) VrfOption {
	return func(vvr *vr.VrVrfReq) {
		vvr.VrfHbflVifIdx = idx
	}
}

func VrfHbfrVifIdx(idx int32) VrfOption {
	return func(vvr *vr.VrVrfReq) {
		vvr.VrfHbfrVifIdx = idx
	}
}

func VrfMarker(marker int32) VrfOption {
	return func(vvr *vr.VrVrfReq) {
		vvr.VrfMarker = marker
	}
}

func (vr_msg *VrMessage) AddVrfTable(setters ...VrfOption) (int32, error) {
	r := vr.NewVrVrfReq()
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
		errmsg := fmt.Errorf("failed to create vrf_table with non-zero resp-code: %v", resp_code)
		return resp_code, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg *VrMessage) GetVrfTable(setters ...VrfOption) (*vr.VrVrfReq, error) {
	r := vr.NewVrVrfReq()
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
		errmsg := fmt.Errorf("failed to create vrf_table with non-zero resp-code: %v", resp_code)
		return nil, errmsg
	}

	vrf := vr.NewVrVrfReq()
	if err := vrf.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse binary into vr_vrf_req: %s", err)
		return nil, errmsg
	}

	return vrf, nil
}

func (vr_msg *VrMessage) DelVrfTable(setters ...VrfOption) (int32, error) {
	r := vr.NewVrVrfReq()
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
		errmsg := fmt.Errorf("failed to delete vrf_table with non-zero resp-code: %v", resp_code)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg VrMessage) DumpVrfTable(setters ...VrfOption) ([]vr.VrVrfReq, error) {
	r := vr.NewVrVrfReq()
	r.HOp = vr.SandeshOp_DUMP

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	vrf_list := []vr.VrVrfReq{}
	vr_resp, multipart, err := vr_msg.syncMultipart(r)
	if err != nil {
		return vrf_list, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to dump vrf_table. non-zero resp-code: %v", resp_code)
		return vrf_list, errmsg
	}

	for _, m := range multipart {
		buf := bytes.NewBuffer(m.data)
		vr_msg.sandesh.transport.Buffer = buf
		for vr_msg.sandesh.transport.Buffer.Len() > 4 {
			nh := vr.NewVrVrfReq()
			if err := nh.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
				fmt.Printf("failed to parse vr_vrf_req: %v", err)
				break
			}
			vrf_list = append(vrf_list, *nh)
		}
	}

	return vrf_list, nil
}
