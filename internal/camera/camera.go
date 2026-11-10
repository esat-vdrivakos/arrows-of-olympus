package camera

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera manages virtual canvas coordinates (1280x720), trauma screen shake, and hitstop freeze frames.
type Camera struct {
	Width        float64
	Height       float64
	Trauma       float64
	ShakeOffsetX float64
	ShakeOffsetY float64
	HitstopTimer float64
	FlashTimer   float64
	rng          *rand.Rand
}

// New creates a new 1280x720 virtual camera.
func New(width, height float64) *Camera {
	return &Camera{
		Width:  width,
		Height: height,
		rng:    rand.New(rand.NewSource(42)),
	}
}

// AddTrauma adds screen-shake trauma in [0, 1].
func (c *Camera) AddTrauma(amount float64) {
	c.Trauma = math.Min(1.0, c.Trauma+amount)
}

// TriggerHitstop pauses gameplay motion briefly for impactful hits (seconds).
func (c *Camera) TriggerHitstop(duration float64) {
	if duration > c.HitstopTimer {
		c.HitstopTimer = duration
	}
}

// TriggerFlash triggers a subtle golden/lightning screen flash.
func (c *Camera) TriggerFlash(duration float64) {
	if duration > c.FlashTimer {
		c.FlashTimer = duration
	}
}

// Reset clears transient trauma, shake, and hitstop state on scene exit.
func (c *Camera) Reset() {
	c.Trauma = 0
	c.ShakeOffsetX = 0
	c.ShakeOffsetY = 0
	c.HitstopTimer = 0
	c.FlashTimer = 0
}

// Update decays camera trauma and returns true if simulation is currently paused by hitstop.
func (c *Camera) Update(dt float64) bool {
	if c.FlashTimer > 0 {
		c.FlashTimer = math.Max(0, c.FlashTimer-dt)
	}
	if c.Trauma > 0 {
		shake := c.Trauma * c.Trauma
		maxOffset := 12.0
		c.ShakeOffsetX = maxOffset * shake * (c.rng.Float64()*2.0 - 1.0)
		c.ShakeOffsetY = maxOffset * shake * (c.rng.Float64()*2.0 - 1.0)
		c.Trauma = math.Max(0, c.Trauma-1.9*dt)
	} else {
		c.ShakeOffsetX = 0
		c.ShakeOffsetY = 0
	}
	if c.HitstopTimer > 0 {
		c.HitstopTimer = math.Max(0, c.HitstopTimer-dt)
		return c.HitstopTimer > 0
	}
	return false
}

// Apply translates an ebiten.DrawImageOptions by the active screen-shake offset.
func (c *Camera) Apply(op *ebiten.DrawImageOptions) {
	if c.ShakeOffsetX != 0 || c.ShakeOffsetY != 0 {
		op.GeoM.Translate(math.Round(c.ShakeOffsetX), math.Round(c.ShakeOffsetY))
	}
}
