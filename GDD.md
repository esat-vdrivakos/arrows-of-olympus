# Game Design Document (GDD): **Arrows of Olympus**

> **Game Title**: Arrows of Olympus
> **Target 2D Genre**: Top-Down Auto-Shooter Survival / Roguelite Arena (*Archer.io* / *Archero* + *Vampire Survivors* hybrid)
> **Camera Projection & Spatial Model**: 3/4 Top-Down Orthographic Arena (`1280x720` virtual canvas, bounded Pantheon courtyard with marble pillars)
> **Simulation Model**: Real-Time Continuous (60 TPS fixed timestep in Ebitengine v2)
> **Target Platform**: Desktop (macOS/Windows/Linux) & WebAssembly (Browser / Game Dev Studio Playtest Tab)
> **Virtual Canvas & Aspect Ratio**: 16:9 Widescreen (`1280x720` Native 2D HD)
> **Author / Lead Designer**: User & Antigravity Game Studio

---

## 0. The 30-Minute Grand Execution Plan

To ship a complete, playable, visually rich **Arrows of Olympus** in 30 minutes, we divide execution into **4 focused phases** (leveraging parallel asset creation where possible):

| Phase | Time Budget | Deliverable | Skills & Tools |
| :--- | :---: | :--- | :--- |
| **Phase 1: Sprites & Olympus Arena Background** | **Minutes 0–10** | 1. **Greek Hero** sprite sheets / stances (`idle`, `run`, `bow_attack`, `sword_attack`, `hit`, `death`)<br>2. **Mobs**: **Skeleton Spartoi** (starting mob) & **Greek Hoplite Warriors** (Unshielded & Bronze-Shielded variants)<br>3. **Bosses**: **Zeus's Champion** (mid-boss) & **Titan Cronus** (final boss)<br>4. **Olympus Background**: `1280x720` Mount Olympus marble pantheon courtyard<br>5. **`.sprite` DSL Props**: Golden Arrows, Sword Slash Waves, Lightning bolts, Cronus scythe waves, **Olympian Treasure Chests**, **Power-Up Relics** (Attack Speed, Damage, Movement Speed, Heal), Ambrosia XP gems, Marble pillars | `nano-banana` + `sprite-design` |
| **Phase 2: "12 Gods of Olympus" Music & Combat SFX** | **Minutes 10–15** | 1. **`olympus_12_gods.music`**: Phrygian/Dorian epic Greek soundtrack evoking the 12 Olympians (Kithara/Lyre plucks, Aulos melodies, Bronze Horn brass, War Timpani, anchored sub-bass)<br>2. **Combat `.sfx` Suite**: `bow_shoot.sfx`, `sword_slash.sfx`, `arrow_hit.sfx`, `shield_clang.sfx`, `chest_open.sfx`, `powerup_pickup.sfx`, `xp_ambrosia.sfx`, `boon_levelup.sfx`, `thunder_boss.sfx`, `hero_hurt.sfx` | `procedural-music` + `procedural-sfx` |
| **Phase 3: Core *Archer.io* Ranged & Melee Engine** | **Minutes 15–25** | 1. **Start-of-Run Weapon Select**: Choose between **Bow of Apollo (Ranged Auto-Shooter)** or **Sword of Ares (Melee Cleave & Slash-Wave)**<br>2. **Movement & Auto-Combat**: 8-way movement + continuous auto-attack (ranged arrows vs. wide melee cleave arcs)<br>3. **Map Chests & NPC Drop Power-Ups**: Breakable Golden Olympus Chests + monster drops granting **Attack Speed**, **Attack Damage**, **Movement Speed**, and **Nectar Heal** pickups<br>4. **12 Olympian Boons System**: Level-up overlay offering 3 random God Upgrades that dynamically adapt to Bow or Sword<br>5. **Wave Spawner & Boss AI**: Escalating waves (Skeletons -> Unshielded & Shielded Warriors -> Zeus's Champion -> Titan Cronus) | `ebitengineer` |
| **Phase 4: Polish, Verification & WASM Playtest** | **Minutes 25–30** | 1. Compile & verify with `godoctor` (`go vet`, native & `js/wasm` builds)<br>2. Launch **Game Dev Studio** (`sidecar://agy-game-studio/studio/playtest`) for live 60 FPS playtesting & DVR scrubbing | `godoctor` + `game-dev-studio` |

---

## 1. Executive Summary & Elevator Pitch
- **Elevator Pitch**: Step into the sandals of a mortal Greek hero ascending Mount Olympus in a fast-paced 2D top-down survival action game. Choose your sacred weapon at the start of each run—the **Golden Bow of Apollo** for long-range precision volleys or the **Bronze Xiphos Sword of Ares** for high-damage sweeping cleaves—stack divine *Archer.io*-style combat multipliers from the 12 Olympian Gods, and overthrow **Zeus's Champion** and **Titan Cronus**.
- **Core Inspiration**: *Archero* / *Archer.io* (stackable projectile & slash-wave modifiers: Multi-Attack, Diagonal Cones, Ricochet, Piercing) blended with *Vampire Survivors* (continuous movement + auto-combat flow) and *Hades* (start-of-run Weapon Aspect choice + Olympian God Boon selection UI).
- **Target Mood & Pacing**: Heroic mythological arcade action with distinct ranged vs. melee risk-reward playstyles, instant combat feedback, and triumphant ancient Greek modal music.

---

## 2. Core Gameplay Loop, Simulation & Mechanics
- **Start-of-Game Weapon Selection (`WeaponSelectScene`)**:
  Before entering the Olympus Courtyard, the player chooses their Hero's armament:
  1. **Golden Bow of Apollo (Ranged Gameplay)**:
     - **Range**: Long (`520 px`) | **Base Damage**: `28` | **Cooldown**: `0.42s` | **Hero Bonus**: `+10% Move Speed`
     - **Playstyle**: Kite monsters from afar and fill the arena with bouncing, piercing arrow volleys. Must flank or pierce **Shielded Hoplite Warriors** (`50%` frontal arrow block).
  2. **Bronze Sword of Ares (Melee Gameplay)**:
     - **Range**: Close-Mid Cleave Arc (`135 px` radius, `110°` wide sweep + short crescent slash wave) | **Base Damage**: `48` | **Cooldown**: `0.38s` | **Hero Bonus**: `+35 Max HP` & `15% Damage Reduction`
     - **Playstyle**: Aggressive close-quarters brawler that cleaves multiple Skeletons at once and **shatters/bypasses Shielded Hoplite frontal blocks** (`100%` damage through shields!).
- **Primary 30-Second Gameplay Loop**:
  1. **Move & Engage**: Weave through incoming enemies and boss telegraphs across the Mount Olympus marble arena while your hero **automatically aims and attacks** the nearest target with their chosen Bow or Sword.
  2. **Slay, Break Chests & Collect Power-Ups**: Defeated enemies and **Golden Olympus Chests** scattered around the map drop **Ambrosia XP Shards**, **Nectar Hearts**, and **Divine Stat Power-Ups** (**Attack Speed**, **Attack Damage**, **Movement Speed**).
  3. **Ascend with Olympian Boons**: Every level-up pauses the action and presents **3 of the 12 Gods of Olympus Boons**—mutating either your Arrow volleys or your Sword slash waves.
- **In-Game Power-Ups & Map Chests System**:
  - **Olympian Treasure Chests (`chest_olympus`)**:
    - Spawn periodically around the arena courtyard (`35 HP` or open on hero contact/attack).
    - Guarantee `1–2` stat power-up orbs plus a shower of Ambrosia XP shards.
  - **NPC Drop Table**:
    - **Skeleton Spartoi**: `100%` XP Shard + `12%` chance to drop a random Power-Up or Nectar Heart.
    - **Greek Warriors (Unshielded & Shielded)**: `100%` Large XP Shard + `25%` chance to drop a Stat Power-Up (**Damage**, **Attack Speed**, or **Movement Speed**).
    - **Zeus's Champion (Mid-Boss)**: Drops a **Golden Chest Burst** granting all 3 Stat Power-Ups + Full Heal.
  - **Collectible Power-Up Types**:
    1. **Attack Speed Orb (*Apollo's Sun Quiver*) — Gold/Yellow**: Grants `+15% Attack Speed` permanently for the run + a `6s` **Rapid-Fire / Rapid-Slash Overdrive** (`+40%` bonus speed aura).
    2. **Attack Damage Orb (*Ares's War Ember*) — Crimson/Orange**: Grants `+15% Base Damage` permanently for the run + a `6s` **Bloodfire Empower** (`+35%` bonus damage aura).
    3. **Movement Speed Orb (*Hermes's Winged Plume*) — Cyan/Azure**: Grants `+12% Move Speed` permanently for the run + a `6s` **Zephyr Haste** (`+35%` movement speed trail).
    4. **Nectar Chalice (*Hebe's Nectar*) — Emerald/Pink**: Restores `+30 HP` immediately.
- **Spatial & Simulation Rules**:
  - **Movement & Interaction Topology**: 360-degree / 8-way continuous movement at `220 px/sec` base speed inside a `1280x720` Olympus arena featuring 4 marble pillars for tactical kiting.
  - **Player Abilities & Dual-Weapon *Archer.io* Upgrade Patterns**:
    - **Artemis (Goddess of the Hunt)**: *Twin Strike* (**Bow**: +1 parallel front arrow | **Sword**: Double-slash combo per swing).
    - **Apollo (God of Archery & Sun)**: *Sunburst Cone* (**Bow**: +2 diagonal arrows at ±20° | **Sword**: +2 diagonal golden slash-waves at ±30°).
    - **Hermes (Messenger God)**: *Winged Sandals* (+25% attack speed & +15% move speed).
    - **Hephaestus (God of the Forge)**: *Adamant Edge* (**Bow**: Arrows pierce enemies + +25% dmg | **Sword**: +40% cleave radius + launches a long-range piercing blade beam).
    - **Poseidon (God of the Sea)**: *Tidal Force* (**Bow**: Arrows ricochet off walls/foes | **Sword**: Sweeps unleash bouncing tidal shockwaves + heavy knockback).
    - **Zeus (King of Olympus)**: *Thunderbolt Strike* (Hits trigger chain lightning arcs across nearby enemies).
    - **Athena (Goddess of Wisdom)**: *Aegis Barrier* (Orbiting golden shields block enemy projectiles and damage mobs on contact).
    - **Ares (God of War)**: *Bloodlust Crit* (+25% Critical Hit chance for 2.5x damage).
    - **Hades (Lord of the Underworld)**: *Soul Siphon* (Slaying enemies restores +3 HP).
    - **Demeter (Goddess of Harvest)**: *Golden Ambrosia* (+35% Max HP & instant heal).
    - **Dionysus (God of Wine)**: *Intoxicating Miasma* (Attacks apply stacking poison DoT).
    - **Hera (Queen of Heaven)**: *Sovereign Rearguard* (**Bow**: Rear arrow | **Sword**: 360° Whirlwind back-slash + +20% XP gain).
- **Enemy Roster & Systemic Challenges**:
  1. **Starting Mob — Skeleton Spartoi (`skeleton`)**:
     - Fast, brittle undead soldiers summoned from dragon teeth (`40 HP`, `115 px/s`). Swarm the player in packs during Waves 1–3.
  2. **Regular Mob — Greek Warriors (`warrior_unshielded` & `warrior_shielded`)**:
     - **Unshielded Warrior**: Aggressive bronze-sword infantry (`85 HP`, `100 px/s`) that lunges when within `140 px`.
     - **Shielded Hoplite Warrior**: Carries a heavy bronze Aspis shield (`150 HP`, `75 px/s`, **blocks 50% frontal projectile damage** and glows gold when deflecting arrows—encouraging piercing/ricochet boons or flanking angles).
  3. **Mid-Boss (Wave 2 Climax) — Zeus's Champion (`boss_zeus_champion`)**:
     - Towering golden-armored demigod (`1,200 HP`). Alternates between **Thunder Spear Dash** (telegraphed red rush line) and **8-Way Lightning Nova** (radial electric orbs).
  4. **Final Boss (Wave 3 Climax) — Titan Cronus (`boss_cronus`)**:
     - Colossal Lord of Time (`2,500 HP`). Executes three signature phases:
       - *Adamant Sickle Cleave*: Fires sweeping crescent energy waves toward the hero.
       - *Chronos Time-Rift*: Spawns a slow-field vortex on the arena floor while summoning Skeleton & Warrior reinforcements.
       - *Titan's Wrath (Enrage < 40% HP)*: Radial spiral barrage of golden time-shards.
- **Win & Loss Conditions**:
  - **Win Condition**: Survive all 3 Olympian Trial Waves (Wave 1: Skeleton Swarm & Warriors -> Wave 2: Phalanx Assault + **Zeus's Champion** -> Wave 3: Titan Legion + **Titan Cronus**) and defeat **Titan Cronus** with `HP > 0`.
  - **Loss Condition**: Hero HP reaches `0`. Immediate retry available with full run statistics (Damage Dealt, Boons Collected, Enemies Slain, Time Survived).

---

## 3. Controls & Input Mapping Scheme
- **Primary Input Devices**: Keyboard, Mouse/Pointer, Gamepad, Touch (Virtual Joystick)
- **Default Control Scheme**:

| Logical Action | Keyboard / Mouse | Gamepad Button | Touch Gesture |
| :--- | :--- | :--- | :--- |
| `ActionNavigate / Move` | `WASD` or `Arrow Keys` (or Hold Left Click to steer) | Left Stick / D-Pad | Virtual Joystick Drag |
| `ActionPrimary / Confirm` | **Auto-Fire** in combat / `1`, `2`, `3` or `Left Click` on Boon cards | Button South (`A` / `Cross`) | Tap Boon Card |
| `ActionSecondary / Dash` | `Space` or `Shift` (Short Dodge Roll with 3s cooldown) | Button East (`B` / `Circle`) | Double-Tap / Dash Button |
| `ActionPause` | `Escape` / `P` | Start / Options | Pause Icon Button |

---

## 4. Visual Art Direction & Sprite Pipeline (`sprite-design` + `nano-banana`)
- **Aesthetic Direction & Color Palette**: Mythic Greek 2D HD arcade art—sunlit white/gold Olympus marble, Aegean azure energy accents, Corinthian bronze armor, and crimson Spartan cloaks.
- **Sprite Resolution Tier**:
  - `[x]` **Native 2D HD (`256x256` per frame, scaled cleanly on `1280x720` canvas)** for Characters & Bosses.
  - `[x]` **Crisp `.sprite` Procedural Assets (`32x32`–`64x64`)** for high-density projectiles, XP gems, orbiting shields, and arena pillars.
- **Selected Visual Pipeline Mode**: **Path C: Hybrid Pipeline**
  - **AI Character Sheets & Arena Background (`nano-banana` + `sprite-design` normalize)**:
    1. `assets/images/olympus_bg.png`: `1280x720` top-down Mount Olympus marble courtyard with golden inlays, clouds at the perimeter, and ancient Greek mosaic emblem in the center.
    2. `assets/sprites/hero_greek.*`: Greek Hero (`idle`, `run`, `bow_attack`, `sword_attack`, `hit`, `death`).
    3. `assets/sprites/mob_skeleton.*`: Skeleton Spartoi (`idle`, `run`, `attack`, `hit`, `death`).
    4. `assets/sprites/mob_warrior.*`: Greek Hoplite Warrior (`idle`, `run`, `attack`, `shield_guard`, `hit`, `death`).
    5. `assets/sprites/boss_zeus_champion.*`: Zeus's Champion (`idle`, `move`, `attack`, `cast`, `death`).
    6. `assets/sprites/boss_cronus.*`: Titan Cronus (`idle`, `move`, `cleave`, `summon`, `death`).
  - **Pure-Code `.sprite` DSL Assets (`sprite-design`)**:
    - Golden Arrow, Bronze Sword Slash Arc / Wave, Lightning Bolt, Cronus Crescent Wave, **Olympian Treasure Chest (Closed & Open)**, **Power-Up Orbs (Attack Speed, Damage, Movement Speed)**, Ambrosia XP Gem, Nectar Heart, Athena Orbiting Shield, Olympus Marble Pillar.

---

## 5. Audio & Soundscape Pipeline (`procedural-sfx` + `procedural-music`)
- **Selected Audio Pipeline Mode**: **Pure-Code Procedural (`.sfx` + `.music`)** for instant compilation and zero-latency playback.
- **Sound Effects (`procedural-sfx` `.sfx` DSL)**:
  - **Gameplay / Combat**: `bow_shoot.sfx`, `sword_slash.sfx`, `arrow_hit.sfx`, `shield_clang.sfx`, `dash_whoosh.sfx`, `thunder_cast.sfx`, `cronus_cleave.sfx`, `hero_hurt.sfx`
  - **Pickups, Chests & Progression**: `chest_open.sfx`, `powerup_pickup.sfx`, `xp_ambrosia.sfx`, `boon_levelup.sfx`, `heal_nectar.sfx`, `victory_olympus.sfx`
  - **UI & Menus**: `ui_select.sfx`, `ui_confirm.sfx`
- **Background Music (`procedural-music` `.music` DSL)**:
  - **Track**: `assets/audio/olympus_12_gods.music`
  - **Tempo, Time Signature & Key**: `136 BPM`, `4/4` driving heroic rhythm, **D Phrygian Dominant / Dorian** modal Greek scale.
  - **Synthesis Engine & Palette**: Pure-Go SEGA YM2612 4-Op FM + SN76489 PSG (compatible with both Native & WASM out of the box), featuring:
    - *Apollo's Golden Lyre* (rapid plucked arpeggios)
    - *Ares's War Timpani & Sub-Bass* (anchored in `C1–C2` `< 260 Hz`)
    - *Athena's Bronze Fanfare* (heroic brass theme representing the 12 Olympians)

---

## 6. Scene Flow & UI / HUD Layout
- **Scene State Machine**: `TitleScene` -> `WeaponSelectScene` (*Bow of Apollo* vs. *Sword of Ares*) -> `GameplayScene` -> `BoonSelectOverlay` (Pause-on-Level-Up with 3 God Cards) / `PauseOverlay` -> `VictoryScene` / `GameOverScene`
- **HUD & UI Regions (`1280x720`)**:
  - **Top-Center**: **Ambrosia XP Bar** (Golden progress bar + current Hero Level & active Wave banner, plus Boss HP Bar when Zeus's Champion or Titan Cronus spawns).
  - **Top-Left**: **Hero HP Bar**, Dodge Dash cooldown indicator, and active **Olympian God Boon Icons** row.
  - **Top-Right**: **Wave Timer**, Enemies Remaining counter, and Kill Score.
  - **World-Space**: Floating damage numbers (white normal, gold critical, blue "BLOCKED" on Shielded Warriors) and overhead health bars.

---

## 7. Ebitengine Architecture & Studio Workflow (`ebitengineer` + `game-dev-studio`)
- **Core Runtime Subsystems (`ebitengineer`)**:
  - **Zero-Allocation Game Loop**: Pre-allocated pools for arrows, enemy projectiles, floating combat text, and hit particles; 4-byte aligned (`len(pcm) &^ 3`) 44.1 kHz stereo PCM audio buffers.
  - **Auto-Aim Targeting**: Spatial distance + line-of-sight scan picks the nearest living enemy every shot cooldown, firing `N` parallel and angled arrows based on the player's stacked Olympian Boons.
- **Interactive Workbench Verification (`game-dev-studio`)**:
  - **Sprite Editor (`sidecar://agy-game-studio/studio/sprite`)**: Inspect Hero, Skeleton, Warrior, Zeus's Champion, and Titan Cronus sheets.
  - **Sound Editor (`sidecar://agy-game-studio/studio/sound`)**: Audition `olympus_12_gods.music` and combat `.sfx`.
  - **Playtest Tab (`sidecar://agy-game-studio/studio/playtest`)**: Live 60 FPS WebAssembly playtesting and 10-second DVR replay scrubbing.
