package audio

import (
	"bytes"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const SampleRate = 44100

// Manager coordinates a single 44.1 kHz stereo PCM audio.Context, looping BGM,
// WASM gesture unlock seeking, and cooldown-rate-limited SFX playback.
type Manager struct {
	ctx          *audio.Context
	bgmPlayer    *audio.Player
	bgmLen       int64
	bgmStart     time.Time
	wasmUnlocked bool
	muted        bool
	bgmWanted    bool
	bgmVolume    float64
	sfxVolume    float64
	sfxCues      map[string][]byte
	lastPlayed   map[string]time.Time
}

// NewManager creates the shared audio manager and initializes the looping BGM player.
func NewManager(bgmPCM []byte, sfxCues map[string][]byte) (*Manager, error) {
	ctx := audio.CurrentContext()
	if ctx == nil {
		ctx = audio.NewContext(SampleRate)
	}
	m := &Manager{
		ctx:        ctx,
		bgmVolume:  0.55,
		sfxVolume:  0.80,
		sfxCues:    sfxCues,
		lastPlayed: make(map[string]time.Time, len(sfxCues)),
		bgmStart:   time.Now(),
	}
	if err := m.SetBGM(bgmPCM); err != nil {
		return nil, err
	}
	return m, nil
}

// SetBGM installs the looping BGM track (e.g. once background synthesis finishes).
// It must be called from the game goroutine. Playback starts if BGM was requested earlier.
func (m *Manager) SetBGM(bgmPCM []byte) error {
	alignedLen := int64(len(bgmPCM) &^ 3)
	if m == nil || alignedLen == 0 || m.bgmPlayer != nil {
		return nil
	}
	loop := audio.NewInfiniteLoop(bytes.NewReader(bgmPCM[:alignedLen]), alignedLen)
	p, err := m.ctx.NewPlayer(loop)
	if err != nil {
		return err
	}
	p.SetVolume(m.bgmVolume)
	m.bgmPlayer = p
	m.bgmLen = alignedLen
	if m.bgmWanted && !m.muted {
		m.bgmStart = time.Now()
		p.Play()
	}
	return nil
}

// StartBGM begins or resumes looping BGM playback.
func (m *Manager) StartBGM() {
	if m == nil {
		return
	}
	m.bgmWanted = true
	if m.bgmPlayer == nil || m.muted {
		return
	}
	if !m.bgmPlayer.IsPlaying() {
		m.bgmStart = time.Now()
		m.bgmPlayer.Play()
	}
}

// UnlockOnUserGesture synchronizes BGM playback to a 4-byte aligned offset on first input (for WASM).
func (m *Manager) UnlockOnUserGesture() {
	if m == nil || m.wasmUnlocked || m.bgmPlayer == nil || m.bgmLen <= 0 {
		return
	}
	m.wasmUnlocked = true
	elapsed := time.Since(m.bgmStart).Seconds()
	samples := int64(elapsed * SampleRate)
	offset := ((samples * 4) % m.bgmLen) &^ 3
	_ = m.bgmPlayer.Seek(time.Duration(offset) * time.Second / (SampleRate * 4))
	if !m.muted && !m.bgmPlayer.IsPlaying() {
		m.bgmPlayer.Play()
	}
}

// ToggleMute toggles global audio mute.
func (m *Manager) ToggleMute() bool {
	if m == nil {
		return false
	}
	m.muted = !m.muted
	if m.bgmPlayer != nil {
		if m.muted {
			m.bgmPlayer.Pause()
		} else {
			m.bgmPlayer.SetVolume(m.bgmVolume)
			m.bgmPlayer.Play()
		}
	}
	return m.muted
}

// IsMuted reports whether audio is currently muted.
func (m *Manager) IsMuted() bool {
	if m == nil {
		return true
	}
	return m.muted
}

// PlaySFX plays a pre-compiled 4-byte aligned PCM sound cue with a cooldown gate.
func (m *Manager) PlaySFX(name string, cooldown time.Duration) {
	if m == nil || m.ctx == nil || m.muted {
		return
	}
	pcm, ok := m.sfxCues[name]
	if !ok {
		return
	}
	aligned := len(pcm) &^ 3
	if aligned == 0 {
		return
	}
	if cooldown <= 0 {
		cooldown = 45 * time.Millisecond
	}
	now := time.Now()
	if prev, exists := m.lastPlayed[name]; exists && now.Sub(prev) < cooldown {
		return
	}
	m.lastPlayed[name] = now

	p := m.ctx.NewPlayerFromBytes(pcm[:aligned])
	p.SetVolume(m.sfxVolume)
	p.Play()
}
