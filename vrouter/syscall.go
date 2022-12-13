// Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vrouter

import "syscall"

// from linux/include/linux/socket.h
const SOL_NETLINK = 270

type GenlMsghdr struct {
	Cmd      uint8
	Version  uint8
	Reserved uint16
}

const SizeofGenlMsghdr = 4

// reserved static generic netlink identifiers:
const (
	GENL_ID_GENERATE  = 0
	GENL_ID_CTRL      = syscall.NLMSG_MIN_TYPE
	GENL_ID_VFS_DQUOT = syscall.NLMSG_MIN_TYPE + 1
	GENL_ID_PMCRAID   = syscall.NLMSG_MIN_TYPE + 2
)

const (
	CTRL_CMD_UNSPEC       = 0
	CTRL_CMD_NEWFAMILY    = 1
	CTRL_CMD_DELFAMILY    = 2
	CTRL_CMD_GETFAMILY    = 3
	CTRL_CMD_NEWOPS       = 4
	CTRL_CMD_DELOPS       = 5
	CTRL_CMD_GETOPS       = 6
	CTRL_CMD_NEWMCAST_GRP = 7
	CTRL_CMD_DELMCAST_GRP = 8
)

const (
	CTRL_ATTR_UNSPEC       = 0
	CTRL_ATTR_FAMILY_ID    = 1
	CTRL_ATTR_FAMILY_NAME  = 2
	CTRL_ATTR_VERSION      = 3
	CTRL_ATTR_HDRSIZE      = 4
	CTRL_ATTR_MAXATTR      = 5
	CTRL_ATTR_OPS          = 6
	CTRL_ATTR_MCAST_GROUPS = 7
)

const (
	CTRL_ATTR_MCAST_GRP_UNSPEC = 0
	CTRL_ATTR_MCAST_GRP_NAME   = 1
	CTRL_ATTR_MCAST_GRP_ID     = 2
)

const (
	NL_ATTR_VR_MESSAGE_PROTOCOL = 1
	SANDESH_REQUEST             = 1
)
