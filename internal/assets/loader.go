package assets

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"net/http"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"

	rootassets "github.com/drivakos/arrows-of-olympus/assets"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procmusic"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procsfx"
	"github.com/drivakos/arrows-of-olympus/internal/sprite"
)

// Bundle holds all validated, pre-sliced GPU images, sprite sheets, fonts, and 4-byte aligned PCM audio.
type Bundle struct {
	Background  *ebiten.Image
	FontRegular *text.GoTextFaceSource
	FontBold    *text.GoTextFaceSource

	HeroSheet              *sprite.GridSpriteSheet
	SkeletonSheet          *sprite.GridSpriteSheet
	WarriorUnshieldedSheet *sprite.GridSpriteSheet
	WarriorShieldedSheet   *sprite.GridSpriteSheet
	BossZeusSheet          *sprite.GridSpriteSheet
	BossCronusSheet        *sprite.GridSpriteSheet
	ChestSheet             *sprite.GridSpriteSheet
	PowerUpRelicSheet      *sprite.GridSpriteSheet
	XPGemSheet             *sprite.GridSpriteSheet
	PillarSheet            *sprite.GridSpriteSheet

	AudioCtx *audio.Context
	BGMPlayer *audio.Player
	BGMPMC    []byte
	SFX       map[string][]byte
	lastSFX   map[string]time.Time
}

// AlignPCM masks a stereo 16-bit PCM byte slice length to a 4-byte boundary.
func AlignPCM(pcm []byte) []byte {
	return pcm[:len(pcm)&^3]
}

// LoadAll loads, validates, and pre-slices all embedded game assets.
func LoadAll(audioCtx *audio.Context) (*Bundle, error) {
	b := &Bundle{
		AudioCtx: audioCtx,
		SFX:      make(map[string][]byte),
		lastSFX:  make(map[string]time.Time),
	}

	regSrc, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, fmt.Errorf("assets: load regular font: %w", err)
	}
	boldSrc, err := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		return nil, fmt.Errorf("assets: load bold font: %w", err)
	}
	b.FontRegular = regSrc
	b.FontBold = boldSrc

	bgImg, err := loadPNG("images/olympus_bg.png")
	if err != nil {
		return nil, err
	}
	b.Background = bgImg

	if b.HeroSheet, err = loadGridSheet("sprites/hero_greek.png", 64, 64); err != nil {
		return nil, err
	}
	if b.SkeletonSheet, err = loadGridSheet("sprites/mob_skeleton.png", 64, 64); err != nil {
		return nil, err
	}
	if b.WarriorUnshieldedSheet, err = loadGridSheet("sprites/mob_warrior_unshielded.png", 64, 64); err != nil {
		return nil, err
	}
	if b.WarriorShieldedSheet, err = loadGridSheet("sprites/mob_warrior_shielded.png", 64, 64); err != nil {
		return nil, err
	}
	if b.BossZeusSheet, err = loadGridSheet("sprites/boss_zeus_champion.png", 96, 96); err != nil {
		return nil, err
	}
	if b.BossCronusSheet, err = loadGridSheet("sprites/boss_cronus.png", 112, 112); err != nil {
		return nil, err
	}
	if b.ChestSheet, err = loadGridSheet("sprites/prop_chest.png", 48, 48); err != nil {
		return nil, err
	}
	if b.PowerUpRelicSheet, err = loadGridSheet("sprites/powerup_relic.png", 36, 36); err != nil {
		return nil, err
	}
	if b.XPGemSheet, err = loadGridSheet("sprites/pickup_xp.png", 24, 24); err != nil {
		return nil, err
	}
	if b.PillarSheet, err = loadGridSheet("sprites/prop_pillar.png", 64, 64); err != nil {
		return nil, err
	}

	// Compile 12 Gods of Olympus BGM score to 4-byte aligned PCM
	bgmBytes, err := rootassets.FS.ReadFile("audio/olympus_12_gods.music")
	if err != nil {
		return nil, fmt.Errorf("assets: read olympus_12_gods.music: %w", err)
	}
	bgmPCM, err := procmusic.CompileDSL(string(bgmBytes))
	if err != nil {
		return nil, fmt.Errorf("assets: compile olympus_12_gods.music: %w", err)
	}
	b.BGMPMC = AlignPCM(bgmPCM)

	if audioCtx != nil && len(b.BGMPMC) >= 4 {
		loop := audio.NewInfiniteLoop(bytes.NewReader(b.BGMPMC), int64(len(b.BGMPMC)&^3))
		player, err := audioCtx.NewPlayer(loop)
		if err != nil {
			return nil, fmt.Errorf("assets: create BGM player: %w", err)
		}
		player.SetVolume(0.55)
		b.BGMPlayer = player
	}

	sfxFiles := map[string]string{
		"bow_shoot":      "audio/bow_shoot.sfx",
		"sword_slash":    "audio/sword_slash.sfx",
		"arrow_hit":      "audio/arrow_hit.sfx",
		"shield_clang":   "audio/shield_clang.sfx",
		"chest_open":     "audio/chest_open.sfx",
		"powerup_pickup": "audio/powerup_pickup.sfx",
		"xp_ambrosia":    "audio/xp_ambrosia.sfx",
		"boon_levelup":   "audio/boon_levelup.sfx",
		"thunder_boss":   "audio/thunder_boss.sfx",
		"hero_hurt":      "audio/hero_hurt.sfx",
	}
	for cue, path := range sfxFiles {
		raw, err := rootassets.FS.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("assets: read %s: %w", path, err)
		}
		pcm, err := procsfx.CompileDSL(string(raw))
		if err != nil {
			return nil, fmt.Errorf("assets: compile %s: %w", path, err)
		}
		b.SFX[cue] = AlignPCM(pcm)
	}

	return b, nil
}

