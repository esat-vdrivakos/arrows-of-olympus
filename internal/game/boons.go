package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
)

// WeaponType distinguishes between Ranged (Bow) and Melee (Sword) playstyles chosen at run start.
type WeaponType int

const (
	WeaponBow WeaponType = iota
	WeaponSword
)

func (w WeaponType) String() string {
	if w == WeaponSword {
		return "Bronze Sword of Ares (Melee Cleave)"
	}
	return "Golden Bow of Apollo (Ranged Volley)"
}

// UpgradeID identifies each of the 12 Olympian God Power-Ups / Boons.
type UpgradeID int

const (
	UpgradeAttackSpeed UpgradeID = iota // Hermes: Exponential Attack Speed
	UpgradeDamage                       // Ares: Exponential Attack Damage
	UpgradeMoveSpeed                    // Zephyr/Hermes: Exponential Movement Speed
	UpgradeMultiStrike                  // Artemis: Parallel Arrows / Double Sword Cleaves
	UpgradeDiagonalCone                 // Apollo: Diagonal Arrows / Crescent Slash-Waves
	UpgradeRicochetWave                 // Poseidon: Bouncing Ricochet / Tidal Knockback
	UpgradePiercingEdge                 // Hephaestus: Piercing Arrows / Expanding Cleave Size
	UpgradeChainBolt                    // Zeus: Chain Lightning Bolts on Hit
	UpgradeAegisOrbit                   // Athena: Orbiting Golden Aegis Shields
	UpgradeSoulVigor                    // Hades & Demeter: Exponential Max HP, Heal & Lifesteal
	UpgradeVenomMiasma                  // Dionysus: Exponential Poison Burn DoT
	UpgradeRoyalSurge                   // Hera: Rear Shot / Whirlwind + Exponential XP & Crit
	UpgradeCount
)

// UpgradeDef describes a power-up option shown in the 3-choice Olympian Modal.
type UpgradeDef struct {
	ID          UpgradeID
	GodName     string
	Title       string
	StatTag     string
	BaseFactor  float64 // Exponential base b in Multiplier(L) = b^L
	AccentColor color.NRGBA
}

// AllUpgrades defines the 12 Gods of Olympus Power-Ups that appear in the 3-choice modal.
var AllUpgrades = [UpgradeCount]UpgradeDef{
	UpgradeAttackSpeed: {
		ID:          UpgradeAttackSpeed,
		GodName:     "HERMES",
		Title:       "Alacrity of Olympus",
		StatTag:     "ATTACK SPEED",
		BaseFactor:  1.38,
		AccentColor: color.NRGBA{255, 215, 64, 255},
	},
	UpgradeDamage: {
		ID:          UpgradeDamage,
		GodName:     "ARES",
		Title:       "Wrath of the War God",
		StatTag:     "ATTACK DAMAGE",
		BaseFactor:  1.45,
		AccentColor: color.NRGBA{245, 70, 85, 255},
	},
	UpgradeMoveSpeed: {
		ID:          UpgradeMoveSpeed,
		GodName:     "ZEPHYR & HERMES",
		Title:       "Winged Sandals",
		StatTag:     "MOVEMENT SPEED",
		BaseFactor:  1.26,
		AccentColor: color.NRGBA{80, 215, 255, 255},
	},
	UpgradeMultiStrike: {
		ID:          UpgradeMultiStrike,
		GodName:     "ARTEMIS",
		Title:       "Twin Huntress Volley",
		StatTag:     "MULTI-ATTACK",
		BaseFactor:  1.50,
		AccentColor: color.NRGBA{110, 240, 150, 255},
	},
	UpgradeDiagonalCone: {
		ID:          UpgradeDiagonalCone,
		GodName:     "APOLLO",
		Title:       "Sunburst Spread",
		StatTag:     "DIAGONAL VOLLEY",
		BaseFactor:  1.48,
		AccentColor: color.NRGBA{255, 190, 50, 255},
	},
	UpgradeRicochetWave: {
		ID:          UpgradeRicochetWave,
		GodName:     "POSEIDON",
		Title:       "Tidal Ricochet",
		StatTag:     "BOUNCE & KNOCKBACK",
		BaseFactor:  1.46,
		AccentColor: color.NRGBA{65, 165, 255, 255},
	},
	UpgradePiercingEdge: {
		ID:          UpgradePiercingEdge,
		GodName:     "HEPHAESTUS",
		Title:       "Adamant Forging",
		StatTag:     "PIERCE & CLEAVE SIZE",
		BaseFactor:  1.42,
		AccentColor: color.NRGBA{255, 135, 60, 255},
	},
	UpgradeChainBolt: {
		ID:          UpgradeChainBolt,
		GodName:     "ZEUS",
		Title:       "Thunderbolt Arc",
		StatTag:     "CHAIN LIGHTNING",
		BaseFactor:  1.55,
		AccentColor: color.NRGBA{130, 220, 255, 255},
	},
	UpgradeAegisOrbit: {
		ID:          UpgradeAegisOrbit,
		GodName:     "ATHENA",
		Title:       "Golden Aegis Orbit",
		StatTag:     "ORBITAL SHIELDS",
		BaseFactor:  1.50,
		AccentColor: color.NRGBA{255, 230, 120, 255},
	},
	UpgradeSoulVigor: {
		ID:          UpgradeSoulVigor,
		GodName:     "DEMETER & HADES",
		Title:       "Ambrosia Vitality",
		StatTag:     "MAX HP & LIFESTEAL",
		BaseFactor:  1.40,
		AccentColor: color.NRGBA{135, 245, 135, 255},
	},
	UpgradeVenomMiasma: {
		ID:          UpgradeVenomMiasma,
		GodName:     "DIONYSUS",
		Title:       "Nectar Venom",
		StatTag:     "EXPONENTIAL DOT",
		BaseFactor:  1.58,
		AccentColor: color.NRGBA{210, 115, 255, 255},
	},
	UpgradeRoyalSurge: {
		ID:          UpgradeRoyalSurge,
		GodName:     "HERA",
		Title:       "Sovereign Crown",
		StatTag:     "REAR ATTACK & CRIT",
		BaseFactor:  1.44,
		AccentColor: color.NRGBA{255, 150, 215, 255},
	},
}

