package assets

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"time"

	"github.com/drivakos/arrows-of-olympus/assets"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procmusic"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procsfx"
	"github.com/drivakos/arrows-of-olympus/internal/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

// SpriteTemplate holds a compiled, GPU-uploaded sprite sheet and its shared animation tags
// so multiple runtime entities can instantiate zero-allocation AnimationControllers.
type SpriteTemplate struct {
	Name   string
	FrameW int
	FrameH int
	Sheet  *sprite.GridSpriteSheet
	Tags   map[string]*sprite.AnimationTag
	First  string
}

// NewController creates a lightweight AnimationController sharing the template's GPU atlas and tag map.
func (st *SpriteTemplate) NewController() *sprite.AnimationController {
	ctrl := &sprite.AnimationController{
		Sheet:       st.Sheet,
		Tags:        st.Tags,
		PingPongDir: 1,
	}
	if st.First != "" {
		_ = ctrl.Play(st.First)
	}
	return ctrl
}

// Registry stores all pre-loaded, validated game assets in memory.
type Registry struct {
	OlympusBG  *ebiten.Image
	WhitePixel *ebiten.Image

	Sprites map[string]*SpriteTemplate
	SFX     map[string][]byte

	FontBoldSource    *text.GoTextFaceSource
	FontRegularSource *text.GoTextFaceSource

	FaceTitle    *text.GoTextFace
	FaceSubtitle *text.GoTextFace
	FaceHeading  *text.GoTextFace
	FaceBody     *text.GoTextFace
	FaceSmall    *text.GoTextFace
	FaceTiny     *text.GoTextFace
	FaceCombat   *text.GoTextFace
}

var requiredSprites = []string{
	"hero_greek",
	"mob_skeleton",
	"mob_skeleton_archer",
	"mob_warrior_unshielded",
	"mob_warrior_shielded",
	"boss_zeus_champion",
	"boss_cronus",
	"pickup_xp",
	"powerup_relic",
	"prop_chest",
	"prop_pillar",
	"proj_arrow",
	"fx_slash_wave",
}

