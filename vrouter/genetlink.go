// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package vrouter

import (
	"fmt"
	"syscall"
)

const GENL_HDRLEN = 0x4
const NETLINK_RESPONSE_HEADER_LEN = syscall.NLA_HDRLEN +
	syscall.NLMSG_HDRLEN +
	GENL_HDRLEN

type GenlFamily struct {
	id       uint16
	mcGroups map[string]uint32
}

func (nlmsg *NlMsgBuilder) PutGenlMsghdr(cmd uint8, version uint8) *GenlMsghdr {
	pos := nlmsg.AlignGrow(syscall.NLMSG_ALIGNTO, SizeofGenlMsghdr)
	res := genlMsghdrAt(nlmsg.buf, pos)
	res.Cmd = cmd
	res.Version = version
	return res
}

func (nlmsg *NlMsgParser) CheckGenlMsghdr(cmd int, fallbackCmd int) (*GenlMsghdr, error) {
	pos, err := nlmsg.AlignAdvance(syscall.NLMSG_ALIGNTO, SizeofGenlMsghdr)
	if err != nil {
		return nil, err
	}

	gh := genlMsghdrAt(nlmsg.data, pos)
	if cmd >= 0 && gh.Cmd != uint8(cmd) && (fallbackCmd < 0 || gh.Cmd != uint8(fallbackCmd)) {
		return nil, fmt.Errorf("generic netlink response has wrong cmd (got %d, expected %d (or fallback: %d))",
			gh.Cmd, cmd, fallbackCmd)
	}

	// Deliberately ignore the version field in the genl header.
	// It's unclear exactly what its meaning is, and how we should
	// handle it.  E.g., if the version is higher than we expect,
	// should we still try to handle the message?  It's unclear,
	// but the fact that ODP bumped the kernel
	// OVS_DATAPATH_VERSION from 1 to 2 while expecting existing
	// userspace to keep working suggests that we should be
	// libreral in what we accept.

	return gh, nil
}

func (s *NetlinkSocket) LookupGenlFamily(name string) (family GenlFamily, err error) {
	req := NewNlMsgBuilder(RequestFlags, GENL_ID_CTRL)

	req.PutGenlMsghdr(CTRL_CMD_GETFAMILY, 0)
	req.PutStringAttr(CTRL_ATTR_FAMILY_NAME, name)

	resp, err := s.Request(req)
	if err != nil {
		return
	}

	_, err = resp.ExpectNlMsghdr(GENL_ID_CTRL)
	if err != nil {
		return
	}

	// For now response command is always CTRL_CMD_NEWFAMILY, though it should have
	// been CTRL_CMD_GETFAMILY. For stability forever, we utilize fallbacking here.
	_, err = resp.CheckGenlMsghdr(CTRL_CMD_GETFAMILY, CTRL_CMD_NEWFAMILY)
	if err != nil {
		return
	}

	attrs, err := resp.TakeAttrs()
	if err != nil {
		return
	}

	family.id, err = attrs.GetUint16(CTRL_ATTR_FAMILY_ID)
	if err != nil {
		return
	}

	mcGroupAttrs, err := attrs.GetNestedAttrs(CTRL_ATTR_MCAST_GROUPS, true)
	if err != nil || mcGroupAttrs == nil {
		return
	}

	family.mcGroups = make(map[string]uint32)
	for _, data := range mcGroupAttrs {
		groupAttrs, err := ParseNestedAttrs(data)
		if err != nil {
			return family, err
		}

		id, err := groupAttrs.GetUint32(CTRL_ATTR_MCAST_GRP_ID)
		if err != nil {
			return family, err
		}

		name, err := groupAttrs.GetString(CTRL_ATTR_MCAST_GRP_NAME)
		if err != nil {
			return family, err
		}

		family.mcGroups[name] = id
	}

	return
}