// ExponentialMultiplier computes b^level (for level == 0 returns 1.0).
// Every level-up compounds exponentially: e.g. 1.45^1 = 1.45x, 1.45^2 = 2.10x, 1.45^3 = 3.05x, 1.45^4 = 4.42x.
func ExponentialMultiplier(id UpgradeID, level int) float64 {
	if level <= 0 {
		return 1.0
	}
	base := AllUpgrades[id].BaseFactor
	if base <= 1.0 {
		base = 1.35
	}
	return math.Pow(base, float64(level))
}

// FormatUpgradeDescription returns a weapon-aware description showing the exponential jump from currentLevel to currentLevel+1.
func FormatUpgradeDescription(id UpgradeID, currentLevel int, weapon WeaponType) (string, string) {
	nextLevel := currentLevel + 1
	curMult := ExponentialMultiplier(id, currentLevel)
	nextMult := ExponentialMultiplier(id, nextLevel)
	expBadge := fmt.Sprintf("Lv.%d -> Lv.%d  (x%.2f -> x%.2f EXP!)", currentLevel, nextLevel, curMult, nextMult)

	var desc string
	switch id {
	case UpgradeAttackSpeed:
		desc = fmt.Sprintf("Exponentially accelerates fire/swing rate to x%.2f speed!", nextMult)
	case UpgradeDamage:
		desc = fmt.Sprintf("Exponentially amplifies all weapon & spell damage to x%.2f!", nextMult)
	case UpgradeMoveSpeed:
		desc = fmt.Sprintf("Exponentially boosts Hero speed to x%.2f + reduces Dash cooldown!", nextMult)
	case UpgradeMultiStrike:
		if weapon == WeaponBow {
			desc = fmt.Sprintf("Fires +%d parallel front arrows per volley with x%.2f volley power!", nextLevel, nextMult)
		} else {
			desc = fmt.Sprintf("Unleashes +%d rapid multi-cleave sword strikes per swing (x%.2f power)!", nextLevel, nextMult)
		}
	case UpgradeDiagonalCone:
		if weapon == WeaponBow {
			desc = fmt.Sprintf("Fires +%d diagonal sunburst arrows at wide angles (x%.2f cone dmg)!", nextLevel*2, nextMult)
		} else {
			desc = fmt.Sprintf("Launches +%d golden crescent slash-waves outward (x%.2f wave dmg)!", nextLevel, nextMult)
		}
	case UpgradeRicochetWave:
		desc = fmt.Sprintf("Attacks bounce +%d times between foes/walls with x%.2f tidal impact!", nextLevel, nextMult)
	case UpgradePiercingEdge:
		if weapon == WeaponBow {
			desc = fmt.Sprintf("Arrows pierce +%d enemies & ignore Hoplite shields (x%.2f velocity)!", nextLevel, nextMult)
		} else {
			desc = fmt.Sprintf("Expands Sword cleave arc & range by x%.2f and crushes armor!", nextMult)
		}
	case UpgradeChainBolt:
		desc = fmt.Sprintf("Hits arc Zeus lightning to %d nearby foes for x%.2f storm damage!", nextLevel+1, nextMult)
	case UpgradeAegisOrbit:
		desc = fmt.Sprintf("Summons %d orbiting Golden Aegis shields dealing x%.2f contact dmg!", nextLevel, nextMult)
	case UpgradeSoulVigor:
		desc = fmt.Sprintf("Heals 40%% HP immediately, scales Max HP by x%.2f & adds lifesteal!", nextMult)
	case UpgradeVenomMiasma:
		desc = fmt.Sprintf("Inflicts Dionysian poison burning for %.0f DPS (x%.2f exponential DoT)!", 18.0*nextMult, nextMult)
	case UpgradeRoyalSurge:
		if weapon == WeaponBow {
			desc = fmt.Sprintf("Fires rear guard arrows + boosts Crit Damage & XP gain by x%.2f!", nextMult)
		} else {
			desc = fmt.Sprintf("Adds 360-degree whirlwind back-slash + x%.2f Crit & XP gain!", nextMult)
		}
	}
	return expBadge, desc
}

// RollThreeUpgrades selects 3 distinct random upgrades from the 12 Olympian pool,
// guaranteeing at least one core stat power-up (Attack Speed, Damage, or Movement Speed) in every roll.
func RollThreeUpgrades(rng *rand.Rand) [3]UpgradeID {
	coreStats := []UpgradeID{UpgradeAttackSpeed, UpgradeDamage, UpgradeMoveSpeed}
	first := coreStats[rng.IntN(len(coreStats))]

	pool := make([]UpgradeID, 0, UpgradeCount-1)
	for i := UpgradeID(0); i < UpgradeCount; i++ {
		if i != first {
			pool = append(pool, i)
		}
	}
	// Fisher-Yates shuffle on remaining pool
	for i := len(pool) - 1; i > 0; i-- {
		j := rng.IntN(i + 1)
		pool[i], pool[j] = pool[j], pool[i]
	}

	out := [3]UpgradeID{first, pool[0], pool[1]}
	// Shuffle the 3 card positions so the core stat isn't always slot 0
	for i := 2; i > 0; i-- {
		j := rng.IntN(i + 1)
		out[i], out[j] = out[j], out[i]
	}
	return out
}
