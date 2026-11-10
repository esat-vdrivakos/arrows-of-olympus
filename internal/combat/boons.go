package combat

import (
	"image/color"
	"math/rand"
)

// WeaponType distinguishes the start-of-run Hero armament.
type WeaponType int

const (
	WeaponBow WeaponType = iota
	WeaponSword
)

// WeaponSpec describes base stats and flavor for a playable Olympus weapon.
type WeaponSpec struct {
	Type            WeaponType
	Name            string
	Subtitle        string
	Playstyle       string
	SpecialTrait    string
	BaseDamage      float64
	Cooldown        float64
	Range           float64
	MoveSpeedBonus  float64
	BonusMaxHP      float64
	DamageReduction float64
	AccentColor     color.RGBA
}

var (
	BowOfApolloSpec = WeaponSpec{
		Type:            WeaponBow,
		Name:            "GOLDEN BOW OF APOLLO",
		Subtitle:        "Sacred Ranged Auto-Shooter",
		Playstyle:       "Long-range sun arrows (520 px) with rapid volleys & +10% Move Speed.",
		SpecialTrait:    "Stacks multi-shot, diagonal cones, piercing & wall ricochets.",
		BaseDamage:      28.0,
		Cooldown:        0.42,
		Range:           520.0,
		MoveSpeedBonus:  0.10,
		BonusMaxHP:      0.0,
		DamageReduction: 0.0,
		AccentColor:     color.RGBA{235, 185, 55, 255},
	}

	SwordOfAresSpec = WeaponSpec{
		Type:            WeaponSword,
		Name:            "BRONZE SWORD OF ARES",
		Subtitle:        "Sacred Melee Cleave & Slash-Wave",
		Playstyle:       "Wide 110° bronze cleave arc + crescent slash-wave, +35 Max HP & 15% Armor.",
		SpecialTrait:    "SHATTERS Shielded Hoplite blocks for 100% damage & cleaves packs!",
		BaseDamage:      48.0,
		Cooldown:        0.38,
		Range:           155.0,
		MoveSpeedBonus:  0.0,
		BonusMaxHP:      35.0,
		DamageReduction: 0.15,
		AccentColor:     color.RGBA{225, 65, 75, 255},
	}
)

// GetWeaponSpec returns the WeaponSpec for the given WeaponType.
func GetWeaponSpec(w WeaponType) WeaponSpec {
	if w == WeaponSword {
		return SwordOfAresSpec
	}
	return BowOfApolloSpec
}

// BoonID identifies each of the 12 Olympian Gods' upgrades.
type BoonID int

const (
	BoonArtemis BoonID = iota
	BoonApollo
	BoonHermes
	BoonHephaestus
	BoonPoseidon
	BoonZeus
	BoonAthena
	BoonAres
	BoonHades
	BoonDemeter
	BoonDionysus
	BoonHera
	BoonCount
)

// BoonDef defines a single Olympian God Boon and its weapon-specific effects.
type BoonDef struct {
	ID          BoonID
	GodName     string
	GodTitle    string
	BoonTitle   string
	Symbol      string
	BowDesc     string
	SwordDesc   string
	AccentColor color.RGBA
}

// DescriptionFor returns the dynamic description matching the player's equipped weapon.
func (b BoonDef) DescriptionFor(w WeaponType) string {
	if w == WeaponSword {
		return b.SwordDesc
	}
	return b.BowDesc
}

