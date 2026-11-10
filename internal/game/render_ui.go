package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/drivakos/arrows-of-olympus/internal/combat"
	"github.com/drivakos/arrows-of-olympus/internal/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Draw renders the current scene and overlays onto the 1280x720 virtual canvas.
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.Scene {
	case SceneTitle:
		g.drawTitleScene(screen)
	case SceneWeaponSelect:
		g.drawWeaponSelectScene(screen)
	case SceneGameplay:
		g.Arena.Draw(screen)
		g.drawHUD(screen)
		if g.Cam.FlashTimer > 0 {
			alpha := uint8(math.Min(120, g.Cam.FlashTimer*380))
			g.drawRect(screen, 0, 0, VirtualWidth, VirtualHeight, color.RGBA{alpha, alpha, uint8(float64(alpha) * 0.75), alpha})
		}
	case SceneBoonSelect:
		g.Arena.Draw(screen)
		g.drawHUD(screen)
		g.drawBoonSelectOverlay(screen)
	case ScenePaused:
		g.Arena.Draw(screen)
		g.drawHUD(screen)
		g.drawPauseOverlay(screen)
	case SceneVictory:
		g.Arena.Draw(screen)
		g.drawHUD(screen)
		g.drawEndScreenOverlay(screen, true)
	case SceneGameOver:
		g.Arena.Draw(screen)
		g.drawHUD(screen)
		g.drawEndScreenOverlay(screen, false)
	}
}

func (g *Game) drawTitleScene(screen *ebiten.Image) {
	// Background
	g.uiOp.GeoM.Reset()
	g.uiOp.ColorScale.Reset()
	b := g.Reg.OlympusBG.Bounds()
	if b.Dx() > 0 && b.Dy() > 0 {
		g.uiOp.GeoM.Scale(VirtualWidth/float64(b.Dx()), VirtualHeight/float64(b.Dy()))
	}
	screen.DrawImage(g.Reg.OlympusBG, &g.uiOp)

	// Dark marble vignette panel
	g.drawPanel(screen, 170, 75, 940, 570, color.RGBA{12, 14, 26, 228}, color.RGBA{235, 190, 65, 255})

	// Title & Mythic Subtitle
	g.drawTextCentered(screen, "ARROWS OF OLYMPUS", g.Reg.FaceTitle, 640, 115, color.RGBA{255, 215, 75, 255})
	g.drawTextCentered(screen, "ASCEND THE PANTHEON  •  12 OLYMPIAN BOONS  •  OVERTHROW TITAN CRONUS", g.Reg.FaceHeading, 640, 172, color.RGBA{130, 215, 255, 255})

	// Showcase Sprites (Hero Bow, Hero Sword, Zeus's Champion, Titan Cronus)
	g.drawPanel(screen, 220, 215, 840, 195, color.RGBA{20, 24, 42, 235}, color.RGBA{110, 150, 210, 220})

	g.drawUISprite(screen, g.PreviewHeroBow, 325, 295, 1.55)
	g.drawTextCentered(screen, "BOW OF APOLLO", g.Reg.FaceSmall, 325, 365, color.RGBA{255, 210, 75, 255})
	g.drawTextCentered(screen, "Ranged Auto-Shooter", g.Reg.FaceTiny, 325, 384, color.RGBA{200, 210, 230, 255})

	g.drawUISprite(screen, g.PreviewHeroSword, 525, 295, 1.55)
	g.drawTextCentered(screen, "SWORD OF ARES", g.Reg.FaceSmall, 525, 365, color.RGBA{255, 95, 95, 255})
	g.drawTextCentered(screen, "Melee Cleave & Waves", g.Reg.FaceTiny, 525, 384, color.RGBA{200, 210, 230, 255})

	g.drawUISprite(screen, g.PreviewZeus, 740, 290, 1.15)
	g.drawTextCentered(screen, "ZEUS'S CHAMPION", g.Reg.FaceSmall, 740, 365, color.RGBA{120, 210, 255, 255})
	g.drawTextCentered(screen, "Wave II Mid-Boss", g.Reg.FaceTiny, 740, 384, color.RGBA{200, 210, 230, 255})

	g.drawUISprite(screen, g.PreviewCronus, 950, 286, 1.08)
	g.drawTextCentered(screen, "TITAN CRONUS", g.Reg.FaceSmall, 950, 365, color.RGBA{255, 110, 140, 255})
	g.drawTextCentered(screen, "Final Lord of Time", g.Reg.FaceTiny, 950, 384, color.RGBA{200, 210, 230, 255})

	// Feature summary & Controls
	g.drawTextCentered(screen, "CONTROLS: WASD / Arrows / Mouse Hold to Move  |  SPACE / SHIFT to Dodge Roll  |  Auto-Aim & Auto-Attack", g.Reg.FaceBody, 640, 435, color.RGBA{235, 240, 255, 255})
	g.drawTextCentered(screen, "Break Golden Olympus Chests & Slay Spartoi for Attack Speed, Damage & Haste Relics!", g.Reg.FaceBody, 640, 465, color.RGBA{255, 215, 115, 255})

	// Pulsing Call-To-Action Button
	pulse := 0.75 + 0.25*math.Sin(g.PulseTimer*5.0)
	btnCol := color.RGBA{uint8(220 * pulse), uint8(165 * pulse), 35, 255}
	g.drawPanel(screen, 420, 515, 440, 58, color.RGBA{38, 30, 14, 245}, btnCol)
	g.drawTextCentered(screen, "PRESS ENTER / SPACE OR CLICK TO ASCEND", g.Reg.FaceHeading, 640, 535, color.RGBA{255, 235, 130, 255})

	g.drawTextCentered(screen, "Press [M] Mute Audio  |  [F11] Fullscreen  |  Attract Demo starts automatically on idle", g.Reg.FaceTiny, 640, 608, color.RGBA{165, 175, 200, 255})
}