// PlaySFX plays a pre-compiled SFX cue with rate-limiting cooldown and 4-byte PCM alignment.
func (b *Bundle) PlaySFX(cue string, cooldown time.Duration) {
	if b == nil || b.AudioCtx == nil {
		return
	}
	pcm, ok := b.SFX[cue]
	if !ok || len(pcm) < 4 {
		return
	}
	now := time.Now()
	if last, exists := b.lastSFX[cue]; exists && now.Sub(last) < cooldown {
		return
	}
	b.lastSFX[cue] = now
	p := b.AudioCtx.NewPlayerFromBytes(pcm[:len(pcm)&^3])
	p.SetVolume(0.75)
	p.Play()
}

func loadPNG(relPath string) (*ebiten.Image, error) {
	data, err := rootassets.FS.ReadFile(relPath)
	if err != nil {
		return nil, fmt.Errorf("assets: read %s: %w", relPath, err)
	}
	mime := http.DetectContentType(data)
	if mime != "image/png" {
		return nil, fmt.Errorf("assets: %s has invalid MIME type %q (expected image/png)", relPath, mime)
	}
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("assets: decode %s: %w", relPath, err)
	}
	return ebiten.NewImageFromImage(decoded), nil
}

func loadGridSheet(relPath string, fw, fh int) (*sprite.GridSpriteSheet, error) {
	img, err := loadPNG(relPath)
	if err != nil {
		return nil, err
	}
	sheet, err := sprite.NewGridSpriteSheet(img, fw, fh)
	if err != nil {
		return nil, fmt.Errorf("assets: slice grid %s (%dx%d): %w", relPath, fw, fh, err)
	}
	return sheet, nil
}

// NewHeroAnim creates a configured AnimationController for the Greek Hero.
func (b *Bundle) NewHeroAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.HeroSheet)
	_ = ctrl.AddTag("idle", 0, 3, 120*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("run", 4, 7, 90*time.Millisecond, sprite.LoopForward)
	_ = ctrl.AddTag("bow_attack", 8, 11, 75*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("sword_attack", 12, 15, 75*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("hit", 16, 17, 90*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("death", 18, 21, 110*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.Play("idle")
	return ctrl
}

// NewSkeletonAnim creates a configured AnimationController for Skeleton Spartoi.
func (b *Bundle) NewSkeletonAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.SkeletonSheet)
	_ = ctrl.AddTag("idle", 0, 3, 120*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("run", 4, 7, 85*time.Millisecond, sprite.LoopForward)
	_ = ctrl.AddTag("attack", 8, 11, 80*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("hit", 12, 13, 90*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("death", 14, 17, 100*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.Play("run")
	return ctrl
}

// NewWarriorUnshieldedAnim creates an AnimationController for Unshielded Greek Warriors.
func (b *Bundle) NewWarriorUnshieldedAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.WarriorUnshieldedSheet)
	_ = ctrl.AddTag("idle", 0, 3, 120*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("run", 4, 7, 90*time.Millisecond, sprite.LoopForward)
	_ = ctrl.AddTag("attack", 8, 11, 80*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("hit", 12, 13, 90*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("death", 14, 17, 110*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.Play("run")
	return ctrl
}

// NewWarriorShieldedAnim creates an AnimationController for Shielded Hoplite Warriors.
func (b *Bundle) NewWarriorShieldedAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.WarriorShieldedSheet)
	_ = ctrl.AddTag("idle", 0, 3, 120*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("run", 4, 7, 95*time.Millisecond, sprite.LoopForward)
	_ = ctrl.AddTag("shield_guard", 8, 11, 100*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("attack", 12, 15, 85*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("hit", 16, 17, 90*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("death", 18, 21, 110*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.Play("run")
	return ctrl
}

// NewBossZeusAnim creates an AnimationController for Zeus's Champion.
func (b *Bundle) NewBossZeusAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.BossZeusSheet)
	_ = ctrl.AddTag("idle", 0, 3, 120*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("move", 4, 7, 95*time.Millisecond, sprite.LoopForward)
	_ = ctrl.AddTag("attack", 8, 11, 85*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("cast", 12, 15, 100*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("death", 16, 19, 120*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.Play("move")
	return ctrl
}

// NewBossCronusAnim creates an AnimationController for Titan Cronus.
func (b *Bundle) NewBossCronusAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.BossCronusSheet)
	_ = ctrl.AddTag("idle", 0, 3, 120*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("move", 4, 7, 100*time.Millisecond, sprite.LoopForward)
	_ = ctrl.AddTag("cleave", 8, 11, 85*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("summon", 12, 15, 110*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.AddTag("death", 16, 19, 130*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.Play("move")
	return ctrl
}

// NewChestAnim creates an AnimationController for Golden Olympus Chests.
func (b *Bundle) NewChestAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.ChestSheet)
	_ = ctrl.AddTag("closed", 0, 3, 150*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.AddTag("open", 4, 7, 100*time.Millisecond, sprite.LoopOnce)
	_ = ctrl.Play("closed")
	return ctrl
}

// NewRelicAnim creates an AnimationController for the unified Divine Power-Up Relic.
func (b *Bundle) NewRelicAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.PowerUpRelicSheet)
	_ = ctrl.AddTag("idle", 0, 3, 110*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.Play("idle")
	return ctrl
}

// NewXPGemAnim creates an AnimationController for Ambrosia XP Gems.
func (b *Bundle) NewXPGemAnim() *sprite.AnimationController {
	ctrl := sprite.NewAnimationController(b.XPGemSheet)
	_ = ctrl.AddTag("idle", 0, 3, 120*time.Millisecond, sprite.LoopPingPong)
	_ = ctrl.Play("idle")
	return ctrl
}