// AllBoons lists all 12 Gods of Olympus Boons in canonical order.
var AllBoons = [BoonCount]BoonDef{
	BoonArtemis: {
		ID:          BoonArtemis,
		GodName:     "ARTEMIS",
		GodTitle:    "Goddess of the Hunt",
		BoonTitle:   "Twin Strike",
		Symbol:      "I",
		BowDesc:     "+1 Parallel Front Arrow per volley & +12% Damage.",
		SwordDesc:   "Unleashes a Double-Slash Cleave per swing & +15% Damage.",
		AccentColor: color.RGBA{75, 210, 115, 255},
	},
	BoonApollo: {
		ID:          BoonApollo,
		GodName:     "APOLLO",
		GodTitle:    "God of Sun & Archery",
		BoonTitle:   "Sunburst Cone",
		Symbol:      "II",
		BowDesc:     "+2 Diagonal Sun Arrows at ±20° on every shot.",
		SwordDesc:   "+2 Diagonal Golden Slash-Waves at ±28° on every swing.",
		AccentColor: color.RGBA{245, 195, 55, 255},
	},
	BoonHermes: {
		ID:          BoonHermes,
		GodName:     "HERMES",
		GodTitle:    "Messenger of Olympus",
		BoonTitle:   "Winged Sandals",
		Symbol:      "III",
		BowDesc:     "+25% Attack Speed & +15% Movement Speed.",
		SwordDesc:   "+25% Swing Speed & +15% Movement Speed.",
		AccentColor: color.RGBA{70, 200, 245, 255},
	},
	BoonHephaestus: {
		ID:          BoonHephaestus,
		GodName:     "HEPHAESTUS",
		GodTitle:    "God of the Forge",
		BoonTitle:   "Adamant Edge",
		Symbol:      "IV",
		BowDesc:     "Arrows PIERCE through enemies (+1 Pierce) & +25% Damage.",
		SwordDesc:   "+35% Cleave Radius & launches a long-range Piercing Blade Beam!",
		AccentColor: color.RGBA{240, 115, 45, 255},
	},
	BoonPoseidon: {
		ID:          BoonPoseidon,
		GodName:     "POSEIDON",
		GodTitle:    "Lord of the Aegean",
		BoonTitle:   "Tidal Force",
		Symbol:      "V",
		BowDesc:     "Arrows RICOCHET off walls/enemies (+1 Bounce) + Knockback.",
		SwordDesc:   "Sweeps launch Bouncing Tidal Shockwaves + Heavy Knockback.",
		AccentColor: color.RGBA{55, 145, 235, 255},
	},
	BoonZeus: {
		ID:          BoonZeus,
		GodName:     "ZEUS",
		GodTitle:    "King of Olympus",
		BoonTitle:   "Thunderbolt Strike",
		Symbol:      "VI",
		BowDesc:     "Arrow hits spark Chain Lightning arcs across nearby foes.",
		SwordDesc:   "Sword cleaves unleash Chain Lightning arcs across nearby foes.",
		AccentColor: color.RGBA{130, 195, 255, 255},
	},
	BoonAthena: {
		ID:          BoonAthena,
		GodName:     "ATHENA",
		GodTitle:    "Goddess of Wisdom",
		BoonTitle:   "Aegis Barrier",
		Symbol:      "VII",
		BowDesc:     "Summons Orbiting Golden Shields that block shots & sear foes.",
		SwordDesc:   "Summons Orbiting Golden Shields that block shots & crush foes.",
		AccentColor: color.RGBA{230, 205, 110, 255},
	},
	BoonAres: {
		ID:          BoonAres,
		GodName:     "ARES",
		GodTitle:    "God of War",
		BoonTitle:   "Bloodlust Crit",
		Symbol:      "VIII",
		BowDesc:     "+25% Critical Hit Chance for 2.5x Golden Crit Damage.",
		SwordDesc:   "+25% Critical Hit Chance for 2.5x Brutal Cleave Damage.",
		AccentColor: color.RGBA{225, 55, 70, 255},
	},
	BoonHades: {
		ID:          BoonHades,
		GodName:     "HADES",
		GodTitle:    "Lord of the Underworld",
		BoonTitle:   "Soul Siphon",
		Symbol:      "IX",
		BowDesc:     "Slaying enemies siphons +4 HP & boosts Damage by +12%.",
		SwordDesc:   "Slaying enemies siphons +5 HP & boosts Damage by +12%.",
		AccentColor: color.RGBA{165, 85, 225, 255},
	},
	BoonDemeter: {
		ID:          BoonDemeter,
		GodName:     "DEMETER",
		GodTitle:    "Goddess of the Harvest",
		BoonTitle:   "Golden Ambrosia",
		Symbol:      "X",
		BowDesc:     "+40 Max HP immediately & restores +50 HP.",
		SwordDesc:   "+45 Max HP immediately & restores +60 HP.",
		AccentColor: color.RGBA{135, 225, 95, 255},
	},
	BoonDionysus: {
		ID:          BoonDionysus,
		GodName:     "DIONYSUS",
		GodTitle:    "God of Revelry & Wine",
		BoonTitle:   "Intoxicating Miasma",
		Symbol:      "XI",
		BowDesc:     "Arrows apply Venomous Miasma dealing stacking Poison DPS.",
		SwordDesc:   "Sword cleaves drench foes in Miasma dealing stacking Poison DPS.",
		AccentColor: color.RGBA{205, 80, 195, 255},
	},
	BoonHera: {
		ID:          BoonHera,
		GodName:     "HERA",
		GodTitle:    "Queen of Heaven",
		BoonTitle:   "Sovereign Rearguard",
		Symbol:      "XII",
		BowDesc:     "Fires Rearguard Arrows behind you & grants +25% XP Gain.",
		SwordDesc:   "Unlocks 360° Whirlwind Back-Slash & grants +25% XP Gain.",
		AccentColor: color.RGBA{115, 215, 210, 255},
	},
}