func (g *Game) drawWeaponSelectScene(screen *ebiten.Image) {
	// Background
	g.uiOp.GeoM.Reset()
	g.uiOp.ColorScale.Reset()
	b := g.Reg.OlympusBG.Bounds()
	if b.Dx() > 0 && b.Dy() > 0 {
		g.uiOp.GeoM.Scale(VirtualWidth/float64(b.Dx()), VirtualHeight/float64(b.Dy()))
	}
	screen.DrawImage(g.Reg.OlympusBG, &g.uiOp)
	g.drawRect(screen, 0, 0, VirtualWidth, VirtualHeight, color.RGBA{8, 10, 20, 185})

	g.drawTextCentered(screen, "CHOOSE YOUR SACRED OLYMPIAN ARMAMENT", g.Reg.FaceSubtitle, 640, 62, color.RGBA{255, 220, 80, 255})
	g.drawTextCentered(screen, "Every Olympian God Boon dynamically adapts to empower either your Bow Volleys or your Sword Cleaves", g.Reg.FaceBody, 640, 102, color.RGBA{195, 215, 245, 255})

	// Card 1: Golden Bow of Apollo
	bowSelected := g.SelectedWeapon == combat.WeaponBow
	g.drawWeaponCard(screen, 150, 155, 460, 425, combat.BowOfApolloSpec, g.PreviewHeroBow, "[1] BOW OF APOLLO (RANGED)", bowSelected)

	// Card 2: Bronze Sword of Ares
	swordSelected := g.SelectedWeapon == combat.WeaponSword
	g.drawWeaponCard(screen, 670, 155, 460, 425, combat.SwordOfAresSpec, g.PreviewHeroSword, "[2] SWORD OF ARES (MELEE)", swordSelected)

	g.drawPanel(screen, 340, 610, 600, 54, color.RGBA{22, 26, 44, 240}, color.RGBA{240, 195, 65, 255})
	g.drawTextCentered(screen, "CLICK A WEAPON CARD OR PRESS [1] / [2] + [ENTER] TO BEGIN TRIAL", g.Reg.FaceHeading, 640, 628, color.RGBA{255, 235, 130, 255})
}

