package entities

import (
	"image/color"
	"math"

	"github.com/drivakos/arrows-of-olympus/internal/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Draw renders the entire Olympus courtyard, entities, projectiles, slash waves, and floating text
// with zero per-frame image allocations.
func (a *Arena) Draw(screen *ebiten.Image) {
	// 1. Draw Mount Olympus Pantheon Background (1280x720)
	a.drawOp.GeoM.Reset()
	a.drawOp.ColorScale.Reset()
	bgBounds := a.Reg.OlympusBG.Bounds()
	if bgBounds.Dx() > 0 && bgBounds.Dy() > 0 {
		a.drawOp.GeoM.Scale(VirtualWidth/float64(bgBounds.Dx()), VirtualHeight/float64(bgBounds.Dy()))
	}
	a.Cam.Apply(&a.drawOp)
	screen.DrawImage(a.Reg.OlympusBG, &a.drawOp)

	// 2. Draw Chronos Time-Rift Vortex if active
	if a.TimeRift.Active {
		pulse := 0.85 + 0.15*math.Sin(a.RunTime*8.0)
		a.drawWorldRing(screen, a.TimeRift.X, a.TimeRift.Y, a.TimeRift.Radius*pulse, 3.5, color.RGBA{130, 85, 220, 140})
		a.drawWorldRing(screen, a.TimeRift.X, a.TimeRift.Y, a.TimeRift.Radius*0.55, 2.5, color.RGBA{75, 195, 240, 140})
	}

	// 3. Draw Boss Telegraphs (Zeus's Champion Thunder Spear Dash line)
	for i := range a.Enemies {
		e := &a.Enemies[i]
		if e.Active && !e.Dying && e.DashTelegraph > 0 {
			a.drawWorldLine(screen, e.X, e.Y, e.DashTargetX, e.DashTargetY, 26.0, color.RGBA{180, 35, 45, 95})
			a.drawWorldLine(screen, e.X, e.Y, e.DashTargetX, e.DashTargetY, 4.0, color.RGBA{255, 215, 70, 220})
		}
	}

	// 4. Draw Pickups & Power-Up Relics
	for i := range a.Pickups {
		p := &a.Pickups[i]
		if !p.Active {
			continue
		}
		var tint color.RGBA
		scale := 1.0
		switch p.Kind {
		case PickupXP:
			tint = color.RGBA{255, 255, 255, 255}
		case PickupAttackSpeed:
			tint = color.RGBA{255, 230, 85, 255}
			scale = 1.15
			a.drawWorldRing(screen, p.X, p.Y, 16.0, 2.0, color.RGBA{255, 215, 55, 180})
		case PickupAttackDamage:
			tint = color.RGBA{255, 95, 85, 255}
			scale = 1.15
			a.drawWorldRing(screen, p.X, p.Y, 16.0, 2.0, color.RGBA{255, 75, 65, 180})
		case PickupMoveSpeed:
			tint = color.RGBA{95, 235, 255, 255}
			scale = 1.15
			a.drawWorldRing(screen, p.X, p.Y, 16.0, 2.0, color.RGBA{75, 215, 255, 180})
		case PickupNectarHeal:
			tint = color.RGBA{120, 255, 155, 255}
			scale = 1.15
			a.drawWorldRing(screen, p.X, p.Y, 16.0, 2.0, color.RGBA{95, 240, 135, 180})
		}
		a.drawSpriteCentered(screen, p.Anim, p.X, p.Y, scale, false, 0, tint)
	}

	// 5. Draw Golden Olympus Chests
	for i := range a.Chests {
		c := &a.Chests[i]
		if !c.Active {
			continue
		}
		a.drawSpriteCentered(screen, c.Anim, c.X, c.Y, 1.15, false, 0, color.RGBA{255, 255, 255, 255})
	}

	// 6. Draw Sword of Ares Melee Cleave Slash Arcs
	for i := range a.SlashArcs {
		sa := &a.SlashArcs[i]
		if !sa.Active {
			continue
		}
		alpha := sa.Timer / sa.MaxTimer
		col := premulAlpha(sa.Color, alpha)
		steps := 14
		if sa.Whirlwind {
			steps = 24
		}
		startAng := sa.Angle - sa.SweepRad*0.5
		stepAng := sa.SweepRad / float64(steps)
		for s := 0; s < steps; s++ {
			a0 := startAng + float64(s)*stepAng
			a1 := a0 + stepAng
			x0 := sa.X + math.Cos(a0)*sa.Radius
			y0 := sa.Y + math.Sin(a0)*sa.Radius
			x1 := sa.X + math.Cos(a1)*sa.Radius
			y1 := sa.Y + math.Sin(a1)*sa.Radius
			a.drawWorldLine(screen, x0, y0, x1, y1, 9.0*alpha, col)
			a.drawWorldLine(screen, sa.X+math.Cos(a0)*(sa.Radius*0.65), sa.Y+math.Sin(a0)*(sa.Radius*0.65), x0, y0, 4.0*alpha, col)
		}
	}

	// 7. Draw Enemies & Bosses
	for i := range a.Enemies {
		e := &a.Enemies[i]
		if !e.Active {
			continue
		}
		scale := 1.12
		if e.Kind == EnemyBossZeusChampion {
			scale = 1.25
		} else if e.Kind == EnemyBossCronus {
			scale = 1.32
		}
		flash := e.HitFlash * 5.0
		if e.ShieldFlash > 0 {
			a.drawWorldRing(screen, e.X+e.FacingX*10, e.Y, e.Radius+6, 3.0, color.RGBA{255, 225, 90, 220})
		}
		tint := color.RGBA{255, 255, 255, 255}
		if e.PoisonTimer > 0 {
			tint = color.RGBA{225, 165, 255, 255}
		}
		a.drawSpriteCentered(screen, e.Anim, e.X, e.Y, scale, e.FlipX, flash, tint)

		// Overhead health bar for damaged mobs
		if !e.Dying && e.HP < e.MaxHP && e.Kind != EnemyBossZeusChampion && e.Kind != EnemyBossCronus {
			barW := 34.0
			barH := 4.0
			bx := e.X - barW*0.5
			by := e.Y - e.Radius - 14.0
			a.drawWorldRect(screen, bx-1, by-1, barW+2, barH+2, color.RGBA{15, 12, 22, 220})
			ratio := math.Max(0, math.Min(1, e.HP/e.MaxHP))
			a.drawWorldRect(screen, bx, by, barW*ratio, barH, color.RGBA{230, 55, 65, 255})
		}
	}

	// 8. Draw Hero & Active Overdrive Auras + Orbiting Athena Shields
	h := &a.Hero
	if h.OverdriveSpd > 0 {
		a.drawWorldRing(screen, h.X, h.Y, 26.0, 2.2, color.RGBA{255, 220, 65, 195})
	}
	if h.OverdriveDmg > 0 {
		a.drawWorldRing(screen, h.X, h.Y, 30.0, 2.2, color.RGBA{255, 75, 70, 195})
	}
	if h.OverdriveMove > 0 {
		a.drawWorldRing(screen, h.X, h.Y, 34.0, 2.0, color.RGBA{75, 220, 255, 195})
	}

	heroScale := 1.22
	a.drawSpriteCentered(screen, h.Anim, h.X, h.Y, heroScale, h.FlipX, h.HitFlash*4.0, color.RGBA{255, 255, 255, 255})

	// Orbiting Athena Aegis Shields
	if h.Stats.OrbitShields > 0 {
		n := h.Stats.OrbitShields
		for s := 0; s < n; s++ {
			ang := h.OrbitAngle + float64(s)*(2.0*math.Pi/float64(n))
			sx := h.X + math.Cos(ang)*58.0
			sy := h.Y + math.Sin(ang)*58.0
			a.drawWorldRect(screen, sx-7, sy-7, 14, 14, color.RGBA{255, 215, 75, 245})
			a.drawWorldRect(screen, sx-4, sy-4, 8, 8, color.RGBA{110, 215, 255, 255})
		}
	}

	// 9. Draw Olympus Marble Pillars
	for i := range a.Pillars {
		pil := &a.Pillars[i]
		a.drawSpriteCentered(screen, pil.Anim, pil.X, pil.Y, 1.18, false, 0, color.RGBA{255, 255, 255, 255})
	}

	// 10. Draw Player Projectiles (Golden Arrows & Sword Crescent Slash Waves)
	for i := range a.PlayerProjs {
		p := &a.PlayerProjs[i]
		if !p.Active {
			continue
		}
		if p.IsSwordWave {
			// Draw glowing crescent wave perpendicular to travel angle
			perpX := -math.Sin(p.Angle)
			perpY := math.Cos(p.Angle)
			halfSpan := p.Radius * 1.25
			x0 := p.X - perpX*halfSpan - math.Cos(p.Angle)*6
			y0 := p.Y - perpY*halfSpan - math.Sin(p.Angle)*6
			xm := p.X + math.Cos(p.Angle)*8
			ym := p.Y + math.Sin(p.Angle)*8
			x1 := p.X + perpX*halfSpan - math.Cos(p.Angle)*6
			y1 := p.Y + perpY*halfSpan - math.Sin(p.Angle)*6
			a.drawWorldLine(screen, x0, y0, xm, ym, 5.5, p.Color)
			a.drawWorldLine(screen, xm, ym, x1, y1, 5.5, p.Color)
			a.drawWorldLine(screen, x0, y0, x1, y1, 2.5, color.RGBA{255, 245, 200, 240})
		} else {
			// Draw Golden Arrow of Apollo shaft + arrowhead
			tailX := p.X - math.Cos(p.Angle)*16.0
			tailY := p.Y - math.Sin(p.Angle)*16.0
			tipX := p.X + math.Cos(p.Angle)*8.0
			tipY := p.Y + math.Sin(p.Angle)*8.0
			a.drawWorldLine(screen, tailX, tailY, tipX, tipY, 3.2, p.Color)
			a.drawWorldLine(screen, p.X, p.Y, tipX, tipY, 5.2, color.RGBA{255, 250, 210, 255})
		}
	}

	// 11. Draw Enemy Projectiles (Lightning Nova Orbs & Cronus Sickle Waves)
	for i := range a.EnemyProjs {
		ep := &a.EnemyProjs[i]
		if !ep.Active {
			continue
		}
		if ep.IsCrescent {
			perpX := -math.Sin(ep.Angle)
			perpY := math.Cos(ep.Angle)
			x0 := ep.X - perpX*18
			y0 := ep.Y - perpY*18
			x1 := ep.X + perpX*18
			y1 := ep.Y + perpY*18
			a.drawWorldLine(screen, x0, y0, x1, y1, 6.0, ep.Color)
		} else {
			a.drawWorldRect(screen, ep.X-ep.Radius, ep.Y-ep.Radius, ep.Radius*2, ep.Radius*2, ep.Color)
			a.drawWorldRect(screen, ep.X-ep.Radius*0.5, ep.Y-ep.Radius*0.5, ep.Radius, ep.Radius, color.RGBA{245, 250, 255, 255})
		}
	}

	// 12. Draw Zeus Chain Lightning Arcs
	for i := range a.LightningArcs {
		la := &a.LightningArcs[i]
		if !la.Active {
			continue
		}
		mx := (la.X1+la.X2)*0.5 + 8.0
		my := (la.Y1+la.Y2)*0.5 - 8.0
		a.drawWorldLine(screen, la.X1, la.Y1, mx, my, 3.2, la.Color)
		a.drawWorldLine(screen, mx, my, la.X2, la.Y2, 3.2, la.Color)
	}

	// 13. Draw Particles
	for i := range a.Particles {
		pt := &a.Particles[i]
		if !pt.Active {
			continue
		}
		alpha := pt.Timer / pt.MaxTimer
		col := premulAlpha(pt.Color, alpha)
		a.drawWorldRect(screen, pt.X-pt.Size*0.5, pt.Y-pt.Size*0.5, pt.Size, pt.Size, col)
	}

	// 14. Draw World-Space Floating Combat Numbers
	for i := range a.FloatingTexts {
		ft := &a.FloatingTexts[i]
		if !ft.Active {
			continue
		}
		face := a.Reg.FaceCombat
		if ft.Big {
			face = a.Reg.FaceHeading
		}
		var top text.DrawOptions
		top.GeoM.Translate(math.Round(ft.X+a.Cam.ShakeOffsetX), math.Round(ft.Y+a.Cam.ShakeOffsetY))
		top.LayoutOptions.PrimaryAlign = text.AlignCenter
		top.ColorScale.ScaleWithColor(ft.Color)
		text.Draw(screen, ft.Text, face, &top)
	}
}

func (a *Arena) drawSpriteCentered(screen *ebiten.Image, ctrl *sprite.AnimationController, x, y, scale float64, flipX bool, flash float64, tint color.RGBA) {
	if ctrl == nil || ctrl.Sheet == nil {
		return
	}
	frame := ctrl.Sheet.Frame(ctrl.CurrentFrame)
	if frame == nil {
		return
	}
	fw := float64(ctrl.Sheet.FrameWidth)
	fh := float64(ctrl.Sheet.FrameHeight)

	a.drawOp.GeoM.Reset()
	a.drawOp.ColorScale.Reset()
	a.drawOp.GeoM.Translate(-fw*0.5, -fh*0.5)
	if flipX {
		a.drawOp.GeoM.Scale(-scale, scale)
	} else {
		a.drawOp.GeoM.Scale(scale, scale)
	}
	a.drawOp.GeoM.Translate(math.Round(x), math.Round(y))
	a.Cam.Apply(&a.drawOp)

	a.drawOp.ColorScale.ScaleWithColor(tint)
	if flash > 0 {
		boost := float32(1.0 + math.Min(2.5, flash))
		a.drawOp.ColorScale.Scale(boost, boost, boost, 1.0)
	}
	screen.DrawImage(frame, &a.drawOp)
}

func (a *Arena) drawWorldRect(screen *ebiten.Image, x, y, w, h float64, col color.RGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	a.drawOp.GeoM.Reset()
	a.drawOp.ColorScale.Reset()
	a.drawOp.GeoM.Scale(w, h)
	a.drawOp.GeoM.Translate(math.Round(x), math.Round(y))
	a.Cam.Apply(&a.drawOp)
	a.drawOp.ColorScale.ScaleWithColor(col)
	screen.DrawImage(a.Reg.WhitePixel, &a.drawOp)
}

func (a *Arena) drawWorldLine(screen *ebiten.Image, x1, y1, x2, y2, thick float64, col color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	length := math.Hypot(dx, dy)
	if length < 1.0 {
		return
	}
	ang := math.Atan2(dy, dx)
	a.drawOp.GeoM.Reset()
	a.drawOp.ColorScale.Reset()
	a.drawOp.GeoM.Translate(0, -0.5)
	a.drawOp.GeoM.Scale(length, thick)
	a.drawOp.GeoM.Rotate(ang)
	a.drawOp.GeoM.Translate(x1, y1)
	a.Cam.Apply(&a.drawOp)
	a.drawOp.ColorScale.ScaleWithColor(col)
	screen.DrawImage(a.Reg.WhitePixel, &a.drawOp)
}

func (a *Arena) drawWorldRing(screen *ebiten.Image, cx, cy, r, thick float64, col color.RGBA) {
	const segs = 16
	step := (2.0 * math.Pi) / float64(segs)
	for i := 0; i < segs; i++ {
		a0 := float64(i) * step
		a1 := a0 + step
		a.drawWorldLine(
			screen,
			cx+math.Cos(a0)*r, cy+math.Sin(a0)*r,
			cx+math.Cos(a1)*r, cy+math.Sin(a1)*r,
			thick, col,
		)
	}
}

func premulAlpha(c color.RGBA, alpha float64) color.RGBA {
	alpha = math.Max(0, math.Min(1, alpha))
	return color.RGBA{
		R: uint8(float64(c.R) * alpha),
		G: uint8(float64(c.G) * alpha),
		B: uint8(float64(c.B) * alpha),
		A: uint8(float64(c.A) * alpha),
	}
}