// ComputedStats aggregates base weapon parameters, stacked Olympian Boons, and permanent/temporary orb buffs.
type ComputedStats struct {
	FrontProjectiles int
	DiagonalPairs    int
	RearAttack       int
	PierceCount      int
	RicochetCount    int
	ChainLightning   int
	OrbitShields     int
	CritChance       float64
	LifeOnKill       float64
	PoisonDPS        float64
	DamageMul        float64
	AttackSpeedMul   float64
	MoveSpeedMul     float64
	CleaveRadiusMul  float64
	XPMul            float64
}

// ComputeStats calculates the player's effective combat multipliers from owned Boons and orb stacks.
func ComputeStats(weapon WeaponType, boons [BoonCount]int, permAtkSpd, permDmg, permMoveSpd float64, overdriveSpd, overdriveDmg, overdriveMove bool) ComputedStats {
	spec := GetWeaponSpec(weapon)
	st := ComputedStats{
		FrontProjectiles: 1 + boons[BoonArtemis],
		DiagonalPairs:    boons[BoonApollo],
		RearAttack:       boons[BoonHera],
		PierceCount:      boons[BoonHephaestus],
		RicochetCount:    boons[BoonPoseidon],
		ChainLightning:   boons[BoonZeus],
		OrbitShields:     boons[BoonAthena] * 2,
		CritChance:       0.06 + float64(boons[BoonAres])*0.25,
		LifeOnKill:       float64(boons[BoonHades]) * 4.0,
		PoisonDPS:        float64(boons[BoonDionysus]) * 18.0,
		DamageMul:        1.0 + permDmg + float64(boons[BoonArtemis])*0.12 + float64(boons[BoonHephaestus])*0.25 + float64(boons[BoonHades])*0.12,
		AttackSpeedMul:   1.0 + permAtkSpd + float64(boons[BoonHermes])*0.25,
		MoveSpeedMul:     1.0 + spec.MoveSpeedBonus + permMoveSpd + float64(boons[BoonHermes])*0.15,
		CleaveRadiusMul:  1.0 + float64(boons[BoonHephaestus])*0.35,
		XPMul:            1.0 + float64(boons[BoonHera])*0.25,
	}
	if st.CritChance > 0.85 {
		st.CritChance = 0.85
	}
	if overdriveSpd {
		st.AttackSpeedMul += 0.40
	}
	if overdriveDmg {
		st.DamageMul += 0.35
	}
	if overdriveMove {
		st.MoveSpeedMul += 0.35
	}
	return st
}

// DraftBoons picks 3 distinct Olympian Boons for the level-up selection overlay.
func DraftBoons(rng *rand.Rand) [3]BoonDef {
	indices := [BoonCount]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	for i := BoonCount - 1; i > 0; i-- {
		j := rng.Intn(int(i) + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}
	return [3]BoonDef{
		AllBoons[indices[0]],
		AllBoons[indices[1]],
		AllBoons[indices[2]],
	}
}