func (g *Game) drawWeaponCard(screen *ebiten.Image, x, y, w, h float64, spec combat.WeaponSpec, preview *sprite.AnimationController, badge string, selected bool) {
	bg := color.RGBA{18, 22, 38, 235}
	border := color.RGBA{90, 105, 140, 220}
	if selected {
		bg = color.RGBA{32, 36, 58, 248}
		border = spec.AccentColor
	}
	g.drawPanel(screen, x, y, w, h, bg, border)
	if selected {
		g.drawRect(screen, x+4, y+4, w-8, 6, spec.AccentColor)
	}

	cx := x + w*0.5
	g.drawTextCentered(screen, badge, g.Reg.FaceSmall, cx, y+22, spec.AccentColor)
	g.drawTextCentered(screen, spec.Name, g.Reg.FaceSubtitle, cx, y+46, color.RGBA{255, 255, 255, 255})
	g.drawTextCentered(screen, spec.Subtitle, g.Reg.FaceBody, cx, y+80, color.RGBA{185, 205, 235, 255})

	// Sprite showcase box
	g.drawPanel(screen, cx-75, y+112, 150, 125, color.RGBA{12, 14, 26, 240}, spec.AccentColor)
	g.drawUISprite(screen, preview, cx, y+175, 1.85)

	// Stats breakdown
	g.drawTextCentered(screen, fmt.Sprintf("Base Damage: %.0f   |   Attack Cooldown: %.2fs   |   Range: %.0f px", spec.BaseDamage, spec.Cooldown, spec.Range), g.Reg.FaceSmall, cx, y+255, color.RGBA{255, 230, 135, 255})
	g.drawTextCentered(screen, spec.Playstyle, g.Reg.FaceBody, cx, y+292, color.RGBA{230, 235, 245, 255})
	g.drawTextCentered(screen, spec.SpecialTrait, g.Reg.FaceSmall, cx, y+328, spec.AccentColor)

	btnLabel := "SELECT ARMAMENT"
	if selected {
		btnLabel = "★ EQUIPPED — CLICK OR PRESS ENTER ★"
	}
	g.drawPanel(screen, x+45, y+364, w-90, 42, color.RGBA{28, 30, 48, 255}, border)
	g.drawTextCentered(screen, btnLabel, g.Reg.FaceSmall, cx, y+377, color.RGBA{255, 245, 190, 255})
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	a := g.Arena
	h := &a.Hero

	// 1. TOP-LEFT: Hero HP Bar, Weapon Badge, Dash Cooldown, Overdrives & Active Boons
	g.drawPanel(screen, 16, 14, 335, 86, color.RGBA{14, 16, 28, 220}, color.RGBA{195, 160, 75, 230})

	// HP Bar
	hpRatio := math.Max(0, math.Min(1, h.HP/math.Max(1, h.MaxHP)))
	g.drawRect(screen, 28, 24, 210, 18, color.RGBA{38, 16, 24, 240})
	g.drawRect(screen, 28, 24, 210*hpRatio, 18, color.RGBA{225, 50, 65, 255})
	g.drawTextLeft(screen, fmt.Sprintf("HP %.0f / %.0f", h.HP, h.MaxHP), g.Reg.FaceSmall, 34, 25, color.RGBA{255, 255, 255, 255})

	// Dash Status
	dashText := "DASH READY [SPACE]"
	dashCol := color.RGBA{95, 235, 255, 255}
	if h.DashCooldown > 0 {
		dashText = fmt.Sprintf("DASH %.1fs", h.DashCooldown)
		dashCol = color.RGBA{150, 165, 190, 255}
	}
	g.drawTextLeft(screen, dashText, g.Reg.FaceTiny, 246, 26, dashCol)

	// Equipped Weapon & Stat Multipliers
	wName := "Bow of Apollo"
	if h.Weapon == combat.WeaponSword {
		wName = "Sword of Ares"
	}
	g.drawTextLeft(screen, fmt.Sprintf("%s  •  DMG x%.2f  SPD x%.2f  MV x%.2f",
		wName, h.Stats.DamageMul, h.Stats.AttackSpeedMul, h.Stats.MoveSpeedMul),
		g.Reg.FaceTiny, 28, 49, h.WeaponSpec.AccentColor)

	// Active Overdrive Pills
	ox := 28.0
	if h.OverdriveSpd > 0 {
		g.drawTextLeft(screen, fmt.Sprintf("⚡QUIVER %.1fs", h.OverdriveSpd), g.Reg.FaceTiny, ox, 70, color.RGBA{255, 225, 70, 255})
		ox += 102
	}
	if h.OverdriveDmg > 0 {
		g.drawTextLeft(screen, fmt.Sprintf("🔥EMBER %.1fs", h.OverdriveDmg), g.Reg.FaceTiny, ox, 70, color.RGBA{255, 95, 85, 255})
		ox += 102
	}
	if h.OverdriveMove > 0 {
		g.drawTextLeft(screen, fmt.Sprintf("💨PLUME %.1fs", h.OverdriveMove), g.Reg.FaceTiny, ox, 70, color.RGBA{90, 225, 255, 255})
		ox += 102
	}
	if ox == 28.0 {
		g.drawTextLeft(screen, fmt.Sprintf("Boons: %d   |   Relic Orbs: %d   |   Chests: %d", h.TotalBoons, h.PowerUpsPicked, a.ChestsOpened), g.Reg.FaceTiny, 28, 70, color.RGBA{195, 205, 225, 255})
	}

	// Row of owned Olympian God Boon Badges under Top-Left HUD
	bx := 16.0
	for id := combat.BoonID(0); id < combat.BoonCount; id++ {
		stacks := h.Boons[id]
		if stacks <= 0 {
			continue
		}
		def := combat.AllBoons[id]
		g.drawPanel(screen, bx, 106, 48, 22, color.RGBA{18, 20, 34, 230}, def.AccentColor)
		g.drawTextCentered(screen, fmt.Sprintf("%s x%d", def.Symbol, stacks), g.Reg.FaceTiny, bx+24, 110, def.AccentColor)
		bx += 52
	}

	// 2. TOP-CENTER: Ambrosia XP Bar, Level & Active Boss Health Bar
	g.drawPanel(screen, 390, 14, 500, 52, color.RGBA{14, 16, 28, 220}, color.RGBA{235, 190, 65, 235})
	xpRatio := math.Max(0, math.Min(1, h.XP/math.Max(1, h.NextLevelXP)))
	g.drawRect(screen, 404, 24, 472, 16, color.RGBA{34, 28, 14, 240})
	g.drawRect(screen, 404, 24, 472*xpRatio, 16, color.RGBA{245, 195, 55, 255})
	g.drawTextCentered(screen, fmt.Sprintf("HERO LEVEL %d   •   AMBROSIA XP %.0f / %.0f", h.Level, h.XP, h.NextLevelXP), g.Reg.FaceSmall, 640, 24, color.RGBA{255, 250, 220, 255})
	g.drawTextCentered(screen, fmt.Sprintf("OLYMPIAN TRIAL WAVE %d / 3", a.Wave), g.Reg.FaceTiny, 640, 45, color.RGBA{145, 215, 255, 255})

	// Boss Health Bar when Zeus's Champion or Titan Cronus is active!
	if boss := a.ActiveBoss(); boss != nil {
		g.drawPanel(screen, 350, 74, 580, 40, color.RGBA{24, 12, 22, 238}, color.RGBA{255, 195, 70, 255})
		bRatio := math.Max(0, math.Min(1, boss.HP/math.Max(1, boss.MaxHP)))
		g.drawRect(screen, 362, 94, 556, 12, color.RGBA{45, 16, 24, 255})
		g.drawRect(screen, 362, 94, 556*bRatio, 12, color.RGBA{245, 65, 85, 255})
		g.drawTextCentered(screen, fmt.Sprintf("★ %s  (%.0f / %.0f HP) ★", boss.Name, boss.HP, boss.MaxHP), g.Reg.FaceSmall, 640, 77, color.RGBA{255, 225, 110, 255})
	} else if a.WaveBannerTimer > 0 {
		g.drawPanel(screen, 340, 76, 600, 34, color.RGBA{18, 22, 38, 230}, color.RGBA{245, 200, 75, 255})
		g.drawTextCentered(screen, a.WaveBannerText, g.Reg.FaceSmall, 640, 85, color.RGBA{255, 230, 110, 255})
	}

	// 3. TOP-RIGHT: Run Timer, Slain Counter & Controls Hint
	g.drawPanel(screen, 930, 14, 334, 62, color.RGBA{14, 16, 28, 220}, color.RGBA{120, 165, 225, 220})
	mins := int(a.RunTime) / 60
	secs := int(a.RunTime) % 60
	g.drawTextLeft(screen, fmt.Sprintf("TIME %02d:%02d   |   SLAIN: %d   |   FOES: %d", mins, secs, a.EnemiesSlain, a.ActiveEnemyCount()), g.Reg.FaceSmall, 944, 24, color.RGBA{240, 245, 255, 255})
	muteLabel := "[M] Audio: ON"
	if g.Audio.IsMuted() {
		muteLabel = "[M] Audio: MUTED"
	}
	g.drawTextLeft(screen, fmt.Sprintf("[P/ESC] Pause   |   %s", muteLabel), g.Reg.FaceTiny, 944, 48, color.RGBA{175, 195, 225, 255})

	// Demo Mode Banner or On-Screen Mobile/Click Dodge Button
	if a.DemoMode {
		g.drawPanel(screen, 410, 650, 460, 44, color.RGBA{18, 20, 36, 235}, color.RGBA{255, 215, 75, 255})
		g.drawTextCentered(screen, "ATTRACT DEMO MODE — PRESS ANY KEY OR CLICK TO PLAY", g.Reg.FaceSmall, 640, 664, color.RGBA{255, 230, 110, 255})
	} else {
		dashBorder := color.RGBA{95, 235, 255, 230}
		dashBg := color.RGBA{16, 28, 44, 215}
		dashLabel := "DODGE [TAP]"
		if h.DashCooldown > 0 {
			dashBorder = color.RGBA{105, 115, 140, 190}
			dashBg = color.RGBA{18, 20, 32, 190}
			dashLabel = fmt.Sprintf("WAIT %.1fs", h.DashCooldown)
		}
		g.drawPanel(screen, 1125, 620, 138, 66, dashBg, dashBorder)
		g.drawTextCentered(screen, dashLabel, g.Reg.FaceSmall, 1194, 636, dashBorder)
		g.drawTextCentered(screen, "SPACE / SHIFT", g.Reg.FaceTiny, 1194, 660, color.RGBA{185, 205, 235, 255})
	}
}

