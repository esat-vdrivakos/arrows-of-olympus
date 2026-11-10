package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/drivakos/arrows-of-olympus/internal/assets"
	"github.com/drivakos/arrows-of-olympus/internal/audio"
	"github.com/drivakos/arrows-of-olympus/internal/camera"
	"github.com/drivakos/arrows-of-olympus/internal/combat"
	"github.com/drivakos/arrows-of-olympus/internal/entities"
	"github.com/drivakos/arrows-of-olympus/internal/input"
	"github.com/drivakos/arrows-of-olympus/internal/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	VirtualWidth  = 1280
	VirtualHeight = 720
)

// SceneID identifies the active state in the game's finite state machine.
type SceneID int

const (
	SceneTitle SceneID = iota
	SceneWeaponSelect
	SceneGameplay
	SceneBoonSelect
	ScenePaused
	SceneVictory
	SceneGameOver
)

// Game implements ebiten.Game and coordinates the scene state machine, audio, input, and arena.
type Game struct {
	Reg   *assets.Registry
	Audio *audio.Manager
	Cam   *camera.Camera
	Input input.State
	Arena *entities.Arena
	RNG   *rand.Rand

	Scene          SceneID
	SelectedWeapon combat.WeaponType
	BoonChoices    [3]combat.BoonDef
	HoveredCard    int

	TitleIdleTimer float64
	PulseTimer     float64

	// Preview controllers for Title & Weapon Select scenes
	PreviewHeroBow   *sprite.AnimationController
	PreviewHeroSword *sprite.AnimationController
	PreviewZeus      *sprite.AnimationController
	PreviewCronus    *sprite.AnimationController

	uiOp ebiten.DrawImageOptions

	bgmCh chan []byte
}

// New initializes the complete Arrows of Olympus game instance.
func New() (*Game, error) {
	reg, err := assets.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("init registry: %w", err)
	}
	aud, err := audio.NewManager(nil, reg.SFX)
	if err != nil {
		return nil, fmt.Errorf("init audio manager: %w", err)
	}
	cam := camera.New(VirtualWidth, VirtualHeight)
	arena, err := entities.NewArena(reg, aud, cam)
	if err != nil {
		return nil, fmt.Errorf("init arena: %w", err)
	}

	heroBow, _ := reg.NewSpriteController("hero_greek")
	_ = heroBow.Play("bow_attack")
	heroBow.Tags["bow_attack"].LoopMode = sprite.LoopForward

	heroSword, _ := reg.NewSpriteController("hero_greek")
	_ = heroSword.Play("sword_attack")
	heroSword.Tags["sword_attack"].LoopMode = sprite.LoopForward

	zeusCtrl, _ := reg.NewSpriteController("boss_zeus_champion")
	cronusCtrl, _ := reg.NewSpriteController("boss_cronus")

	g := &Game{
		Reg:              reg,
		Audio:            aud,
		Cam:              cam,
		Arena:            arena,
		RNG:              rand.New(rand.NewSource(time.Now().UnixNano())),
		Scene:            SceneTitle,
		SelectedWeapon:   combat.WeaponBow,
		HoveredCard:      0,
		PreviewHeroBow:   heroBow,
		PreviewHeroSword: heroSword,
		PreviewZeus:      zeusCtrl,
		PreviewCronus:    cronusCtrl,
	}
	// Synthesizing the BGM takes seconds (longer in WASM); do it off the game loop.
	g.bgmCh = make(chan []byte, 1)
	go func() {
		pcm, err := assets.CompileBGM()
		if err != nil {
			log.Printf("bgm: %v", err)
			return
		}
		g.bgmCh <- pcm
	}()
	g.Audio.StartBGM()
	return g, nil
}

// Layout returns the fixed 1280x720 16:9 widescreen virtual canvas.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return VirtualWidth, VirtualHeight
}

