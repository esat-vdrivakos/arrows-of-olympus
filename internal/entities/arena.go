package entities

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/drivakos/arrows-of-olympus/internal/assets"
	"github.com/drivakos/arrows-of-olympus/internal/audio"
	"github.com/drivakos/arrows-of-olympus/internal/camera"
	"github.com/drivakos/arrows-of-olympus/internal/combat"
	"github.com/drivakos/arrows-of-olympus/internal/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	VirtualWidth  = 1280.0
	VirtualHeight = 720.0

	ArenaMinX = 72.0
	ArenaMaxX = 1208.0
	ArenaMinY = 86.0
	ArenaMaxY = 668.0

	MaxEnemies       = 96
	MaxPlayerProj    = 256
	MaxEnemyProj     = 256
	MaxPickups       = 256
	MaxChests        = 16
	MaxSlashArcs     = 32
	MaxLightningArcs = 32
	MaxFloatingTexts = 128
	MaxParticles     = 384
)

// EnemyKind enumerates all regular mobs and Olympian bosses.
type EnemyKind int

const (
	EnemySkeleton EnemyKind = iota
	EnemyWarriorUnshielded
	EnemyWarriorShielded
	EnemyBossZeusChampion
	EnemyBossCronus
)

// PickupKind enumerates XP gems, Nectar heal chalices, and the 3 Divine Stat Power-Ups.
type PickupKind int

const (
	PickupXP PickupKind = iota
	PickupAttackSpeed
	PickupAttackDamage
	PickupMoveSpeed
	PickupNectarHeal
)

// Pillar represents one of the 4 tactical marble columns in the Pantheon courtyard.
type Pillar struct {
	X, Y   float64
	Radius float64
	Anim   *sprite.AnimationController
}

// Chest represents a breakable Golden Olympus Treasure Chest.
type Chest struct {
	Active    bool
	Opening   bool
	X, Y      float64
	HP, MaxHP float64
	OpenTimer float64
	Anim      *sprite.AnimationController
}

// Pickup represents an Ambrosia XP gem or Divine Power-Up Relic on the arena floor.
type Pickup struct {
	Active bool
	Kind   PickupKind
	X, Y   float64
	VX, VY float64
	Value  float64
	Age    float64
	Anim   *sprite.AnimationController
}

// PlayerProj represents a Golden Arrow, Sword Crescent Wave, Piercing Blade Beam, or Tidal Shockwave.
type PlayerProj struct {
	Active         bool
	IsSwordWave    bool
	IsBladeBeam    bool
	X, Y           float64
	VX, VY         float64
	Angle          float64
	Radius         float64
	Damage         float64
	PierceLeft     int
	BounceLeft     int
	Life           float64
	MaxLife        float64
	Color          color.RGBA
	LastHitEnemyID int
}

// EnemyProj represents a boss lightning orb, crescent sickle wave, or time shard.
type EnemyProj struct {
	Active     bool
	IsCrescent bool
	X, Y       float64
	VX, VY     float64
	Angle      float64
	Radius     float64
	Damage     float64
	Life       float64
	Color      color.RGBA
}

// SlashArc represents a visible melee sword cleave sweep on screen.
type SlashArc struct {
	Active    bool
	X, Y      float64
	Angle     float64
	SweepRad  float64
	Radius    float64
	Timer     float64
	MaxTimer  float64
	Color     color.RGBA
	Whirlwind bool
}

// LightningArc represents a chain lightning bolt between two world points.
type LightningArc struct {
	Active         bool
	X1, Y1, X2, Y2 float64
	Timer          float64
	Color          color.RGBA
}

// TimeRift represents Titan Cronus's arena slow-field vortex.
type TimeRift struct {
	Active bool
	X, Y   float64
	Radius float64
	Timer  float64
}

// FloatingText represents world-space combat numbers and pickup banners.
type FloatingText struct {
	Active   bool
	X, Y     float64
	VY       float64
	Text     string
	Color    color.RGBA
	Timer    float64
	MaxTimer float64
	Big      bool
}

// Particle represents a zero-allocation spark or hit burst particle.
type Particle struct {
	Active   bool
	X, Y     float64
	VX, VY   float64
	Size     float64
	Timer    float64
	MaxTimer float64
	Color    color.RGBA
}

// Enemy represents an active mob or boss in the pre-allocated pool.
type Enemy struct {
	Active      bool
	Dying       bool
	ID          int
	Kind        EnemyKind
	Name        string
	X, Y        float64
	VX, VY      float64
	Radius      float64
	HP, MaxHP   float64
	Speed       float64
	ContactDmg  float64
	FlipX       bool
	FacingX     float64
	FacingY     float64
	HitFlash    float64
	ShieldFlash float64
	PoisonTimer float64
	PoisonDPS   float64
	DeathTimer  float64

	// AI & Boss state
	AttackTimer   float64
	PhaseTimer    float64
	DashTelegraph float64
	DashActive    float64
	DashTargetX   float64
	DashTargetY   float64
	DashDirX      float64
	DashDirY      float64
	SpiralAngle   float64
	Anim          *sprite.AnimationController
}

// Hero stores the player's state, weapon, Boons, and active power-up overdrives.
type Hero struct {
	X, Y            float64
	VX, VY          float64
	Radius          float64
	FlipX           bool
	AimX, AimY      float64
	HP, MaxHP       float64
	BaseSpeed       float64
	Weapon          combat.WeaponType
	WeaponSpec      combat.WeaponSpec
	Boons           [combat.BoonCount]int
	Stats           combat.ComputedStats
	TotalBoons      int
	Level           int
	XP, NextLevelXP float64
	PendingLevelUps int

	AttackCooldown float64
	AttackAnimHold float64
	ComboToggle    bool

	DashTimer    float64
	DashCooldown float64
	DashVX       float64
	DashVY       float64
	InvulnTimer  float64
	HitFlash     float64

	// Permanent & Overdrive Power-up Orbs
	PermAttackSpeed float64
	PermDamage      float64
	PermMoveSpeed   float64
	OverdriveSpd    float64
	OverdriveDmg    float64
	OverdriveMove   float64
	PowerUpsPicked  int

	OrbitAngle float64
	Anim       *sprite.AnimationController
}

// Arena holds the entire gameplay simulation state and pre-allocated pools.
type Arena struct {
	Reg   *assets.Registry
	Audio *audio.Manager
	Cam   *camera.Camera
	RNG   *rand.Rand

	DemoMode bool

	Hero    Hero
	Pillars [4]Pillar
	Chests  [MaxChests]Chest
	Enemies [MaxEnemies]Enemy
	Pickups [MaxPickups]Pickup

	PlayerProjs   [MaxPlayerProj]PlayerProj
	EnemyProjs    [MaxEnemyProj]EnemyProj
	SlashArcs     [MaxSlashArcs]SlashArc
	LightningArcs [MaxLightningArcs]LightningArc
	FloatingTexts [MaxFloatingTexts]FloatingText
	Particles     [MaxParticles]Particle
	TimeRift      TimeRift

	// Wave & Progression state
	Wave            int
	WaveElapsed     float64
	RunTime         float64
	SpawnTimer      float64
	ChestSpawnTimer float64
	NextEnemyID     int

	ZeusSpawned   bool
	ZeusDefeated  bool
	CronusSpawned bool
	CronusSlain   bool
	HeroDead      bool

	// Run statistics
	EnemiesSlain    int
	TotalDamageDone float64
	ChestsOpened    int
	WaveBannerText  string
	WaveBannerTimer float64

	// Reusable DrawImageOptions for zero-allocation rendering
	drawOp ebiten.DrawImageOptions
}