func (g *Game) drawBoonSelectOverlay(screen *ebiten.Image) {
	g.drawRect(screen, 0, 0, VirtualWidth, VirtualHeight, color.RGBA{6, 8, 18, 198})

	h := &g.Arena.Hero
	g.drawTextCentered(screen, "DIVINE ASCENSION — BESTOWED BY THE 12 OLYMPIANS", g.Reg.FaceSubtitle, 640, 74, color.RGBA{255, 220, 80, 255})
	wLabel := "Golden Bow of Apollo"
	if h.Weapon == combat.WeaponSword {
		wLabel = "Bronze Sword of Ares"
	}
	g.drawTextCentered(screen, fmt.Sprintf("Select 1 of 3 Sacred Boons to empower your %s  (Press [1], [2], [3] or Click Card)", wLabel), g.Reg.FaceBody, 640, 112, color.RGBA{195, 215, 245, 255})

	for idx := 0; idx < 3; idx++ {
		boon := g.BoonChoices[idx]
		cx := 125.0 + float64(idx)*355.0
		cy := 165.0
		cw := 320.0
		ch := 390.0
		hovered := g.HoveredCard == idx

		bg := color.RGBA{18, 22, 38, 242}
		border := color.RGBA{95, 110, 145, 220}
		if hovered {
			bg = color.RGBA{30, 36, 60, 252}
			border = boon.AccentColor
		}
		g.drawPanel(screen, cx, cy, cw, ch, bg, border)
		g.drawRect(screen, cx+4, cy+4, cw-8, 7, boon.AccentColor)

		midX := cx + cw*0.5
		currRank := h.Boons[boon.ID]
		g.drawTextCentered(screen, fmt.Sprintf("[%d] OLYMPIAN GOD %s  •  RANK %d -> %d", idx+1, boon.Symbol, currRank, currRank+1), g.Reg.FaceTiny, midX, cy+24, boon.AccentColor)
		g.drawTextCentered(screen, boon.GodName, g.Reg.FaceSubtitle, midX, cy+52, color.RGBA{255, 255, 255, 255})
		g.drawTextCentered(screen, boon.GodTitle, g.Reg.FaceSmall, midX, cy+86, color.RGBA{185, 200, 230, 255})

		// Emblem medallion
		g.drawPanel(screen, midX-48, cy+118, 96, 76, color.RGBA{12, 14, 24, 245}, boon.AccentColor)
		g.drawTextCentered(screen, boon.Symbol, g.Reg.FaceTitle, midX, cy+132, boon.AccentColor)

		// Boon Name & Weapon-Specific Effect
		g.drawTextCentered(screen, boon.BoonTitle, g.Reg.FaceHeading, midX, cy+214, boon.AccentColor)
		desc := boon.DescriptionFor(h.Weapon)
		g.drawWrappedTextCentered(screen, desc, g.Reg.FaceBody, midX, cy+252, cw-36, color.RGBA{235, 240, 255, 255})

		btnText := fmt.Sprintf("PRESS [%d] OR CLICK", idx+1)
		if hovered {
			btnText = fmt.Sprintf("★ CLAIM %s'S BOON [%d] ★", boon.GodName, idx+1)
		}
		g.drawPanel(screen, cx+24, cy+330, cw-48, 40, color.RGBA{24, 28, 46, 255}, border)
		g.drawTextCentered(screen, btnText, g.Reg.FaceSmall, midX, cy+342, color.RGBA{255, 240, 175, 255})
	}
}

