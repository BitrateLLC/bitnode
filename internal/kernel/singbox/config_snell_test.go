package singbox

import (
	"testing"

	"github.com/BitrateLLC/bitnode/internal/model"
)

func TestBuildSnellV5OmitsUsers(t *testing.T) {
	users := []model.UserSpec{{ID: 1, UUID: "11111111-1111-1111-1111-111111111111"}}

	cfg := buildSnell(M{}, &model.NodeSpec{Version: 5, PSK: "abcdefghijkl", ObfsMode: "none"}, users)

	if cfg["version"] != 5 {
		t.Fatalf("version = %v, want 5", cfg["version"])
	}
	// v5 必须走 sing-box 的单用户 Service：一旦下发 users，inbound 会切到
	// MultiService，要求客户端携带 userkey，Clash/Surge 永远连不上。
	if _, ok := cfg["users"]; ok {
		t.Fatal("v5 inbound must not include users")
	}
	if cfg["obfs_mode"] != "none" {
		t.Fatalf("obfs_mode = %v, want none", cfg["obfs_mode"])
	}
	if _, ok := cfg["mode"]; ok {
		t.Fatal("v5 inbound must not include v6 mode")
	}
}

func TestBuildSnellV6KeepsUsers(t *testing.T) {
	users := []model.UserSpec{
		{ID: 1, UUID: "11111111-1111-1111-1111-111111111111"},
		{ID: 2, UUID: "22222222-2222-2222-2222-222222222222"},
	}

	cfg := buildSnell(M{}, &model.NodeSpec{Version: 6, PSK: "abcdefghijkl", Mode: "default"}, users)

	list, ok := cfg["users"].([]M)
	if !ok {
		t.Fatalf("users type = %T, want []M", cfg["users"])
	}
	if len(list) != 2 {
		t.Fatalf("len(users) = %d, want 2", len(list))
	}
	if list[0]["userkey"] != users[0].UUID {
		t.Fatalf("userkey = %v, want %s", list[0]["userkey"], users[0].UUID)
	}
	if cfg["mode"] != "default" {
		t.Fatalf("mode = %v, want default", cfg["mode"])
	}
	if _, ok := cfg["obfs_mode"]; ok {
		t.Fatal("v6 inbound must not include v5 obfs_mode")
	}
}
