package game

import (
	"bytes"
	"image/png"
	"math/rand"
	"net/http"
	"testing"

	"github.com/drivakos/arrows-of-olympus/assets"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procmusic"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procsfx"
	"github.com/drivakos/arrows-of-olympus/internal/camera"
	"github.com/drivakos/arrows-of-olympus/internal/combat"
	"github.com/drivakos/arrows-of-olympus/internal/sprite"
)

func TestEmbeddedAssetsAndAlignment(t *testing.T) {
	// 1. Validate Olympus background PNG
	bgBytes, err := assets.FS.ReadFile("images/olympus_bg.png")
	if err != nil {
		t.Fatalf("read images/olympus_bg.png: %v", err)
	}
	if got := http.DetectContentType(bgBytes); got != "image/png" {
		t.Fatalf("olympus_bg.png MIME = %q, want image/png", got)
	}
	if _, err := png.Decode(bytes.NewReader(bgBytes)); err != nil {
		t.Fatalf("decode olympus_bg.png: %v", err)
	}

	// 2. Validate all .sprite DSL files compile and obey premultiplied alpha
	sprites := []string{
		"hero_greek",
		"mob_skeleton",
		"mob_warrior_unshielded",
		"mob_warrior_shielded",
		"boss_zeus_champion",
		"boss_cronus",
		"pickup_xp",
		"powerup_relic",
		"prop_chest",
		"prop_pillar",
	}
	for _, name := range sprites {
		raw, err := assets.FS.ReadFile("sprites/" + name + ".sprite")
		if err != nil {
			t.Fatalf("read sprite %s: %v", name, err)
		}
		cs, err := sprite.CompileSpriteDSL(string(raw))
		if err != nil {
			t.Fatalf("compile sprite %s: %v", name, err)
		}
		if err := sprite.ValidatePremultipliedRGBA(cs.SheetRGBA); err != nil {
			t.Fatalf("premultiplied RGBA check failed for %s: %v", name, err)
		}
	}

	// 3. Validate all .sfx cues compile and satisfy 4-byte stereo PCM alignment
	sfxs := []string{
		"bow_shoot",
		"sword_slash",
		"arrow_hit",
		"shield_clang",
		"chest_open",
		"powerup_pickup",
		"xp_ambrosia",
		"boon_levelup",
		"thunder_boss",
		"hero_hurt",
	}
	for _, sfx := range sfxs {
		raw, err := assets.FS.ReadFile("audio/" + sfx + ".sfx")
		if err != nil {
			t.Fatalf("read sfx %s: %v", sfx, err)
		}
		pcm, err := procsfx.CompileDSL(string(raw))
		if err != nil {
			t.Fatalf("compile sfx %s: %v", sfx, err)
		}
		aligned := len(pcm) &^ 3
		if aligned == 0 || aligned%4 != 0 {
			t.Fatalf("sfx %s invalid aligned length: %d", sfx, aligned)
		}
	}

	// 4. Validate olympus_12_gods.music compiles and satisfies 4-byte alignment
	bgmRaw, err := assets.FS.ReadFile("audio/olympus_12_gods.music")
	if err != nil {
		t.Fatalf("read olympus_12_gods.music: %v", err)
	}
	bgmPCM, err := procmusic.CompileDSL(string(bgmRaw))
	if err != nil {
		t.Fatalf("compile olympus_12_gods.music: %v", err)
	}
	if len(bgmPCM)&^3 == 0 {
		t.Fatalf("compiled BGM is empty")
	}
}

func TestCombatBoonsAndCamera(t *testing.T) {
	var boons [combat.BoonCount]int
	boons[combat.BoonArtemis] = 2
	boons[combat.BoonApollo] = 1
	boons[combat.BoonAthena] = 1

	stBow := combat.ComputeStats(combat.WeaponBow, boons, 0.15, 0.15, 0.12, true, true, true)
	if stBow.FrontProjectiles != 3 {
		t.Fatalf("expected 3 front projectiles with 2x Artemis, got %d", stBow.FrontProjectiles)
	}
	if stBow.DiagonalPairs != 1 {
		t.Fatalf("expected 1 diagonal pair with 1x Apollo, got %d", stBow.DiagonalPairs)
	}
	if stBow.OrbitShields != 2 {
		t.Fatalf("expected 2 orbit shields with 1x Athena, got %d", stBow.OrbitShields)
	}

	rng := rand.New(rand.NewSource(99))
	draft := combat.DraftBoons(rng)
	if draft[0].ID == draft[1].ID || draft[1].ID == draft[2].ID || draft[0].ID == draft[2].ID {
		t.Fatalf("DraftBoons returned duplicate Boons: %+v", draft)
	}

	cam := camera.New(1280, 720)
	cam.AddTrauma(0.5)
	cam.TriggerHitstop(0.1)
	if !cam.Update(0.05) {
		t.Fatalf("expected hitstop to return true while active")
	}
	if cam.Update(0.06) {
		t.Fatalf("expected hitstop to expire after 0.11s")
	}
}