func (g *Game) drawPauseOverlay(screen *ebiten.Image) {
	g.drawRect(screen, 0, 0, VirtualWidth, VirtualHeight, color.RGBA{8, 10, 20, 190})
	g.drawPanel(screen, 390, 210, 500, 290, color.RGBA{16, 20, 36, 245}, color.RGBA{240, 195, 65, 255})
	g.drawTextCentered(screen, "TRIAL OF OLYMPUS PAUSED", g.Reg.FaceSubtitle, 640, 248, color.RGBA{255, 220, 80, 255})
	g.drawTextCentered(screen, "Press [ESC], [P], [ENTER] or Click to Resume Battle", g.Reg.FaceHeading, 640, 320, color.RGBA{235, 240, 255, 255})
	g.drawTextCentered(screen, "Press [1] to Switch Weapon / Restart Trial", g.Reg.FaceBody, 640, 370, color.RGBA{185, 205, 235, 255})
	g.drawTextCentered(screen, "Press [M] to Toggle Audio Mute", g.Reg.FaceSmall, 640, 420, color.RGBA{160, 175, 205, 255})
}

func (g *Game) drawEndScreenOverlay(screen *ebiten.Image, victory bool) {
	g.drawRect(screen, 0, 0, VirtualWidth, VirtualHeight, color.RGBA{8, 10, 22, 205})
	a := g.Arena
	h := &a.Hero

	border := color.RGBA{255, 215, 70, 255}
	title := "VICTORY! TITAN CRONUS OVERTHROWN!"
	subtitle := "You have ascended Mount Olympus and claimed your seat among the Gods!"
	if !victory {
		border = color.RGBA{235, 65, 80, 255}
		title = "FALLEN ON MOUNT OLYMPUS"
		subtitle = "The legions of Cronus overwhelmed you, but the Olympian Gods await your return!"
	}

	g.drawPanel(screen, 260, 130, 760, 460, color.RGBA{16, 20, 36, 248}, border)
	g.drawTextCentered(screen, title, g.Reg.FaceSubtitle, 640, 168, border)
	g.drawTextCentered(screen, subtitle, g.Reg.FaceBody, 640, 212, color.RGBA{220, 230, 245, 255})

	wName := "Golden Bow of Apollo (Ranged)"
	if h.Weapon == combat.WeaponSword {
		wName = "Bronze Sword of Ares (Melee)"
	}
	mins := int(a.RunTime) / 60
	secs := int(a.RunTime) % 60

	g.drawPanel(screen, 330, 255, 620, 210, color.RGBA{12, 14, 26, 240}, color.RGBA{95, 115, 155, 220})
	g.drawTextCentered(screen, fmt.Sprintf("Sacred Armament: %s", wName), g.Reg.FaceHeading, 640, 276, h.WeaponSpec.AccentColor)
	g.drawTextCentered(screen, fmt.Sprintf("Hero Level Reached: %d   |   Olympian Boons Claimed: %d", h.Level, h.TotalBoons), g.Reg.FaceBody, 640, 314, color.RGBA{240, 245, 255, 255})
	g.drawTextCentered(screen, fmt.Sprintf("Enemies Slain: %d   |   Total Damage Dealt: %.0f", a.EnemiesSlain, a.TotalDamageDone), g.Reg.FaceBody, 640, 348, color.RGBA{240, 245, 255, 255})
	g.drawTextCentered(screen, fmt.Sprintf("Olympus Chests Opened: %d   |   Power-Up Relics: %d   |   Time: %02d:%02d", a.ChestsOpened, h.PowerUpsPicked, mins, secs), g.Reg.FaceBody, 640, 382, color.RGBA{255, 220, 115, 255})
	g.drawTextCentered(screen, fmt.Sprintf("Final Multipliers: Damage x%.2f  •  Attack Speed x%.2f  •  Move Speed x%.2f", h.Stats.DamageMul, h.Stats.AttackSpeedMul, h.Stats.MoveSpeedMul), g.Reg.FaceSmall, 640, 422, color.RGBA{145, 215, 255, 255})

	g.drawPanel(screen, 410, 498, 460, 54, color.RGBA{32, 28, 16, 255}, border)
	g.drawTextCentered(screen, "PRESS ENTER / SPACE OR CLICK TO PLAY AGAIN", g.Reg.FaceHeading, 640, 515, color.RGBA{255, 235, 130, 255})
}