var requiredSFX = []string{
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

// LoadRegistry loads, compiles, and validates all embedded game assets on boot.
// Fails fast with descriptive diagnostics if any asset is missing or invalid.
func LoadRegistry() (*Registry, error) {
	reg := &Registry{
		Sprites: make(map[string]*SpriteTemplate, len(requiredSprites)),
		SFX:     make(map[string][]byte, len(requiredSFX)),
	}

	// 1. Pre-allocate 1x1 white pixel for vector/rect drawing
	wp := ebiten.NewImage(1, 1)
	wp.Fill(color.RGBA{255, 255, 255, 255})
	reg.WhitePixel = wp

	// 2. Load & validate typography (ebiten/v2/text/v2)
	boldSrc, err := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		return nil, fmt.Errorf("boot: load bold font source: %w", err)
	}
	regSrc, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, fmt.Errorf("boot: load regular font source: %w", err)
	}
	reg.FontBoldSource = boldSrc
	reg.FontRegularSource = regSrc

	reg.FaceTitle = &text.GoTextFace{Source: boldSrc, Size: 46}
	reg.FaceSubtitle = &text.GoTextFace{Source: boldSrc, Size: 24}
	reg.FaceHeading = &text.GoTextFace{Source: boldSrc, Size: 18}
	reg.FaceBody = &text.GoTextFace{Source: regSrc, Size: 15}
	reg.FaceSmall = &text.GoTextFace{Source: boldSrc, Size: 13}
	reg.FaceTiny = &text.GoTextFace{Source: regSrc, Size: 12}
	reg.FaceCombat = &text.GoTextFace{Source: boldSrc, Size: 16}

	// 3. Load & binary-validate Olympus background PNG
	bgBytes, err := assets.FS.ReadFile("images/olympus_bg.png")
	if err != nil {
		return nil, fmt.Errorf("boot: read images/olympus_bg.png: %w", err)
	}
	mimeType := http.DetectContentType(bgBytes)
	if mimeType != "image/png" {
		return nil, fmt.Errorf("boot: invalid MIME type %q for images/olympus_bg.png (expected image/png)", mimeType)
	}
	bgDecoded, err := png.Decode(bytes.NewReader(bgBytes))
	if err != nil {
		return nil, fmt.Errorf("boot: decode images/olympus_bg.png: %w", err)
	}
	reg.OlympusBG = ebiten.NewImageFromImage(bgDecoded)

	// 4. Compile & validate all .sprite DSL files into GPU templates
	for _, name := range requiredSprites {
		path := fmt.Sprintf("sprites/%s.sprite", name)
		raw, err := assets.FS.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("boot: read sprite %q: %w", path, err)
		}
		cs, err := sprite.CompileSpriteDSL(string(raw))
		if err != nil {
			return nil, fmt.Errorf("boot: compile sprite %q: %w", path, err)
		}
		if err := sprite.ValidatePremultipliedRGBA(cs.SheetRGBA); err != nil {
			return nil, fmt.Errorf("boot: premultiplied alpha check failed for %q: %w", path, err)
		}
		ebitenSheet := ebiten.NewImageFromImage(cs.SheetRGBA)
		grid, err := sprite.NewGridSpriteSheet(ebitenSheet, cs.FrameW, cs.FrameH)
		if err != nil {
			return nil, fmt.Errorf("boot: slice grid sheet for %q: %w", path, err)
		}
		tags := make(map[string]*sprite.AnimationTag, len(cs.Tags))
		firstTag := ""
		for i, t := range cs.Tags {
			if i == 0 {
				firstTag = t.Name
			}
			mode, err := sprite.ParseLoopMode(t.Loop)
			if err != nil {
				return nil, fmt.Errorf("boot: sprite %q tag %q loop mode: %w", path, t.Name, err)
			}
			dur := time.Duration(t.DurationMs) * time.Millisecond
			if dur <= 0 {
				dur = 100 * time.Millisecond
			}
			loopStart := t.LoopFromFrame
			if loopStart < t.FromFrame || loopStart > t.ToFrame {
				loopStart = t.FromFrame
			}
			tags[t.Name] = &sprite.AnimationTag{
				Name:           t.Name,
				StartFrame:     t.FromFrame,
				LoopStartFrame: loopStart,
				EndFrame:       t.ToFrame,
				FrameDuration:  dur,
				LoopMode:       mode,
			}
		}
		reg.Sprites[name] = &SpriteTemplate{
			Name:   name,
			FrameW: cs.FrameW,
			FrameH: cs.FrameH,
			Sheet:  grid,
			Tags:   tags,
			First:  firstTag,
		}
	}

	// 5. Compile & 4-byte align procedural SFX (.sfx)
	for _, sfxName := range requiredSFX {
		path := fmt.Sprintf("audio/%s.sfx", sfxName)
		raw, err := assets.FS.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("boot: read sfx %q: %w", path, err)
		}
		pcm, err := procsfx.CompileDSL(string(raw))
		if err != nil {
			return nil, fmt.Errorf("boot: compile sfx %q: %w", path, err)
		}
		alignedLen := len(pcm) &^ 3
		if alignedLen == 0 {
			return nil, fmt.Errorf("boot: compiled sfx %q produced empty PCM buffer", path)
		}
		reg.SFX[sfxName] = pcm[:alignedLen]
	}

	log.Printf("boot: loaded %d sprites and %d SFX cues", len(reg.Sprites), len(reg.SFX))
	return reg, nil
}

// CompileBGM synthesizes olympus_12_gods.music into 4-byte aligned stereo PCM.
// It is slow (seconds in WASM), so callers should run it off the game loop.
func CompileBGM() ([]byte, error) {
	raw, err := assets.FS.ReadFile("audio/olympus_12_gods.music")
	if err != nil {
		return nil, fmt.Errorf("bgm: read audio/olympus_12_gods.music: %w", err)
	}
	pcm, err := procmusic.CompileDSL(string(raw))
	if err != nil {
		return nil, fmt.Errorf("bgm: compile audio/olympus_12_gods.music: %w", err)
	}
	aligned := len(pcm) &^ 3
	if aligned == 0 {
		return nil, fmt.Errorf("bgm: compiled BGM produced empty PCM buffer")
	}
	return pcm[:aligned], nil
}

// NewSpriteController returns a new AnimationController for the named sprite template.
func (r *Registry) NewSpriteController(name string) (*sprite.AnimationController, error) {
	tpl, ok := r.Sprites[name]
	if !ok {
		return nil, fmt.Errorf("sprite template %q not found in registry", name)
	}
	return tpl.NewController(), nil
}
