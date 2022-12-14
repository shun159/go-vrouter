// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import (
	"bytes"
	"fmt"

	"github.com/shun159/vr/vr"
)

type RouteOption func(*vr.VrRouteReq)

func RouteVrfId(vrf_id int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrVrfID = vrf_id
	}
}

func RouteFamily(family int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrFamily = family
	}
}

func RoutePrefix(prefix []int8) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrPrefix = prefix
	}
}

func RoutePrefixLen(p_len int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrPrefixLen = p_len
	}
}

func RouteRid(rid int16) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrRid = rid
	}
}

func RouteLabelFlags(flags int16) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrLabelFlags = flags
	}
}

func RouteLabel(label int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrLabel = label
	}
}

func RouteNhId(nh_id int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrNhID = nh_id
	}
}

func RouteMarker(marker []int8) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrMarker = marker
	}
}

func RouteMarkerPlen(marker_plen int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrMarkerPlen = marker_plen
	}
}

func RouteMac(mac []int8) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrMac = mac
	}
}

func RouteReplacePlen(replace_plen int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrReplacePlen = replace_plen
	}
}

func RouteIndex(index int32) RouteOption {
	return func(args *vr.VrRouteReq) {
		args.RtrIndex = index
	}
}

func (vr_msg *VrMessage) AddRoute(setters ...RouteOption) (int32, error) {
	r := vr.NewVrRouteReq()
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
		errmsg := fmt.Errorf("failed to create route with non-zero resp-code: %v", resp_code)
		return resp_code, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg *VrMessage) GetRoute(setters ...RouteOption) (*vr.VrRouteReq, error) {
	r := vr.NewVrRouteReq()
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
		errmsg := fmt.Errorf("failed to create route with non-zero resp-code: %v", resp_code)
		return nil, errmsg
	}

	rt := vr.NewVrRouteReq()
	if err := rt.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse binary into vr_route: %s", err)
		return nil, errmsg
	}

	return rt, nil
}

func (vr_msg *VrMessage) DelRoute(setters ...RouteOption) (int32, error) {
	r := vr.NewVrRouteReq()
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
		errmsg := fmt.Errorf("failed to delete route with non-zero resp-code: %v", resp_code)
		return -1, errmsg
	}

	return vr_resp.RespCode, nil
}

func (vr_msg VrMessage) DumpRoute(setters ...RouteOption) ([]vr.VrRouteReq, error) {
	r := vr.NewVrRouteReq()
	r.HOp = vr.SandeshOp_DUMP

	defer vr_msg.sandesh.protocol.ReadI16(vr_msg.sandesh.context)

	for _, setter := range setters {
		setter(r)
	}

	route_list := []vr.VrRouteReq{}
	vr_resp, multipart, err := vr_msg.syncMultipart(r)
	if err != nil {
		return route_list, err
	}

	if vr_resp.RespCode < 0 {
		resp_code := vr_resp.RespCode
		errmsg := fmt.Errorf("failed to dump route. non-zero resp-code: %v", resp_code)
		return route_list, errmsg
	}

	for _, m := range multipart {
		buf := bytes.NewBuffer(m.data)
		vr_msg.sandesh.transport.Buffer = buf
		for vr_msg.sandesh.transport.Buffer.Len() > 4 {
			nh := vr.NewVrRouteReq()
			if err := nh.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
				fmt.Printf("failed to parse vr_route: %v", err)
				break
			}
			route_list = append(route_list, *nh)
		}
	}

	return route_list, nil
}