func (g *Game) drawUISprite(screen *ebiten.Image, ctrl *sprite.AnimationController, cx, cy, scale float64) {
	if ctrl == nil || ctrl.Sheet == nil {
		return
	}
	fr := ctrl.Sheet.Frame(ctrl.CurrentFrame)
	if fr == nil {
		return
	}
	fw := float64(ctrl.Sheet.FrameWidth)
	fh := float64(ctrl.Sheet.FrameHeight)
	g.uiOp.GeoM.Reset()
	g.uiOp.ColorScale.Reset()
	g.uiOp.GeoM.Translate(-fw*0.5, -fh*0.5)
	g.uiOp.GeoM.Scale(scale, scale)
	g.uiOp.GeoM.Translate(math.Round(cx), math.Round(cy))
	screen.DrawImage(fr, &g.uiOp)
}

func (g *Game) drawPanel(screen *ebiten.Image, x, y, w, h float64, bg, border color.RGBA) {
	g.drawRect(screen, x, y, w, h, border)
	g.drawRect(screen, x+2, y+2, w-4, h-4, bg)
}

func (g *Game) drawRect(screen *ebiten.Image, x, y, w, h float64, col color.RGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	g.uiOp.GeoM.Reset()
	g.uiOp.ColorScale.Reset()
	g.uiOp.GeoM.Scale(w, h)
	g.uiOp.GeoM.Translate(math.Round(x), math.Round(y))
	g.uiOp.ColorScale.ScaleWithColor(col)
	screen.DrawImage(g.Reg.WhitePixel, &g.uiOp)
}

