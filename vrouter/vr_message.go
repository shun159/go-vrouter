// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"syscall"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/shun159/vr"
	vr_raw "github.com/shun159/vr/vr"
)

// Sandesh protocol and transport
type Sandesh struct {
	context   context.Context
	transport *thrift.TMemoryBuffer
	protocol  *vr.TSandeshProtocol
}

// Instantiate Sandesh protocol
func newSandesh() *Sandesh {
	mem_buffer := thrift.NewTMemoryBuffer()
	vrouter := vr.NewTSandeshProtocolTransport(mem_buffer)
	sandesh := &Sandesh{
		context:   context.Background(),
		transport: mem_buffer,
		protocol:  vrouter,
	}

	return sandesh
}

type VrMessage struct {
	sk      *NetlinkSocket
	family  GenlFamily
	sandesh Sandesh
}

const FUEMessage = `Generic netlink family '%s' unavailable; 
the vrouter kernel module is probably not loaded,
try 'modprobe openvswitch'
`

type familyUnavailableError struct {
	family string
}

func (fue familyUnavailableError) Error() string {
	return fmt.Sprintf(FUEMessage, fue.family)
}

func IsKernelLacksVrouterError(err error) bool {
	_, ok := err.(familyUnavailableError)
	return ok
}

func lookupFamily(sk *NetlinkSocket, name string) (GenlFamily, error) {
	family, err := sk.LookupGenlFamily(name)
	if err == nil {
		return family, err
	}
	return GenlFamily{}, err
}

func NewVrMessage() (*VrMessage, error) {
	sandesh := newSandesh()
	sk, err := OpenNetlinkSocket(syscall.NETLINK_GENERIC)
	if err != nil {
		return nil, err
	}

	vr_msg := VrMessage{sk: sk, sandesh: *sandesh}
	family, err := lookupFamily(sk, "vrouter")
	if err != nil {
		return nil, err
	}

	vr_msg.family = family
	return &vr_msg, nil
}

func (vr_msg *VrMessage) Reopen() error {
	sk, err := OpenNetlinkSocket(syscall.NETLINK_GENERIC)
	if err != nil {
		return err
	}

	vr_msg.sk = sk
	return nil
}

func (vr_msg *VrMessage) GetMcGroup(name string) (uint32, error) {
	if mcGroup, ok := vr_msg.family.mcGroups[name]; ok {
		return mcGroup, nil
	}

	errmsg := fmt.Errorf("no genl MC group %s in vrouter family", name)
	return 0, errmsg
}

func (vr_msg *VrMessage) Close() error {
	return vr_msg.sk.Close()
}

type nlResponse struct {
	genlhdr *GenlMsghdr
	data    []byte
}

func (vr_msg *VrMessage) handleNlResponse(resp *NlMsgParser) (*nlResponse, error) {
	nl_resp := &nlResponse{}

	if _, err := resp.ExpectNlMsghdr(vr_msg.family.id); err != nil {
		return nil, err
	}

	genlhdr, err := resp.CheckGenlMsghdr(NL_ATTR_VR_MESSAGE_PROTOCOL, -1)
	if err != nil {
		return nil, err
	}

	attrs, err := ParseNestedAttrs(resp.data)
	if err != nil {
		return nil, err
	}

	pos, err := resp.AlignAdvance(syscall.NLA_ALIGNTO, NETLINK_RESPONSE_HEADER_LEN)
	if err != nil {
		return nil, err
	}

	nl_resp.genlhdr = genlhdr
	nl_resp.data = attrs[0][pos:]
	return nl_resp, nil
}

func (vr_msg *VrMessage) nlTransRequest(vr_req vr.Sandesh) (*nlResponse, error) {
	if err := vr_req.Write(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		return nil, errors.New("failed to encode request into binary")
	}

	req_b := vr_msg.sandesh.transport.Bytes()
	req := NewNlMsgBuilder(RequestFlags, vr_msg.family.id)
	req.PutGenlMsghdr(NL_ATTR_VR_MESSAGE_PROTOCOL, 0)
	req.PutSliceAttr(SANDESH_REQUEST, req_b)

	resp, err := vr_msg.sk.Request(req)
	if err != nil {
		return nil, err
	}

	return vr_msg.handleNlResponse(resp)
}

func (vr_msg *VrMessage) sync(args vr.Sandesh) (*vr_raw.VrResponse, error) {
	resp, err := vr_msg.nlTransRequest(args)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(resp.data)
	vr_msg.sandesh.transport.Buffer = buf

	vr_resp := vr_raw.NewVrResponse()
	if err := vr_resp.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse vr_response: %v", err)
		return nil, errmsg
	}

	return vr_resp, nil
}

func (vr_msg *VrMessage) nlTransMultiRequest(vr_req vr.Sandesh) ([]*nlResponse, error) {
	nl_resps := []*nlResponse{}

	if err := vr_req.Write(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		return nl_resps, errors.New("failed to encode request into binary")
	}

	req_b := vr_msg.sandesh.transport.Bytes()
	req := NewNlMsgBuilder(RequestFlags, vr_msg.family.id)
	req.PutGenlMsghdr(NL_ATTR_VR_MESSAGE_PROTOCOL, 0)
	req.PutSliceAttr(SANDESH_REQUEST, req_b)

	consumer := func(resp *NlMsgParser) error {
		nl_resp, _ := vr_msg.handleNlResponse(resp)
		nl_resps = append(nl_resps, nl_resp)
		return nil
	}

	vr_msg.sk.RequestMulti(req, consumer)
	return nl_resps, nil
}

func (vr_msg *VrMessage) syncMultipart(args vr.Sandesh) (*vr_raw.VrResponse, []*nlResponse, error) {
	nl_resps, err := vr_msg.nlTransMultiRequest(args)
	if err != nil {
		return nil, []*nlResponse{}, err
	}

	buf := bytes.NewBuffer((*nl_resps[0]).data)
	vr_msg.sandesh.transport.Buffer = buf

	vr_resp := vr_raw.NewVrResponse()
	if err := vr_resp.Read(vr_msg.sandesh.context, vr_msg.sandesh.protocol); err != nil {
		errmsg := fmt.Errorf("failed to parse vr_response: %v", err)
		return nil, []*nlResponse{}, errmsg
	}

	return vr_resp, nl_resps[1:], nil
}