// Update advances the active scene at 60 TPS.
func (g *Game) Update() error {
	const dt = 1.0 / 60.0
	dtDur := time.Second / 60
	g.PulseTimer += dt
	select {
	case pcm := <-g.bgmCh:
		if err := g.Audio.SetBGM(pcm); err != nil {
			log.Printf("bgm: %v", err)
		}
	default:
	}

	allowMouseSteer := g.Scene == SceneGameplay && !g.Arena.DemoMode
	g.Input.Poll(g.Arena.Hero.X, g.Arena.Hero.Y, allowMouseSteer)

	if g.Input.AnyActivity {
		g.Audio.UnlockOnUserGesture()
	}
	if g.Input.MuteJustPressed {
		g.Audio.ToggleMute()
	}

	switch g.Scene {
	case SceneTitle:
		g.PreviewHeroBow.Update(dtDur)
		g.PreviewHeroSword.Update(dtDur)
		g.PreviewZeus.Update(dtDur)
		g.PreviewCronus.Update(dtDur)

		g.TitleIdleTimer += dt
		if g.Input.ConfirmJustPressed || g.Input.MouseJustClicked {
			g.TitleIdleTimer = 0
			g.Audio.PlaySFX("boon_levelup", 80*time.Millisecond)
			g.Scene = SceneWeaponSelect
			return nil
		}
		// Arcade Attract / Demo Mode after 12 seconds idle on Title Screen
		if g.TitleIdleTimer >= 12.0 {
			g.TitleIdleTimer = 0
			w := combat.WeaponBow
			if g.RNG.Intn(2) == 1 {
				w = combat.WeaponSword
			}
			g.Arena.ResetRun(w, true)
			g.Scene = SceneGameplay
		}

	case SceneWeaponSelect:
		g.PreviewHeroBow.Update(dtDur)
		g.PreviewHeroSword.Update(dtDur)

		if g.Input.LeftJustPressed || g.Input.Key1JustPressed {
			g.SelectedWeapon = combat.WeaponBow
			g.Audio.PlaySFX("xp_ambrosia", 40*time.Millisecond)
		}
		if g.Input.RightJustPressed || g.Input.Key2JustPressed {
			g.SelectedWeapon = combat.WeaponSword
			g.Audio.PlaySFX("xp_ambrosia", 40*time.Millisecond)
		}

		// Mouse hover / click on the 2 Weapon Cards
		mx, my := g.Input.MouseX, g.Input.MouseY
		if my >= 170 && my <= 575 {
			if mx >= 150 && mx <= 610 {
				g.SelectedWeapon = combat.WeaponBow
				if g.Input.MouseJustClicked {
					g.startNewRun(combat.WeaponBow)
					return nil
				}
			} else if mx >= 670 && mx <= 1130 {
				g.SelectedWeapon = combat.WeaponSword
				if g.Input.MouseJustClicked {
					g.startNewRun(combat.WeaponSword)
					return nil
				}
			}
		}

		if g.Input.ConfirmJustPressed {
			g.startNewRun(g.SelectedWeapon)
			return nil
		}
		if g.Input.PauseJustPressed {
			g.Scene = SceneTitle
		}

	case SceneGameplay:
		// If in Attract / Demo Mode, any user input interrupts immediately back to Title Screen
		if g.Arena.DemoMode {
			if g.Input.AnyActivity {
				g.Arena.DemoMode = false
				g.Cam.Reset()
				g.Scene = SceneTitle
				return nil
			}
			g.Arena.Update(dt, 0, 0, false)
			if g.Arena.CronusSlain || g.Arena.RunTime > 45.0 {
				g.Arena.DemoMode = false
				g.Cam.Reset()
				g.Scene = SceneTitle
			}
			return nil
		}

		if g.Input.PauseJustPressed {
			g.Scene = ScenePaused
			return nil
		}

		// Mobile / Mouse HUD Tap Targets
		if g.Input.MouseJustClicked {
			mx, my := g.Input.MouseX, g.Input.MouseY
			// Tap Top-Right HUD panel (930..1264, 14..76) -> Pause
			if mx >= 930 && mx <= 1264 && my >= 14 && my <= 76 {
				g.Scene = ScenePaused
				return nil
			}
			// Tap On-Screen Dodge Button (Bottom-Right 1120..1250, 610..695) -> Dash
			if mx >= 1120 && mx <= 1250 && my >= 610 && my <= 695 {
				g.Input.DashJustPressed = true
			}
		}

		g.Arena.Update(dt, g.Input.MoveX, g.Input.MoveY, g.Input.DashJustPressed)

		if g.Arena.HeroDead {
			g.Scene = SceneGameOver
			return nil
		}
		if g.Arena.CronusSlain {
			g.Scene = SceneVictory
			return nil
		}
		if g.Arena.Hero.PendingLevelUps > 0 {
			g.BoonChoices = combat.DraftBoons(g.RNG)
			g.HoveredCard = 0
			g.Scene = SceneBoonSelect
			return nil
		}

	case SceneBoonSelect:
		if g.Input.LeftJustPressed {
			g.HoveredCard = (g.HoveredCard + 2) % 3
			g.Audio.PlaySFX("xp_ambrosia", 30*time.Millisecond)
		}
		if g.Input.RightJustPressed {
			g.HoveredCard = (g.HoveredCard + 1) % 3
			g.Audio.PlaySFX("xp_ambrosia", 30*time.Millisecond)
		}
		if g.Input.Key1JustPressed {
			g.selectBoonCard(0)
			return nil
		}
		if g.Input.Key2JustPressed {
			g.selectBoonCard(1)
			return nil
		}
		if g.Input.Key3JustPressed {
			g.selectBoonCard(2)
			return nil
		}

		mx, my := g.Input.MouseX, g.Input.MouseY
		if my >= 180 && my <= 560 {
			for idx := 0; idx < 3; idx++ {
				cx := 125.0 + float64(idx)*355.0
				if mx >= cx && mx <= cx+320.0 {
					g.HoveredCard = idx
					if g.Input.MouseJustClicked {
						g.selectBoonCard(idx)
						return nil
					}
				}
			}
		}
		if g.Input.ConfirmJustPressed {
			g.selectBoonCard(g.HoveredCard)
			return nil
		}

	case ScenePaused:
		if g.Input.PauseJustPressed || g.Input.ConfirmJustPressed || g.Input.MouseJustClicked {
			g.Scene = SceneGameplay
			return nil
		}
		if g.Input.Key1JustPressed {
			g.Scene = SceneWeaponSelect
			return nil
		}

	case SceneVictory, SceneGameOver:
		if g.Input.ConfirmJustPressed || g.Input.MouseJustClicked || g.Input.Key1JustPressed {
			g.Audio.PlaySFX("boon_levelup", 80*time.Millisecond)
			g.Scene = SceneWeaponSelect
			return nil
		}
		if g.Input.PauseJustPressed {
			g.Scene = SceneTitle
			return nil
		}
	}

	return nil
}

func (g *Game) startNewRun(w combat.WeaponType) {
	g.SelectedWeapon = w
	if w == combat.WeaponSword {
		g.Audio.PlaySFX("sword_slash", 40*time.Millisecond)
	} else {
		g.Audio.PlaySFX("bow_shoot", 40*time.Millisecond)
	}
	g.Arena.ResetRun(w, false)
	g.Scene = SceneGameplay
}

func (g *Game) selectBoonCard(idx int) {
	if idx < 0 || idx > 2 {
		return
	}
	chosen := g.BoonChoices[idx]
	g.Arena.ApplyBoon(chosen.ID)
	if g.Arena.Hero.PendingLevelUps > 0 {
		g.BoonChoices = combat.DraftBoons(g.RNG)
		g.HoveredCard = 0
	} else {
		g.Scene = SceneGameplay
	}
}