// NewArena constructs and pre-allocates the Olympus courtyard arena and entity pools.
func NewArena(reg *assets.Registry, aud *audio.Manager, cam *camera.Camera) (*Arena, error) {
	a := &Arena{
		Reg:   reg,
		Audio: aud,
		Cam:   cam,
		RNG:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	heroCtrl, err := reg.NewSpriteController("hero_greek")
	if err != nil {
		return nil, err
	}
	a.Hero.Anim = heroCtrl

	pillarCoords := [4][2]float64{
		{350, 240},
		{930, 240},
		{350, 500},
		{930, 500},
	}
	for i, pt := range pillarCoords {
		pCtrl, err := reg.NewSpriteController("prop_pillar")
		if err != nil {
			return nil, err
		}
		a.Pillars[i] = Pillar{
			X:      pt[0],
			Y:      pt[1],
			Radius: 28.0,
			Anim:   pCtrl,
		}
	}

	for i := range a.Chests {
		cCtrl, err := reg.NewSpriteController("prop_chest")
		if err != nil {
			return nil, err
		}
		a.Chests[i].Anim = cCtrl
	}

	for i := range a.Pickups {
		xpCtrl, err := reg.NewSpriteController("pickup_xp")
		if err != nil {
			return nil, err
		}
		a.Pickups[i].Anim = xpCtrl
	}

	return a, nil
}

// ResetRun initializes a fresh run with the player's selected weapon (Bow of Apollo or Sword of Ares).
func (a *Arena) ResetRun(weapon combat.WeaponType, demoMode bool) {
	a.DemoMode = demoMode
	a.Cam.Reset()

	spec := combat.GetWeaponSpec(weapon)
	a.Hero = Hero{
		X:           VirtualWidth * 0.5,
		Y:           VirtualHeight * 0.55,
		Radius:      18.0,
		AimX:        1.0,
		AimY:        0.0,
		HP:          120.0 + spec.BonusMaxHP,
		MaxHP:       120.0 + spec.BonusMaxHP,
		BaseSpeed:   225.0,
		Weapon:      weapon,
		WeaponSpec:  spec,
		Level:       1,
		XP:          0,
		NextLevelXP: 45.0,
		Anim:        a.Hero.Anim,
	}
	_ = a.Hero.Anim.Play("idle")
	a.RecomputeHeroStats()

	for i := range a.Enemies {
		a.Enemies[i].Active = false
	}
	for i := range a.PlayerProjs {
		a.PlayerProjs[i].Active = false
	}
	for i := range a.EnemyProjs {
		a.EnemyProjs[i].Active = false
	}
	for i := range a.Pickups {
		a.Pickups[i].Active = false
	}
	for i := range a.Chests {
		a.Chests[i].Active = false
	}
	for i := range a.SlashArcs {
		a.SlashArcs[i].Active = false
	}
	for i := range a.LightningArcs {
		a.LightningArcs[i].Active = false
	}
	for i := range a.FloatingTexts {
		a.FloatingTexts[i].Active = false
	}
	for i := range a.Particles {
		a.Particles[i].Active = false
	}
	a.TimeRift.Active = false

	a.Wave = 1
	a.WaveElapsed = 0
	a.RunTime = 0
	a.SpawnTimer = 0.6
	a.ChestSpawnTimer = 5.0
	a.NextEnemyID = 1
	a.ZeusSpawned = false
	a.ZeusDefeated = false
	a.CronusSpawned = false
	a.CronusSlain = false
	a.HeroDead = false
	a.EnemiesSlain = 0
	a.TotalDamageDone = 0
	a.ChestsOpened = 0

	a.TriggerWaveBanner("WAVE I: SKELETON SPARTOI & BRONZE INFANTRY")

	// Spawn 2 initial Olympian Treasure Chests in the courtyard
	a.SpawnChestAt(220, 370)
	a.SpawnChestAt(1060, 370)

	// Spawn initial Skeleton Spartoi pack so combat starts immediately
	a.SpawnEnemy(EnemySkeleton, 210, 160)
	a.SpawnEnemy(EnemySkeleton, 1070, 160)
	a.SpawnEnemy(EnemySkeleton, 210, 590)
	a.SpawnEnemy(EnemySkeleton, 1070, 590)
}

// RecomputeHeroStats recalculates the Hero's effective multipliers from weapon, Boons, and orbs.
func (a *Arena) RecomputeHeroStats() {
	h := &a.Hero
	h.Stats = combat.ComputeStats(
		h.Weapon,
		h.Boons,
		h.PermAttackSpeed,
		h.PermDamage,
		h.PermMoveSpeed,
		h.OverdriveSpd > 0,
		h.OverdriveDmg > 0,
		h.OverdriveMove > 0,
	)
}

// ApplyBoon grants one stack of the chosen Olympian God Boon and triggers immediate effects.
func (a *Arena) ApplyBoon(id combat.BoonID) {
	h := &a.Hero
	h.Boons[id]++
	h.TotalBoons++
	if h.PendingLevelUps > 0 {
		h.PendingLevelUps--
	}
	if id == combat.BoonDemeter {
		bonusHP := 40.0
		healHP := 50.0
		if h.Weapon == combat.WeaponSword {
			bonusHP = 45.0
			healHP = 60.0
		}
		h.MaxHP += bonusHP
		h.HP = math.Min(h.MaxHP, h.HP+healHP)
	}
	a.RecomputeHeroStats()
	a.Audio.PlaySFX("boon_levelup", 80*time.Millisecond)
	a.Cam.TriggerFlash(0.18)
	def := combat.AllBoons[id]
	a.SpawnFloatingText(h.X, h.Y-34, fmt.Sprintf("%s: %s!", def.GodName, def.BoonTitle), def.AccentColor, true)
	a.SpawnParticleBurst(h.X, h.Y, 22, def.AccentColor)
}

// TriggerWaveBanner displays a dramatic center-top wave announcement.
func (a *Arena) TriggerWaveBanner(text string) {
	a.WaveBannerText = text
	a.WaveBannerTimer = 3.6
}

// Update advances the arena simulation by dt seconds.
func (a *Arena) Update(dt float64, moveX, moveY float64, dashPressed bool) {
	if a.Cam.Update(dt) {
		return // Hitstop freeze frame active
	}
	if a.HeroDead || a.CronusSlain {
		return
	}

	a.RunTime += dt
	a.WaveElapsed += dt
	if a.WaveBannerTimer > 0 {
		a.WaveBannerTimer = math.Max(0, a.WaveBannerTimer-dt)
	}

	dtDur := time.Duration(dt * float64(time.Second))
	for i := range a.Pillars {
		a.Pillars[i].Anim.Update(dtDur)
	}

	// Demo AI steering override when running in Attract Mode
	if a.DemoMode {
		moveX, moveY, dashPressed = a.computeDemoAI()
		if a.Hero.PendingLevelUps > 0 {
			draft := combat.DraftBoons(a.RNG)
			a.ApplyBoon(draft[0].ID)
		}
	}

	a.updateHero(dt, dtDur, moveX, moveY, dashPressed)
	a.updateWaveSpawner(dt)
	a.updateChests(dt, dtDur)
	a.updateEnemies(dt, dtDur)
	a.updatePlayerProjs(dt)
	a.updateEnemyProjs(dt)
	a.updatePickups(dt, dtDur)
	a.updateVisualEffects(dt)
}

func (a *Arena) computeDemoAI() (float64, float64, bool) {
	h := &a.Hero
	var vx, vy float64
	dash := false

	// Avoid nearby enemies or enemy projectiles
	for i := range a.Enemies {
		e := &a.Enemies[i]
		if !e.Active || e.Dying {
			continue
		}
		dx := h.X - e.X
		dy := h.Y - e.Y
		dist := math.Hypot(dx, dy)
		keepDist := 180.0
		if h.Weapon == combat.WeaponSword {
			keepDist = 95.0
		}
		if dist < keepDist && dist > 1e-3 {
			weight := (keepDist - dist) / keepDist
			vx += (dx / dist) * weight * 2.2
			vy += (dy / dist) * weight * 2.2
			if dist < 68.0 && h.DashCooldown <= 0 {
				dash = true
			}
		} else if h.Weapon == combat.WeaponSword && dist > 130.0 && dist < 420.0 {
			vx -= (dx / dist) * 0.8
			vy -= (dy / dist) * 0.8
		}
	}

	// Gravitate toward nearby pickups and chests
	for i := range a.Pickups {
		p := &a.Pickups[i]
		if !p.Active {
			continue
		}
		dx := p.X - h.X
		dy := p.Y - h.Y
		dist := math.Hypot(dx, dy)
		if dist > 10 && dist < 280 {
			vx += (dx / dist) * 0.9
			vy += (dy / dist) * 0.9
		}
	}

	// Keep away from arena walls
	cx := VirtualWidth*0.5 - h.X
	cy := VirtualHeight*0.5 - h.Y
	vx += cx * 0.0025
	vy += cy * 0.0025

	mag := math.Hypot(vx, vy)
	if mag > 1e-4 {
		vx /= mag
		vy /= mag
	}
	return vx, vy, dash
}

func (a *Arena) updateHero(dt float64, dtDur time.Duration, moveX, moveY float64, dashPressed bool) {
	h := &a.Hero

	// Update temporary Overdrive power-up timers
	recompute := false
	if h.OverdriveSpd > 0 {
		h.OverdriveSpd = math.Max(0, h.OverdriveSpd-dt)
		if h.OverdriveSpd == 0 {
			recompute = true
		}
	}
	if h.OverdriveDmg > 0 {
		h.OverdriveDmg = math.Max(0, h.OverdriveDmg-dt)
		if h.OverdriveDmg == 0 {
			recompute = true
		}
	}
	if h.OverdriveMove > 0 {
		h.OverdriveMove = math.Max(0, h.OverdriveMove-dt)
		if h.OverdriveMove == 0 {
			recompute = true
		}
	}
	if recompute {
		a.RecomputeHeroStats()
	}

	if h.DashCooldown > 0 {
		h.DashCooldown = math.Max(0, h.DashCooldown-dt)
	}
	if h.InvulnTimer > 0 {
		h.InvulnTimer = math.Max(0, h.InvulnTimer-dt)
	}
	if h.HitFlash > 0 {
		h.HitFlash = math.Max(0, h.HitFlash-dt)
	}
	if h.AttackCooldown > 0 {
		h.AttackCooldown = math.Max(0, h.AttackCooldown-dt)
	}
	if h.AttackAnimHold > 0 {
		h.AttackAnimHold = math.Max(0, h.AttackAnimHold-dt)
	}

	// Trigger Dodge Roll on Space / Shift
	if dashPressed && h.DashCooldown <= 0 && h.DashTimer <= 0 {
		dx, dy := moveX, moveY
		if dx == 0 && dy == 0 {
			dx, dy = h.AimX, h.AimY
		}
		mag := math.Hypot(dx, dy)
		if mag < 1e-4 {
			dx, dy = 1, 0
		} else {
			dx /= mag
			dy /= mag
		}
		h.DashVX = dx * 620.0
		h.DashVY = dy * 620.0
		h.DashTimer = 0.18
		h.DashCooldown = 2.4
		h.InvulnTimer = 0.28
		a.SpawnParticleBurst(h.X, h.Y, 10, color.RGBA{120, 220, 255, 255})
	}

	// Movement integration
	speed := h.BaseSpeed * h.Stats.MoveSpeedMul
	if a.TimeRift.Active && math.Hypot(h.X-a.TimeRift.X, h.Y-a.TimeRift.Y) < a.TimeRift.Radius {
		speed *= 0.62 // Chronos Time-Rift slow effect
	}

	if h.DashTimer > 0 {
		h.DashTimer = math.Max(0, h.DashTimer-dt)
		h.VX = h.DashVX
		h.VY = h.DashVY
		a.SpawnParticle(h.X, h.Y, 0, 0, 5.0, 0.16, color.RGBA{110, 205, 255, 190})
	} else {
		h.VX = moveX * speed
		h.VY = moveY * speed
	}

	h.X += h.VX * dt
	h.Y += h.VY * dt
	h.X, h.Y = a.clampAndResolvePillars(h.X, h.Y, h.Radius)

	if h.OverdriveMove > 0 && (math.Abs(h.VX) > 5 || math.Abs(h.VY) > 5) && a.RNG.Float64() < 0.45 {
		a.SpawnParticle(h.X, h.Y+10, -h.VX*0.15, -h.VY*0.15, 4.0, 0.22, color.RGBA{80, 220, 255, 220})
	}

	// Update Orbiting Athena Aegis Shields
	h.OrbitAngle += 3.2 * dt
	if h.Stats.OrbitShields > 0 {
		a.updateOrbitShieldCollisions()
	}

	// Auto-Aim & Auto-Attack nearest target (Enemy or Closed Chest)
	targetX, targetY, hasTarget := a.findNearestTarget(h.X, h.Y, h.WeaponSpec.Range*1.35)
	if hasTarget {
		dx := targetX - h.X
		dy := targetY - h.Y
		dist := math.Hypot(dx, dy)
		if dist > 1e-4 {
			h.AimX = dx / dist
			h.AimY = dy / dist
		}
		if h.AimX < -0.05 {
			h.FlipX = true
		} else if h.AimX > 0.05 {
			h.FlipX = false
		}

		// Fire Bow or Swing Sword when within effective range
		maxTriggerDist := h.WeaponSpec.Range
		if h.Weapon == combat.WeaponSword {
			maxTriggerDist = 230.0 * h.Stats.CleaveRadiusMul
			if h.Boons[combat.BoonHephaestus] > 0 || h.Boons[combat.BoonPoseidon] > 0 || h.Boons[combat.BoonApollo] > 0 {
				maxTriggerDist = 380.0
			}
		}
		if dist <= maxTriggerDist && h.AttackCooldown <= 0 {
			a.executeHeroAttack()
		}
	} else if math.Abs(moveX) > 0.05 {
		h.FlipX = moveX < 0
		h.AimX = moveX
		h.AimY = moveY
	}

	// Update Hero Sprite Animation Tag
	if h.HitFlash > 0.08 {
		_ = h.Anim.Play("hit")
	} else if h.AttackAnimHold > 0 {
		if h.Weapon == combat.WeaponSword {
			_ = h.Anim.Play("sword_attack")
		} else {
			_ = h.Anim.Play("bow_attack")
		}
	} else if math.Hypot(h.VX, h.VY) > 15.0 {
		_ = h.Anim.Play("run")
	} else {
		_ = h.Anim.Play("idle")
	}
	h.Anim.Update(dtDur)
}

func (a *Arena) findNearestTarget(x, y, maxRange float64) (float64, float64, bool) {
	bestDist := maxRange
	var bestX, bestY float64
	found := false

	for i := range a.Enemies {
		e := &a.Enemies[i]
		if !e.Active || e.Dying {
			continue
		}
		d := math.Hypot(e.X-x, e.Y-y)
		if d < bestDist {
			bestDist = d
			bestX = e.X
			bestY = e.Y
			found = true
		}
	}
	if found {
		return bestX, bestY, true
	}

	// If no enemy in range, auto-target nearby closed Golden Chests
	for i := range a.Chests {
		c := &a.Chests[i]
		if !c.Active || c.Opening {
			continue
		}
		d := math.Hypot(c.X-x, c.Y-y)
		if d < bestDist {
			bestDist = d
			bestX = c.X
			bestY = c.Y
			found = true
		}
	}
	return bestX, bestY, found
}

func (a *Arena) executeHeroAttack() {
	h := &a.Hero
	cd := h.WeaponSpec.Cooldown / math.Max(0.25, h.Stats.AttackSpeedMul)
	h.AttackCooldown = cd
	h.AttackAnimHold = math.Min(0.26, cd*0.85)
	baseAimAngle := math.Atan2(h.AimY, h.AimX)
	baseDmg := h.WeaponSpec.BaseDamage * h.Stats.DamageMul

	if h.Weapon == combat.WeaponBow {
		a.Audio.PlaySFX("bow_shoot", 45*time.Millisecond)
		projSpeed := 620.0

		// 1. Parallel Front Arrows (Base + Artemis Twin Strike)
		nFront := h.Stats.FrontProjectiles
		perpX := -math.Sin(baseAimAngle)
		perpY := math.Cos(baseAimAngle)
		spacing := 11.0
		for i := 0; i < nFront; i++ {
			offset := (float64(i) - float64(nFront-1)*0.5) * spacing
			sx := h.X + math.Cos(baseAimAngle)*18.0 + perpX*offset
			sy := h.Y + math.Sin(baseAimAngle)*18.0 + perpY*offset
			a.spawnPlayerProj(sx, sy, baseAimAngle, projSpeed, 7.0, baseDmg, h.Stats.PierceCount, h.Stats.RicochetCount, 1.35, false, false, color.RGBA{255, 215, 70, 255})
		}

		// 2. Apollo Sunburst Diagonal Cones (±18°, ±32°, ...)
		for p := 1; p <= h.Stats.DiagonalPairs; p++ {
			spread := float64(p) * (18.0 * math.Pi / 180.0)
			a.spawnPlayerProj(h.X, h.Y, baseAimAngle-spread, projSpeed, 6.5, baseDmg*0.85, h.Stats.PierceCount, h.Stats.RicochetCount, 1.25, false, false, color.RGBA{255, 190, 55, 255})
			a.spawnPlayerProj(h.X, h.Y, baseAimAngle+spread, projSpeed, 6.5, baseDmg*0.85, h.Stats.PierceCount, h.Stats.RicochetCount, 1.25, false, false, color.RGBA{255, 190, 55, 255})
		}

		// 3. Hera Sovereign Rearguard Arrows
		if h.Stats.RearAttack > 0 {
			rearAngle := baseAimAngle + math.Pi
			a.spawnPlayerProj(h.X, h.Y, rearAngle, projSpeed*0.95, 6.5, baseDmg*0.90, h.Stats.PierceCount, h.Stats.RicochetCount, 1.2, false, false, color.RGBA{130, 235, 225, 255})
		}
		return
	}

	// === SWORD OF ARES MELEE CLEAVE & SLASH-WAVE ===
	a.Audio.PlaySFX("sword_slash", 45*time.Millisecond)
	cleaveRadius := 135.0 * h.Stats.CleaveRadiusMul
	sweepRad := 115.0 * math.Pi / 180.0
	whirlwind := h.Stats.RearAttack > 0 // Hera unlocks 360° Whirlwind Back-Slash!
	if whirlwind {
		sweepRad = math.Pi * 2.0
	}

	// Artemis Twin Strike adds bonus multi-hit damage + double visual slash arc
	cleaveHits := 1
	if h.Boons[combat.BoonArtemis] > 0 {
		cleaveHits = 2
	}
	a.spawnSlashArc(h.X, h.Y, baseAimAngle, sweepRad, cleaveRadius, color.RGBA{255, 95, 75, 240}, whirlwind)
	if cleaveHits > 1 {
		a.spawnSlashArc(h.X, h.Y, baseAimAngle+0.18, sweepRad*0.92, cleaveRadius*1.08, color.RGBA{255, 215, 90, 220}, false)
	}

	// Apply instant melee cleave damage (SHATTERS Shielded Hoplite blocks for 100% damage!)
	totalCleaveDmg := baseDmg * (1.0 + float64(h.Boons[combat.BoonArtemis])*0.38)
	a.applyMeleeCleave(h.X, h.Y, baseAimAngle, sweepRad, cleaveRadius, totalCleaveDmg)

	// Launch standard short-range Crescent Slash Wave forward
	waveSpeed := 440.0
	a.spawnPlayerProj(h.X+math.Cos(baseAimAngle)*22, h.Y+math.Sin(baseAimAngle)*22, baseAimAngle, waveSpeed, 18.0, baseDmg*0.65, 1+h.Stats.PierceCount, h.Stats.RicochetCount, 0.42, true, false, color.RGBA{255, 110, 75, 245})

	// Hephaestus Adamant Edge launches a long-range Piercing Blade Beam!
	if h.Boons[combat.BoonHephaestus] > 0 {
		a.spawnPlayerProj(h.X+math.Cos(baseAimAngle)*26, h.Y+math.Sin(baseAimAngle)*26, baseAimAngle, 560.0, 22.0, baseDmg*0.85, 3+h.Stats.PierceCount, h.Stats.RicochetCount, 0.85, true, true, color.RGBA{255, 175, 60, 255})
	}

	// Apollo Sunburst Cone launches diagonal Golden Slash-Waves at ±28°!
	for p := 1; p <= h.Stats.DiagonalPairs; p++ {
		spread := float64(p) * (28.0 * math.Pi / 180.0)
		a.spawnPlayerProj(h.X, h.Y, baseAimAngle-spread, 460.0, 16.0, baseDmg*0.65, 1+h.Stats.PierceCount, h.Stats.RicochetCount, 0.52, true, false, color.RGBA{255, 210, 75, 245})
		a.spawnPlayerProj(h.X, h.Y, baseAimAngle+spread, 460.0, 16.0, baseDmg*0.65, 1+h.Stats.PierceCount, h.Stats.RicochetCount, 0.52, true, false, color.RGBA{255, 210, 75, 245})
	}

	// Poseidon Tidal Force launches bouncing Azure Shockwaves!
	if h.Boons[combat.BoonPoseidon] > 0 {
		a.spawnPlayerProj(h.X, h.Y, baseAimAngle, 480.0, 20.0, baseDmg*0.75, 2, 2+h.Stats.RicochetCount, 0.95, true, false, color.RGBA{70, 195, 255, 250})
	}
}

func (a *Arena) applyMeleeCleave(ox, oy, aimAngle, sweepRad, radius, dmg float64) {
	halfSweep := sweepRad * 0.5
	hitAny := false

	for i := range a.Enemies {
		e := &a.Enemies[i]
		if !e.Active || e.Dying {
			continue
		}
		dx := e.X - ox
		dy := e.Y - oy
		dist := math.Hypot(dx, dy)
		if dist > radius+e.Radius {
			continue
		}
		if sweepRad < math.Pi*1.95 {
			angDiff := normalizeAngle(math.Atan2(dy, dx) - aimAngle)
			if math.Abs(angDiff) > halfSweep {
				continue
			}
		}
		// Melee cleave shatters frontal shields for 100% damage!
		a.damageEnemy(e, dmg, dx, dy, true)
		hitAny = true
	}

	// Also cleave open Golden Olympus Chests in range!
	for i := range a.Chests {
		c := &a.Chests[i]
		if !c.Active || c.Opening {
			continue
		}
		if math.Hypot(c.X-ox, c.Y-oy) <= radius+24.0 {
			a.damageChest(c, dmg)
			hitAny = true
		}
	}

	if hitAny {
		a.Cam.AddTrauma(0.14)
	}
}

func (a *Arena) updateOrbitShieldCollisions() {
	h := &a.Hero
	n := h.Stats.OrbitShields
	orbitRadius := 58.0

	for s := 0; s < n; s++ {
		ang := h.OrbitAngle + float64(s)*(2.0*math.Pi/float64(n))
		sx := h.X + math.Cos(ang)*orbitRadius
		sy := h.Y + math.Sin(ang)*orbitRadius

		// Destroy incoming enemy projectiles
		for j := range a.EnemyProjs {
			ep := &a.EnemyProjs[j]
			if !ep.Active {
				continue
			}
			if math.Hypot(ep.X-sx, ep.Y-sy) <= ep.Radius+14.0 {
				ep.Active = false
				a.Audio.PlaySFX("shield_clang", 65*time.Millisecond)
				a.SpawnParticleBurst(sx, sy, 6, color.RGBA{255, 225, 110, 255})
			}
		}
	}
}

func (a *Arena) updateWaveSpawner(dt float64) {
	// Periodic Golden Chest spawns around the courtyard
	a.ChestSpawnTimer -= dt
	if a.ChestSpawnTimer <= 0 {
		a.ChestSpawnTimer = 14.0
		if a.countActiveChests() < 3 {
			cx := 180.0 + a.RNG.Float64()*(VirtualWidth-360.0)
			cy := 150.0 + a.RNG.Float64()*(VirtualHeight-300.0)
			a.SpawnChestAt(cx, cy)
		}
	}

	// Wave progression transitions
	if a.Wave == 1 && a.WaveElapsed >= 32.0 {
		a.Wave = 2
		a.WaveElapsed = 0
		a.TriggerWaveBanner("WAVE II: PHALANX ASSAULT & ZEUS'S CHAMPION!")
		a.SpawnEnemy(EnemyBossZeusChampion, VirtualWidth*0.5, 135.0)
		a.ZeusSpawned = true
		a.Audio.PlaySFX("thunder_boss", 100*time.Millisecond)
		a.Cam.AddTrauma(0.45)
	} else if a.Wave == 2 && (a.ZeusDefeated || a.WaveElapsed >= 45.0) && !a.CronusSpawned {
		a.Wave = 3
		a.WaveElapsed = 0
		a.TriggerWaveBanner("FINAL WAVE III: TITAN CRONUS, LORD OF TIME!")
		a.SpawnEnemy(EnemyBossCronus, VirtualWidth*0.5, 140.0)
		a.CronusSpawned = true
		a.Audio.PlaySFX("thunder_boss", 100*time.Millisecond)
		a.Cam.AddTrauma(0.65)
		a.Cam.TriggerFlash(0.25)
	}

	// Continuous mob spawning
	a.SpawnTimer -= dt
	if a.SpawnTimer > 0 {
		return
	}

	activeMobs := a.countActiveEnemies()
	if activeMobs >= 26 {
		a.SpawnTimer = 1.0
		return
	}

	switch a.Wave {
	case 1:
		a.SpawnTimer = 2.0
		// Spawn Skeleton pack + occasional Unshielded/Shielded Greek Warrior
		for k := 0; k < 2; k++ {
			sx, sy := a.randomPerimeterPoint()
			a.SpawnEnemy(EnemySkeleton, sx, sy)
		}
		if a.WaveElapsed > 8.0 {
			sx, sy := a.randomPerimeterPoint()
			if a.WaveElapsed > 18.0 && a.RNG.Float64() < 0.45 {
				a.SpawnEnemy(EnemyWarriorShielded, sx, sy)
			} else {
				a.SpawnEnemy(EnemyWarriorUnshielded, sx, sy)
			}
		}
	case 2:
		a.SpawnTimer = 2.3
		sx, sy := a.randomPerimeterPoint()
		r := a.RNG.Float64()
		switch {
		case r < 0.40:
			a.SpawnEnemy(EnemySkeleton, sx, sy)
		case r < 0.72:
			a.SpawnEnemy(EnemyWarriorUnshielded, sx, sy)
		default:
			a.SpawnEnemy(EnemyWarriorShielded, sx, sy)
		}
	default:
		a.SpawnTimer = 2.6
		sx, sy := a.randomPerimeterPoint()
		r := a.RNG.Float64()
		switch {
		case r < 0.35:
			a.SpawnEnemy(EnemySkeleton, sx, sy)
		case r < 0.68:
			a.SpawnEnemy(EnemyWarriorUnshielded, sx, sy)
		default:
			a.SpawnEnemy(EnemyWarriorShielded, sx, sy)
		}
	}
}

func (a *Arena) randomPerimeterPoint() (float64, float64) {
	side := a.RNG.Intn(4)
	switch side {
	case 0: // Top
		return 110.0 + a.RNG.Float64()*(VirtualWidth-220.0), 102.0
	case 1: // Bottom
		return 110.0 + a.RNG.Float64()*(VirtualWidth-220.0), 648.0
	case 2: // Left
		return 96.0, 120.0 + a.RNG.Float64()*(VirtualHeight-240.0)
	default: // Right
		return 1184.0, 120.0 + a.RNG.Float64()*(VirtualHeight-240.0)
	}
}

// SpawnChestAt activates a Golden Olympus Chest in the courtyard.
func (a *Arena) SpawnChestAt(x, y float64) {
	x, y = a.clampAndResolvePillars(x, y, 22.0)
	for i := range a.Chests {
		c := &a.Chests[i]
		if !c.Active {
			c.Active = true
			c.Opening = false
			c.X = x
			c.Y = y
			c.HP = 30.0
			c.MaxHP = 30.0
			c.OpenTimer = 0
			_ = c.Anim.Play("closed")
			a.SpawnParticleBurst(x, y, 10, color.RGBA{255, 215, 75, 255})
			return
		}
	}
}

func (a *Arena) updateChests(dt float64, dtDur time.Duration) {
	h := &a.Hero
	for i := range a.Chests {
		c := &a.Chests[i]
		if !c.Active {
			continue
		}
		c.Anim.Update(dtDur)
		if c.Opening {
			c.OpenTimer -= dt
			if c.OpenTimer <= 0 {
				c.Active = false
			}
			continue
		}
		// Hero can also pop open a chest by touching it directly!
		if math.Hypot(h.X-c.X, h.Y-c.Y) <= h.Radius+24.0 {
			a.openChest(c)
		}
	}
}

func (a *Arena) damageChest(c *Chest, dmg float64) {
	if !c.Active || c.Opening {
		return
	}
	c.HP -= dmg
	a.SpawnParticleBurst(c.X, c.Y, 5, color.RGBA{255, 220, 90, 255})
	if c.HP <= 0 {
		a.openChest(c)
	}
}

func (a *Arena) openChest(c *Chest) {
	if c.Opening {
		return
	}
	c.Opening = true
	c.OpenTimer = 0.55
	a.ChestsOpened++
	_ = c.Anim.Play("open")
	a.Audio.PlaySFX("chest_open", 60*time.Millisecond)
	a.SpawnFloatingText(c.X, c.Y-24, "OLYMPUS CHEST!", color.RGBA{255, 220, 75, 255}, true)
	a.SpawnParticleBurst(c.X, c.Y, 20, color.RGBA{255, 215, 65, 255})

	// Drop 1-2 guaranteed Divine Stat Power-Up Relics + shower of Ambrosia XP shards
	kinds := [4]PickupKind{PickupAttackSpeed, PickupAttackDamage, PickupMoveSpeed, PickupNectarHeal}
	p1 := kinds[a.RNG.Intn(3)] // Guaranteed Attack Speed, Damage, or Move Speed
	a.SpawnPickup(p1, c.X-18, c.Y, 1.0)
	if a.RNG.Float64() < 0.65 {
		p2 := kinds[a.RNG.Intn(4)]
		a.SpawnPickup(p2, c.X+18, c.Y, 1.0)
	}
	for k := 0; k < 4; k++ {
		ang := float64(k) * (math.Pi * 0.5)
		a.SpawnPickup(PickupXP, c.X+math.Cos(ang)*22, c.Y+math.Sin(ang)*22, 12.0)
	}
}

// SpawnEnemy activates an enemy or boss from the pre-allocated pool.
func (a *Arena) SpawnEnemy(kind EnemyKind, x, y float64) {
	for i := range a.Enemies {
		e := &a.Enemies[i]
		if e.Active {
			continue
		}
		var spriteName string
		switch kind {
		case EnemySkeleton:
			spriteName = "mob_skeleton"
		case EnemyWarriorUnshielded:
			spriteName = "mob_warrior_unshielded"
		case EnemyWarriorShielded:
			spriteName = "mob_warrior_shielded"
		case EnemyBossZeusChampion:
			spriteName = "boss_zeus_champion"
		case EnemyBossCronus:
			spriteName = "boss_cronus"
		}
		ctrl, err := a.Reg.NewSpriteController(spriteName)
		if err != nil {
			return
		}
		*e = Enemy{
			Active:  true,
			ID:      a.NextEnemyID,
			Kind:    kind,
			X:       x,
			Y:       y,
			FacingX: 1,
			FacingY: 0,
			Anim:    ctrl,
		}
		a.NextEnemyID++

		switch kind {
		case EnemySkeleton:
			e.Name = "Skeleton Spartoi"
			e.Radius = 16.0
			e.MaxHP = 42.0
			e.HP = e.MaxHP
			e.Speed = 115.0
			e.ContactDmg = 10.0
		case EnemyWarriorUnshielded:
			e.Name = "Greek Warrior"
			e.Radius = 18.0
			e.MaxHP = 88.0
			e.HP = e.MaxHP
			e.Speed = 102.0
			e.ContactDmg = 14.0
		case EnemyWarriorShielded:
			e.Name = "Shielded Hoplite"
			e.Radius = 20.0
			e.MaxHP = 150.0
			e.HP = e.MaxHP
			e.Speed = 76.0
			e.ContactDmg = 16.0
		case EnemyBossZeusChampion:
			e.Name = "ZEUS'S CHAMPION"
			e.Radius = 34.0
			e.MaxHP = 1200.0
			e.HP = e.MaxHP
			e.Speed = 92.0
			e.ContactDmg = 22.0
			e.PhaseTimer = 2.8
		case EnemyBossCronus:
			e.Name = "TITAN CRONUS, LORD OF TIME"
			e.Radius = 42.0
			e.MaxHP = 2500.0
			e.HP = e.MaxHP
			e.Speed = 78.0
			e.ContactDmg = 26.0
			e.PhaseTimer = 2.5
		}
		_ = e.Anim.Play("idle")
		a.SpawnParticleBurst(x, y, 8, color.RGBA{220, 190, 255, 220})
		return
	}
}

func (a *Arena) updateEnemies(dt float64, dtDur time.Duration) {
	h := &a.Hero

	for i := range a.Enemies {
		e := &a.Enemies[i]
		if !e.Active {
			continue
		}
		e.Anim.Update(dtDur)

		if e.Dying {
			e.DeathTimer -= dt
			if e.DeathTimer <= 0 {
				e.Active = false
			}
			continue
		}

		if e.HitFlash > 0 {
			e.HitFlash = math.Max(0, e.HitFlash-dt)
		}
		if e.ShieldFlash > 0 {
			e.ShieldFlash = math.Max(0, e.ShieldFlash-dt)
		}

		// Dionysus Poison Miasma DoT
		if e.PoisonTimer > 0 {
			e.PoisonTimer = math.Max(0, e.PoisonTimer-dt)
			e.HP -= e.PoisonDPS * dt
			if a.RNG.Float64() < 0.15 {
				a.SpawnParticle(e.X, e.Y-8, 0, -25, 4.0, 0.28, color.RGBA{205, 80, 215, 220})
			}
			if e.HP <= 0 {
				a.killEnemy(e)
				continue
			}
		}

		dx := h.X - e.X
		dy := h.Y - e.Y
		dist := math.Hypot(dx, dy)
		if dist > 1e-4 {
			e.FacingX = dx / dist
			e.FacingY = dy / dist
			e.FlipX = e.FacingX < 0
		}

		switch e.Kind {
		case EnemyBossZeusChampion:
			a.updateBossZeusChampion(e, dt, dist)
		case EnemyBossCronus:
			a.updateBossCronus(e, dt, dist)
		default:
			spd := e.Speed
			if e.Kind == EnemyWarriorUnshielded && dist < 145.0 {
				spd *= 1.38 // Aggressive warrior lunge
			}
			e.VX = e.FacingX * spd
			e.VY = e.FacingY * spd
			e.X += e.VX * dt
			e.Y += e.VY * dt

			if e.ShieldFlash > 0.05 && e.Kind == EnemyWarriorShielded {
				_ = e.Anim.Play("shield_guard")
			} else if dist < e.Radius+h.Radius+16.0 {
				_ = e.Anim.Play("attack")
			} else {
				_ = e.Anim.Play("run")
			}
		}

		// Soft separation between enemies so mobs don't collapse into a single point
		for j := i + 1; j < len(a.Enemies); j++ {
			other := &a.Enemies[j]
			if !other.Active || other.Dying {
				continue
			}
			ox := e.X - other.X
			oy := e.Y - other.Y
			od := math.Hypot(ox, oy)
			minD := (e.Radius + other.Radius) * 0.78
			if od > 1e-3 && od < minD {
				push := (minD - od) * 0.5
				nx, ny := ox/od, oy/od
				e.X += nx * push
				e.Y += ny * push
				other.X -= nx * push
				other.Y -= ny * push
			}
		}

		e.X, e.Y = a.clampAndResolvePillars(e.X, e.Y, e.Radius)

		// Contact damage against Hero
		if math.Hypot(h.X-e.X, h.Y-e.Y) <= h.Radius+e.Radius-4.0 {
			a.damageHero(e.ContactDmg)
		}
	}
}

func (a *Arena) updateBossZeusChampion(e *Enemy, dt, dist float64) {
	h := &a.Hero

	if e.DashTelegraph > 0 {
		e.DashTelegraph -= dt
		_ = e.Anim.Play("cast")
		if e.DashTelegraph <= 0 {
			e.DashActive = 0.42
			a.Audio.PlaySFX("thunder_boss", 80*time.Millisecond)
		}
		return
	}

	if e.DashActive > 0 {
		e.DashActive -= dt
		e.X += e.DashDirX * 480.0 * dt
		e.Y += e.DashDirY * 480.0 * dt
		_ = e.Anim.Play("attack")
		a.SpawnParticle(e.X, e.Y, 0, 0, 6.0, 0.22, color.RGBA{120, 200, 255, 230})
		return
	}

	e.PhaseTimer -= dt
	if e.PhaseTimer <= 0 {
		if a.RNG.Float64() < 0.52 {
			// Phase A: Thunder Spear Dash with red/gold telegraph line!
			e.DashTelegraph = 0.72
			e.DashDirX = e.FacingX
			e.DashDirY = e.FacingY
			e.DashTargetX = e.X + e.DashDirX*340.0
			e.DashTargetY = e.Y + e.DashDirY*340.0
			e.PhaseTimer = 3.2
		} else {
			// Phase B: 8-Way Lightning Nova
			_ = e.Anim.Play("cast")
			a.Audio.PlaySFX("thunder_boss", 80*time.Millisecond)
			for k := 0; k < 8; k++ {
				ang := float64(k) * (math.Pi / 4.0)
				a.spawnEnemyProj(e.X, e.Y, ang, 255.0, 9.0, 16.0, 3.2, false, color.RGBA{110, 205, 255, 255})
			}
			e.PhaseTimer = 2.9
		}
		return
	}

	e.X += e.FacingX * e.Speed * dt
	e.Y += e.FacingY * e.Speed * dt
	if dist < 80 {
		_ = e.Anim.Play("attack")
	} else {
		_ = e.Anim.Play("move")
	}
	_ = h
}

func (a *Arena) updateBossCronus(e *Enemy, dt, dist float64) {
	h := &a.Hero
	enraged := e.HP < e.MaxHP*0.42

	e.PhaseTimer -= dt
	if e.PhaseTimer <= 0 {
		roll := a.RNG.Float64()
		switch {
		case roll < 0.50:
			// Signature 1: Adamant Sickle Cleave (3 sweeping crescent waves toward Hero)
			_ = e.Anim.Play("cleave")
			a.Audio.PlaySFX("thunder_boss", 80*time.Millisecond)
			baseAng := math.Atan2(h.Y-e.Y, h.X-e.X)
			for k := -1; k <= 1; k++ {
				ang := baseAng + float64(k)*0.26
				a.spawnEnemyProj(e.X, e.Y, ang, 290.0, 16.0, 20.0, 3.4, true, color.RGBA{255, 95, 120, 255})
			}
			e.PhaseTimer = 2.4
		default:
			// Signature 2: Chronos Time-Rift Vortex + Spartoi Summon
			_ = e.Anim.Play("summon")
			a.Audio.PlaySFX("thunder_boss", 80*time.Millisecond)
			a.TimeRift = TimeRift{
				Active: true,
				X:      h.X,
				Y:      h.Y,
				Radius: 115.0,
				Timer:  5.0,
			}
			a.SpawnEnemy(EnemySkeleton, math.Max(110, e.X-90), e.Y+40)
			a.SpawnEnemy(EnemyWarriorShielded, math.Min(1170, e.X+90), e.Y+40)
			e.PhaseTimer = 3.6
		}
	}

	// Signature 3: Titan's Wrath (<42% HP Enrage spiral barrage of golden time-shards)
	if enraged {
		e.AttackTimer -= dt
		if e.AttackTimer <= 0 {
			e.AttackTimer = 0.32
			e.SpiralAngle += 0.55
			a.spawnEnemyProj(e.X, e.Y, e.SpiralAngle, 240.0, 8.0, 14.0, 3.5, false, color.RGBA{255, 205, 70, 255})
			a.spawnEnemyProj(e.X, e.Y, e.SpiralAngle+math.Pi, 240.0, 8.0, 14.0, 3.5, false, color.RGBA{255, 205, 70, 255})
		}
	}

	e.X += e.FacingX * e.Speed * dt
	e.Y += e.FacingY * e.Speed * dt
	if dist < 95 {
		_ = e.Anim.Play("cleave")
	} else {
		_ = e.Anim.Play("move")
	}
}

func (a *Arena) damageEnemy(e *Enemy, rawDmg, hitDX, hitDY float64, isMelee bool) {
	if !e.Active || e.Dying {
		return
	}
	h := &a.Hero
	dmg := rawDmg

	// Shielded Hoplite Warrior blocks 50% of FRONTAL ranged arrows (Melee Sword of Ares shatters right through!)
	blocked := false
	if !isMelee && e.Kind == EnemyWarriorShielded {
		mag := math.Hypot(hitDX, hitDY)
		if mag > 1e-4 {
			nx, ny := hitDX/mag, hitDY/mag
			// Projectile traveling opposite to the Hoplite's facing direction hits the frontal Aspis shield
			dot := nx*e.FacingX + ny*e.FacingY
			if dot < -0.15 {
				dmg *= 0.50
				blocked = true
				e.ShieldFlash = 0.25
				a.Audio.PlaySFX("shield_clang", 55*time.Millisecond)
				a.SpawnFloatingText(e.X, e.Y-26, fmt.Sprintf("BLOCKED %.0f", dmg), color.RGBA{110, 200, 255, 255}, false)
				a.SpawnParticleBurst(e.X+e.FacingX*12, e.Y+e.FacingY*12, 7, color.RGBA{255, 225, 110, 255})
			}
		}
	}

	// Ares Bloodlust Critical Hit roll
	isCrit := false
	if !blocked && a.RNG.Float64() < h.Stats.CritChance {
		dmg *= 2.5
		isCrit = true
	}

	e.HP -= dmg
	e.HitFlash = 0.14
	a.TotalDamageDone += dmg

	if !blocked {
		a.Audio.PlaySFX("arrow_hit", 40*time.Millisecond)
		if isCrit {
			a.SpawnFloatingText(e.X, e.Y-24, fmt.Sprintf("CRIT %.0f!", dmg), color.RGBA{255, 215, 55, 255}, true)
		} else {
			a.SpawnFloatingText(e.X, e.Y-20, fmt.Sprintf("%.0f", dmg), color.RGBA{245, 245, 255, 255}, false)
		}
		a.SpawnParticleBurst(e.X, e.Y, 6, color.RGBA{255, 195, 95, 255})
	}

	// Knockback (boosted by Poseidon Tidal Force or Melee Sword of Ares)
	kb := 8.0
	if isMelee || h.Boons[combat.BoonPoseidon] > 0 {
		kb = 22.0
	}
	if e.Kind == EnemyBossZeusChampion || e.Kind == EnemyBossCronus {
		kb *= 0.18
	}
	mag := math.Hypot(hitDX, hitDY)
	if mag > 1e-4 {
		e.X += (hitDX / mag) * kb
		e.Y += (hitDY / mag) * kb
	}

	// Dionysus Intoxicating Miasma Poison DoT
	if h.Stats.PoisonDPS > 0 {
		e.PoisonTimer = 3.5
		e.PoisonDPS = h.Stats.PoisonDPS
	}

	// Zeus Thunderbolt Strike Chain Lightning
	if h.Stats.ChainLightning > 0 {
		a.triggerChainLightning(e, rawDmg*0.55, h.Stats.ChainLightning)
	}

	if e.HP <= 0 {
		a.killEnemy(e)
	}
}

func (a *Arena) triggerChainLightning(source *Enemy, arcDmg float64, jumps int) {
	currX, currY := source.X, source.Y
	lastID := source.ID

	for j := 0; j < jumps; j++ {
		bestDist := 210.0
		var target *Enemy
		for i := range a.Enemies {
			cand := &a.Enemies[i]
			if !cand.Active || cand.Dying || cand.ID == lastID {
				continue
			}
			d := math.Hypot(cand.X-currX, cand.Y-currY)
			if d < bestDist {
				bestDist = d
				target = cand
			}
		}
		if target == nil {
			break
		}
		a.spawnLightningArc(currX, currY, target.X, target.Y, color.RGBA{135, 215, 255, 255})
		target.HP -= arcDmg
		target.HitFlash = 0.12
		a.TotalDamageDone += arcDmg
		a.SpawnFloatingText(target.X, target.Y-18, fmt.Sprintf("⚡%.0f", arcDmg), color.RGBA{135, 215, 255, 255}, false)
		currX, currY = target.X, target.Y
		lastID = target.ID
		if target.HP <= 0 {
			a.killEnemy(target)
		}
	}
}

func (a *Arena) killEnemy(e *Enemy) {
	if e.Dying {
		return
	}
	e.Dying = true
	e.DeathTimer = 0.42
	_ = e.Anim.Play("death")
	a.EnemiesSlain++

	h := &a.Hero
	if h.Stats.LifeOnKill > 0 {
		h.HP = math.Min(h.MaxHP, h.HP+h.Stats.LifeOnKill)
	}

	a.SpawnParticleBurst(e.X, e.Y, 16, color.RGBA{255, 160, 75, 255})

	// Drop Table per GDD Section 2
	switch e.Kind {
	case EnemySkeleton:
		a.SpawnPickup(PickupXP, e.X, e.Y, 14.0)
		if a.RNG.Float64() < 0.14 {
			kinds := [4]PickupKind{PickupAttackSpeed, PickupAttackDamage, PickupMoveSpeed, PickupNectarHeal}
			a.SpawnPickup(kinds[a.RNG.Intn(4)], e.X+10, e.Y, 1.0)
		}
	case EnemyWarriorUnshielded, EnemyWarriorShielded:
		a.SpawnPickup(PickupXP, e.X, e.Y, 26.0)
		if a.RNG.Float64() < 0.26 {
			kinds := [3]PickupKind{PickupAttackSpeed, PickupAttackDamage, PickupMoveSpeed}
			a.SpawnPickup(kinds[a.RNG.Intn(3)], e.X+12, e.Y, 1.0)
		}
	case EnemyBossZeusChampion:
		a.ZeusDefeated = true
		a.Cam.AddTrauma(0.60)
		a.Cam.TriggerHitstop(0.08)
		a.Cam.TriggerFlash(0.28)
		a.SpawnFloatingText(e.X, e.Y-32, "ZEUS'S CHAMPION DEFEATED!", color.RGBA{255, 225, 75, 255}, true)
		// Golden Chest Burst: all 3 Stat Power-Ups + Full Heal + XP shower!
		a.SpawnPickup(PickupAttackSpeed, e.X-28, e.Y, 1.0)
		a.SpawnPickup(PickupAttackDamage, e.X+28, e.Y, 1.0)
		a.SpawnPickup(PickupMoveSpeed, e.X, e.Y-28, 1.0)
		a.SpawnPickup(PickupNectarHeal, e.X, e.Y+28, 1.0)
		for k := 0; k < 8; k++ {
			ang := float64(k) * (math.Pi / 4.0)
			a.SpawnPickup(PickupXP, e.X+math.Cos(ang)*36, e.Y+math.Sin(ang)*36, 20.0)
		}
	case EnemyBossCronus:
		a.CronusSlain = true
		a.Cam.AddTrauma(0.85)
		a.Cam.TriggerHitstop(0.12)
		a.Cam.TriggerFlash(0.40)
		a.Audio.PlaySFX("boon_levelup", 50*time.Millisecond)
		a.SpawnFloatingText(e.X, e.Y-36, "TITAN CRONUS OVERTHROWN!", color.RGBA{255, 225, 85, 255}, true)
	}
}

func (a *Arena) damageHero(dmg float64) {
	h := &a.Hero
	if h.InvulnTimer > 0 || a.HeroDead || a.CronusSlain {
		return
	}
	effectiveDmg := dmg * (1.0 - h.WeaponSpec.DamageReduction)
	h.HP -= effectiveDmg
	h.InvulnTimer = 0.55
	h.HitFlash = 0.25
	a.Audio.PlaySFX("hero_hurt", 80*time.Millisecond)
	a.Cam.AddTrauma(0.32)
	a.SpawnFloatingText(h.X, h.Y-24, fmt.Sprintf("-%.0f", effectiveDmg), color.RGBA{255, 75, 85, 255}, true)
	a.SpawnParticleBurst(h.X, h.Y, 12, color.RGBA{255, 70, 80, 255})

	if h.HP <= 0 {
		h.HP = 0
		if a.DemoMode {
			// In Demo Mode keep the spectacle running
			h.HP = h.MaxHP
			return
		}
		a.HeroDead = true
		_ = h.Anim.Play("death")
	}
}

func (a *Arena) spawnPlayerProj(x, y, angle, speed, radius, dmg float64, pierce, bounce int, life float64, isWave, isBeam bool, col color.RGBA) {
	for i := range a.PlayerProjs {
		p := &a.PlayerProjs[i]
		if !p.Active {
			*p = PlayerProj{
				Active:         true,
				IsSwordWave:    isWave,
				IsBladeBeam:    isBeam,
				X:              x,
				Y:              y,
				VX:             math.Cos(angle) * speed,
				VY:             math.Sin(angle) * speed,
				Angle:          angle,
				Radius:         radius,
				Damage:         dmg,
				PierceLeft:     pierce,
				BounceLeft:     bounce,
				Life:           life,
				MaxLife:        life,
				Color:          col,
				LastHitEnemyID: -1,
			}
			return
		}
	}
}

func (a *Arena) updatePlayerProjs(dt float64) {
	for i := range a.PlayerProjs {
		p := &a.PlayerProjs[i]
		if !p.Active {
			continue
		}
		p.Life -= dt
		if p.Life <= 0 {
			p.Active = false
			continue
		}
		p.X += p.VX * dt
		p.Y += p.VY * dt

		// Wall bounce (Poseidon Ricochet) or deactivation
		bouncedWall := false
		if p.X < ArenaMinX || p.X > ArenaMaxX {
			if p.BounceLeft > 0 {
				p.BounceLeft--
				p.VX = -p.VX
				p.X = math.Max(ArenaMinX, math.Min(ArenaMaxX, p.X))
				bouncedWall = true
			} else {
				p.Active = false
				continue
			}
		}
		if p.Y < ArenaMinY || p.Y > ArenaMaxY {
			if p.BounceLeft > 0 {
				p.BounceLeft--
				p.VY = -p.VY
				p.Y = math.Max(ArenaMinY, math.Min(ArenaMaxY, p.Y))
				bouncedWall = true
			} else {
				p.Active = false
				continue
			}
		}
		if bouncedWall {
			p.Angle = math.Atan2(p.VY, p.VX)
			p.LastHitEnemyID = -1
			a.SpawnParticleBurst(p.X, p.Y, 4, p.Color)
		}

		// Pillar collision
		for k := range a.Pillars {
			pil := &a.Pillars[k]
			dx := p.X - pil.X
			dy := p.Y - pil.Y
			dist := math.Hypot(dx, dy)
			if dist < p.Radius+pil.Radius-4 {
				if p.BounceLeft > 0 && dist > 1e-3 {
					p.BounceLeft--
					nx, ny := dx/dist, dy/dist
					dot := p.VX*nx + p.VY*ny
					p.VX -= 2 * dot * nx
					p.VY -= 2 * dot * ny
					p.Angle = math.Atan2(p.VY, p.VX)
					p.X = pil.X + nx*(p.Radius+pil.Radius)
					p.Y = pil.Y + ny*(p.Radius+pil.Radius)
					p.LastHitEnemyID = -1
				} else if !p.IsBladeBeam {
					p.Active = false
				}
				break
			}
		}
		if !p.Active {
			continue
		}

		// Hit Golden Chests
		for cIdx := range a.Chests {
			c := &a.Chests[cIdx]
			if !c.Active || c.Opening {
				continue
			}
			if math.Hypot(p.X-c.X, p.Y-c.Y) <= p.Radius+20.0 {
				a.damageChest(c, p.Damage)
				if p.PierceLeft > 0 {
					p.PierceLeft--
				} else {
					p.Active = false
					break
				}
			}
		}
		if !p.Active {
			continue
		}

		// Hit Enemies
		for eIdx := range a.Enemies {
			e := &a.Enemies[eIdx]
			if !e.Active || e.Dying || e.ID == p.LastHitEnemyID {
				continue
			}
			if math.Hypot(p.X-e.X, p.Y-e.Y) <= p.Radius+e.Radius {
				a.damageEnemy(e, p.Damage, p.VX, p.VY, p.IsSwordWave)
				p.LastHitEnemyID = e.ID

				if p.PierceLeft > 0 {
					p.PierceLeft--
				} else if p.BounceLeft > 0 {
					p.BounceLeft--
					if !a.ricochetToNextEnemy(p, e) {
						p.VX = -p.VX
						p.VY = -p.VY
						p.Angle = math.Atan2(p.VY, p.VX)
					}
				} else {
					p.Active = false
				}
				break
			}
		}
	}
}

func (a *Arena) ricochetToNextEnemy(p *PlayerProj, from *Enemy) bool {
	bestDist := 280.0
	var next *Enemy
	for i := range a.Enemies {
		cand := &a.Enemies[i]
		if !cand.Active || cand.Dying || cand.ID == from.ID {
			continue
		}
		d := math.Hypot(cand.X-from.X, cand.Y-from.Y)
		if d < bestDist {
			bestDist = d
			next = cand
		}
	}
	if next == nil {
		return false
	}
	spd := math.Hypot(p.VX, p.VY)
	ang := math.Atan2(next.Y-p.Y, next.X-p.X)
	p.VX = math.Cos(ang) * spd
	p.VY = math.Sin(ang) * spd
	p.Angle = ang
	return true
}

func (a *Arena) spawnEnemyProj(x, y, angle, speed, radius, dmg, life float64, isCrescent bool, col color.RGBA) {
	for i := range a.EnemyProjs {
		ep := &a.EnemyProjs[i]
		if !ep.Active {
			*ep = EnemyProj{
				Active:     true,
				IsCrescent: isCrescent,
				X:          x,
				Y:          y,
				VX:         math.Cos(angle) * speed,
				VY:         math.Sin(angle) * speed,
				Angle:      angle,
				Radius:     radius,
				Damage:     dmg,
				Life:       life,
				Color:      col,
			}
			return
		}
	}
}

func (a *Arena) updateEnemyProjs(dt float64) {
	h := &a.Hero
	for i := range a.EnemyProjs {
		ep := &a.EnemyProjs[i]
		if !ep.Active {
			continue
		}
		ep.Life -= dt
		if ep.Life <= 0 {
			ep.Active = false
			continue
		}
		ep.X += ep.VX * dt
		ep.Y += ep.VY * dt
		if ep.X < ArenaMinX || ep.X > ArenaMaxX || ep.Y < ArenaMinY || ep.Y > ArenaMaxY {
			ep.Active = false
			continue
		}
		for k := range a.Pillars {
			pil := &a.Pillars[k]
			if math.Hypot(ep.X-pil.X, ep.Y-pil.Y) < ep.Radius+pil.Radius-4 {
				ep.Active = false
				a.SpawnParticleBurst(ep.X, ep.Y, 5, ep.Color)
				break
			}
		}
		if !ep.Active {
			continue
		}
		if math.Hypot(ep.X-h.X, ep.Y-h.Y) <= ep.Radius+h.Radius-2 {
			ep.Active = false
			a.damageHero(ep.Damage)
		}
	}
}

// SpawnPickup drops an Ambrosia XP gem or Divine Power-Up Relic onto the floor.
func (a *Arena) SpawnPickup(kind PickupKind, x, y, value float64) {
	x, y = a.clampAndResolvePillars(x, y, 12.0)
	for i := range a.Pickups {
		p := &a.Pickups[i]
		if !p.Active {
			spriteName := "pickup_xp"
			if kind != PickupXP {
				spriteName = "powerup_relic"
			}
			ctrl, err := a.Reg.NewSpriteController(spriteName)
			if err != nil {
				return
			}
			*p = Pickup{
				Active: true,
				Kind:   kind,
				X:      x,
				Y:      y,
				Value:  value,
				Anim:   ctrl,
			}
			return
		}
	}
}

func (a *Arena) updatePickups(dt float64, dtDur time.Duration) {
	h := &a.Hero
	for i := range a.Pickups {
		p := &a.Pickups[i]
		if !p.Active {
			continue
		}
		p.Age += dt
		p.Anim.Update(dtDur)

		dx := h.X - p.X
		dy := h.Y - p.Y
		dist := math.Hypot(dx, dy)

		magnetRadius := 145.0
		if p.Kind != PickupXP {
			magnetRadius = 110.0
		}
		if dist < magnetRadius && dist > 1e-3 {
			pullSpd := 380.0
			p.X += (dx / dist) * pullSpd * dt
			p.Y += (dy / dist) * pullSpd * dt
		}

		if dist <= h.Radius+16.0 {
			p.Active = false
			a.collectPickup(p)
		}
	}
}

func (a *Arena) collectPickup(p *Pickup) {
	h := &a.Hero
	switch p.Kind {
	case PickupXP:
		a.Audio.PlaySFX("xp_ambrosia", 35*time.Millisecond)
		h.XP += p.Value * h.Stats.XPMul
		for h.XP >= h.NextLevelXP {
			h.XP -= h.NextLevelXP
			h.Level++
			h.NextLevelXP = math.Round(h.NextLevelXP*1.28 + 15.0)
			h.PendingLevelUps++
			a.Audio.PlaySFX("boon_levelup", 80*time.Millisecond)
		}
	case PickupAttackSpeed:
		a.Audio.PlaySFX("powerup_pickup", 50*time.Millisecond)
		h.PermAttackSpeed += 0.15
		h.OverdriveSpd = 6.0
		h.PowerUpsPicked++
		a.RecomputeHeroStats()
		a.SpawnFloatingText(h.X, h.Y-28, "APOLLO'S QUIVER: +15% ATK SPD & OVERDRIVE!", color.RGBA{255, 220, 65, 255}, true)
		a.SpawnParticleBurst(h.X, h.Y, 16, color.RGBA{255, 220, 65, 255})
	case PickupAttackDamage:
		a.Audio.PlaySFX("powerup_pickup", 50*time.Millisecond)
		h.PermDamage += 0.15
		h.OverdriveDmg = 6.0
		h.PowerUpsPicked++
		a.RecomputeHeroStats()
		a.SpawnFloatingText(h.X, h.Y-28, "ARES'S EMBER: +15% DAMAGE & BLOODFIRE!", color.RGBA{255, 85, 75, 255}, true)
		a.SpawnParticleBurst(h.X, h.Y, 16, color.RGBA{255, 85, 75, 255})
	case PickupMoveSpeed:
		a.Audio.PlaySFX("powerup_pickup", 50*time.Millisecond)
		h.PermMoveSpeed += 0.12
		h.OverdriveMove = 6.0
		h.PowerUpsPicked++
		a.RecomputeHeroStats()
		a.SpawnFloatingText(h.X, h.Y-28, "HERMES'S PLUME: +12% SPEED & ZEPHYR HASTE!", color.RGBA{80, 220, 255, 255}, true)
		a.SpawnParticleBurst(h.X, h.Y, 16, color.RGBA{80, 220, 255, 255})
	case PickupNectarHeal:
		a.Audio.PlaySFX("powerup_pickup", 50*time.Millisecond)
		heal := 35.0
		h.HP = math.Min(h.MaxHP, h.HP+heal)
		h.PowerUpsPicked++
		a.SpawnFloatingText(h.X, h.Y-28, "HEBE'S NECTAR: +35 HP!", color.RGBA{110, 245, 145, 255}, true)
		a.SpawnParticleBurst(h.X, h.Y, 14, color.RGBA{110, 245, 145, 255})
	}
}

func (a *Arena) spawnSlashArc(x, y, angle, sweepRad, radius float64, col color.RGBA, whirlwind bool) {
	for i := range a.SlashArcs {
		sa := &a.SlashArcs[i]
		if !sa.Active {
			*sa = SlashArc{
				Active:    true,
				X:         x,
				Y:         y,
				Angle:     angle,
				SweepRad:  sweepRad,
				Radius:    radius,
				Timer:     0.18,
				MaxTimer:  0.18,
				Color:     col,
				Whirlwind: whirlwind,
			}
			return
		}
	}
}

func (a *Arena) spawnLightningArc(x1, y1, x2, y2 float64, col color.RGBA) {
	for i := range a.LightningArcs {
		la := &a.LightningArcs[i]
		if !la.Active {
			*la = LightningArc{
				Active: true,
				X1:     x1,
				Y1:     y1,
				X2:     x2,
				Y2:     y2,
				Timer:  0.18,
				Color:  col,
			}
			return
		}
	}
}

// SpawnFloatingText adds a world-space floating damage or status label.
func (a *Arena) SpawnFloatingText(x, y float64, msg string, col color.RGBA, big bool) {
	for i := range a.FloatingTexts {
		ft := &a.FloatingTexts[i]
		if !ft.Active {
			dur := 0.75
			if big {
				dur = 1.15
			}
			*ft = FloatingText{
				Active:   true,
				X:        x + (a.RNG.Float64()*14 - 7),
				Y:        y,
				VY:       -42.0,
				Text:     msg,
				Color:    col,
				Timer:    dur,
				MaxTimer: dur,
				Big:      big,
			}
			return
		}
	}
}

// SpawnParticle activates a single particle in the pool.
func (a *Arena) SpawnParticle(x, y, vx, vy, size, life float64, col color.RGBA) {
	for i := range a.Particles {
		pt := &a.Particles[i]
		if !pt.Active {
			*pt = Particle{
				Active:   true,
				X:        x,
				Y:        y,
				VX:       vx,
				VY:       vy,
				Size:     size,
				Timer:    life,
				MaxTimer: life,
				Color:    col,
			}
			return
		}
	}
}

// SpawnParticleBurst spawns count radial particles at (x, y).
func (a *Arena) SpawnParticleBurst(x, y float64, count int, col color.RGBA) {
	for k := 0; k < count; k++ {
		ang := a.RNG.Float64() * math.Pi * 2.0
		spd := 45.0 + a.RNG.Float64()*140.0
		life := 0.18 + a.RNG.Float64()*0.24
		sz := 3.0 + a.RNG.Float64()*3.5
		a.SpawnParticle(x, y, math.Cos(ang)*spd, math.Sin(ang)*spd, sz, life, col)
	}
}

func (a *Arena) updateVisualEffects(dt float64) {
	if a.TimeRift.Active {
		a.TimeRift.Timer -= dt
		if a.TimeRift.Timer <= 0 {
			a.TimeRift.Active = false
		}
	}
	for i := range a.SlashArcs {
		if a.SlashArcs[i].Active {
			a.SlashArcs[i].Timer -= dt
			if a.SlashArcs[i].Timer <= 0 {
				a.SlashArcs[i].Active = false
			}
		}
	}
	for i := range a.LightningArcs {
		if a.LightningArcs[i].Active {
			a.LightningArcs[i].Timer -= dt
			if a.LightningArcs[i].Timer <= 0 {
				a.LightningArcs[i].Active = false
			}
		}
	}
	for i := range a.FloatingTexts {
		ft := &a.FloatingTexts[i]
		if ft.Active {
			ft.Timer -= dt
			ft.Y += ft.VY * dt
			if ft.Timer <= 0 {
				ft.Active = false
			}
		}
	}
	for i := range a.Particles {
		pt := &a.Particles[i]
		if pt.Active {
			pt.Timer -= dt
			pt.X += pt.VX * dt
			pt.Y += pt.VY * dt
			if pt.Timer <= 0 {
				pt.Active = false
			}
		}
	}
}

func (a *Arena) clampAndResolvePillars(x, y, radius float64) (float64, float64) {
	x = math.Max(ArenaMinX+radius, math.Min(ArenaMaxX-radius, x))
	y = math.Max(ArenaMinY+radius, math.Min(ArenaMaxY-radius, y))

	for i := range a.Pillars {
		p := &a.Pillars[i]
		dx := x - p.X
		dy := y - p.Y
		dist := math.Hypot(dx, dy)
		minD := radius + p.Radius
		if dist < minD {
			if dist < 1e-4 {
				x = p.X + minD
			} else {
				x = p.X + (dx/dist)*minD
				y = p.Y + (dy/dist)*minD
			}
		}
	}
	return x, y
}

func (a *Arena) countActiveEnemies() int {
	n := 0
	for i := range a.Enemies {
		if a.Enemies[i].Active && !a.Enemies[i].Dying {
			n++
		}
	}
	return n
}

func (a *Arena) countActiveChests() int {
	n := 0
	for i := range a.Chests {
		if a.Chests[i].Active && !a.Chests[i].Opening {
			n++
		}
	}
	return n
}

// ActiveBoss returns a pointer to the currently active Boss (Cronus or Zeus's Champion), or nil.
func (a *Arena) ActiveBoss() *Enemy {
	for i := range a.Enemies {
		e := &a.Enemies[i]
		if e.Active && !e.Dying && (e.Kind == EnemyBossCronus || e.Kind == EnemyBossZeusChampion) {
			return e
		}
	}
	return nil
}

// ActiveEnemyCount returns the number of living enemies in the courtyard.
func (a *Arena) ActiveEnemyCount() int {
	return a.countActiveEnemies()
}

func normalizeAngle(rad float64) float64 {
	for rad > math.Pi {
		rad -= 2 * math.Pi
	}
	for rad < -math.Pi {
		rad += 2 * math.Pi
	}
	return rad
}
