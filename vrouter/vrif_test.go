package vrouter_test

import (
	"testing"

	"github.com/shun159/go-vrouter/vrouter"
	"github.com/shun159/vr"
	"golang.org/x/sys/unix"
)

func checkCloseVrif(vr_msg *vrouter.VrMessage, t *testing.T) {
	err := vr_msg.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func createTapDevice(t *testing.T) *vrouter.Pkt0Device {
	pkt0 := &vrouter.Pkt0Device{TxBufferCount: 100}

	err := pkt0.Init()
	if err != nil {
		t.Fatal(err)
	}

	vif_pkt0_mac := []int8{}
	for _, o := range pkt0.HardwareAddr {
		vif_pkt0_mac = append(vif_pkt0_mac, int8(o))
	}

	return pkt0
}

func destroyTapDevice(pkt0 *vrouter.Pkt0Device, t *testing.T) {
	err := pkt0.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigureVrouter(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	if _, err := vr_msg.UpdateVRouter(
		vrouter.Perfs(1),
		vrouter.Perfr1(1),
	); err != nil {
		t.Fatal(err)
	}
}

func TestResetVrouter(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	if _, err := vr_msg.ResetVRouter(); err != nil {
		t.Fatal(err)
	}
}

func TestDumpVrouterInterface(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	vif_list, err := vr_msg.DumpVif(
		vrouter.VifRid(0),
		vrouter.VifCore(0),
		vrouter.VifMarker(-1),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(vif_list) != 2 {
		t.Fatalf("number of vif at initial state should be 2")
	}
}

func TestGetVrouterInterface(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	vif, err := vr_msg.GetVif(vrouter.VifIdx(4353))
	if err != nil {
		t.Fatal(err)
	}

	if vif.VifrIdx != 4353 {
		t.Fatalf("the vifr_idx should be 4353")
	}
}

func TestResetStatsVrouterInterface(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	if _, err := vr_msg.ResetStatsVif(); err != nil {
		t.Fatal(err)
	}
}

func TestAddVrouterInterface(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	pkt0 := createTapDevice(t)
	vif_pkt0_mac := []int8{}
	for _, o := range pkt0.HardwareAddr {
		vif_pkt0_mac = append(vif_pkt0_mac, int8(o))
	}

	ret, err := vr_msg.AddVif(
		vrouter.VifName("pkt0"),
		vrouter.VifRid(0),
		vrouter.VifType(vr.VIF_TYPE_AGENT),
		vrouter.VifFlags(vr.VIF_FLAG_L3_ENABLED),
		vrouter.VifTransport(vr.VIF_TRANSPORT_SOCKET),
		vrouter.VifOsIdx(int32(pkt0.Index)),
		vrouter.VifMac(vif_pkt0_mac),
	)

	if ret != 0 || err != nil {
		t.Fatal(err)
	}

	destroyTapDevice(pkt0, t)
}

func TestDumpNexthop(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	nh_list, err := vr_msg.DumpNexthop(
		vrouter.NhMarker(-1),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(nh_list) != 1 {
		t.Fatalf("the number of nexthop at initial state should be 1")
	}
}

func TestAddNexthop(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	ret, err := vr_msg.AddNexthop(
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
	)

	if ret != 0 || err != nil {
		t.Fatal(err)
	}
}

func TestVxlanReq(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	if _, err := vr_msg.AddVxlan(
		vrouter.VxlanVnid(1),
		vrouter.VxlanNhid(1),
	); err != nil {
		t.Fatal(err)
	}

	if _, err := vr_msg.GetVxlan(
		vrouter.VxlanVnid(1),
	); err != nil {
		t.Fatal(err)
	}

	if _, err := vr_msg.DumpVxlan(
		vrouter.VxlanVnid(-1),
	); err != nil {
		t.Fatal(err)
	}

	if _, err := vr_msg.DelVxlan(
		vrouter.VxlanVnid(1),
	); err != nil {
		t.Fatal(err)
	}
}

func TestDelNexthop(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	if _, err := vr_msg.DelNexthop(
		vrouter.NhID(1),
	); err != nil {
		t.Fatal(err)
	}
}

func TestAddRoute(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	ret, err := vr_msg.AddRoute(
		vrouter.RouteRid(0),
		vrouter.RouteVrfId(0),
		vrouter.RouteFamily(unix.AF_INET),
		vrouter.RouteMac([]int8{0, 0, 0, 0, 0, 0}),
		vrouter.RoutePrefix([]int8{100, 100, 100, 100}),
		vrouter.RoutePrefixLen(32),
		vrouter.RouteNhId(vr.NH_DISCARD_ID),
	)

	if err != nil {
		t.Fatal(err)
	}

	if ret != 0 {
		t.Fatalf("failed to add route: resp_code should be 0")
	}
}

func TestGetRoute(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	rt, err := vr_msg.GetRoute(
		vrouter.RouteVrfId(0),
		vrouter.RouteFamily(unix.AF_INET),
		vrouter.RouteMac([]int8{0, 0, 0, 0, 0, 0}),
		vrouter.RoutePrefix([]int8{100, 100, 100, 100}),
		vrouter.RoutePrefixLen(32),
	)

	if err != nil {
		t.Fatal(err)
	}

	if rt == nil {
		t.Fatal(rt)
	}
}

func TestDumpRoute(t *testing.T) {
	vr_msg, err := vrouter.NewVrMessage()
	if err != nil {
		t.Fatal(err)
	}

	defer checkCloseVrif(vr_msg, t)

	rt, err := vr_msg.DumpRoute(
		vrouter.RouteVrfId(0),
		vrouter.RouteFamily(unix.AF_INET),
		vrouter.RouteMac([]int8{0, 0, 0, 0, 0, 0}),
		vrouter.RoutePrefix([]int8{0, 0, 0, 0}),
		vrouter.RouteMarker([]int8{100, 100, 100, 0}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(rt) < 1 {
		t.Fatal(rt)
	}
}