func (g *Game) drawTextCentered(screen *ebiten.Image, str string, face *text.GoTextFace, cx, y float64, col color.RGBA) {
	var op text.DrawOptions
	op.GeoM.Translate(math.Round(cx), math.Round(y))
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	op.ColorScale.ScaleWithColor(col)
	text.Draw(screen, str, face, &op)
}

func (g *Game) drawTextLeft(screen *ebiten.Image, str string, face *text.GoTextFace, x, y float64, col color.RGBA) {
	var op text.DrawOptions
	op.GeoM.Translate(math.Round(x), math.Round(y))
	op.LayoutOptions.PrimaryAlign = text.AlignStart
	op.ColorScale.ScaleWithColor(col)
	text.Draw(screen, str, face, &op)
}

func (g *Game) drawWrappedTextCentered(screen *ebiten.Image, str string, face *text.GoTextFace, cx, y, maxW float64, col color.RGBA) {
	// Simple 2-line word wrapper for Boon descriptions
	if len(str) <= 38 {
		g.drawTextCentered(screen, str, face, cx, y, col)
		return
	}
	splitIdx := len(str) / 2
	for i := splitIdx; i < len(str); i++ {
		if str[i] == ' ' {
			splitIdx = i
			break
		}
	}
	line1 := str[:splitIdx]
	line2 := str[splitIdx+1:]
	g.drawTextCentered(screen, line1, face, cx, y, col)
	g.drawTextCentered(screen, line2, face, cx, y+22, col)
	_ = maxW
}
