package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// State captures a single frame's normalized input across keyboard, mouse, and gamepad.
type State struct {
	MoveX          float64
	MoveY          float64
	DashJustPressed bool
	PauseJustPressed bool
	ConfirmJustPressed bool
	MuteJustPressed bool
	Key1JustPressed bool
	Key2JustPressed bool
	Key3JustPressed bool
	LeftJustPressed bool
	RightJustPressed bool
	MouseX         float64
	MouseY         float64
	MouseJustClicked bool
	MouseHeld      bool
	AnyActivity    bool
	gamepadIDs     []ebiten.GamepadID
	touchIDs       []ebiten.TouchID
	justTouchedIDs []ebiten.TouchID
}

// Poll samples all active input devices for the current tick without heap allocations.
func (s *State) Poll(heroX, heroY float64, allowMouseSteer bool) {
	s.MoveX = 0
	s.MoveY = 0
	s.DashJustPressed = false
	s.PauseJustPressed = false
	s.ConfirmJustPressed = false
	s.MuteJustPressed = false
	s.Key1JustPressed = false
	s.Key2JustPressed = false
	s.Key3JustPressed = false
	s.LeftJustPressed = false
	s.RightJustPressed = false
	s.AnyActivity = false

	// Fullscreen toggle (F11 or Alt+Enter)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) ||
		(ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Keyboard movement
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.MoveX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.MoveX += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.MoveY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		s.MoveY += 1
	}

	s.LeftJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	s.RightJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight)

	s.DashJustPressed = inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	s.PauseJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP)
	s.ConfirmJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
	s.MuteJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyM)

	s.Key1JustPressed = inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad1)
	s.Key2JustPressed = inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad2)
	s.Key3JustPressed = inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad3)

	// Mouse / Pointer / Touch
	mx, my := ebiten.CursorPosition()
	s.MouseX = float64(mx)
	s.MouseY = float64(my)
	s.MouseJustClicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	s.MouseHeld = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Mobile Touch & Tap controls (zero-allocation polling)
	s.justTouchedIDs = inpututil.AppendJustPressedTouchIDs(s.justTouchedIDs[:0])
	if len(s.justTouchedIDs) > 0 {
		tx, ty := ebiten.TouchPosition(s.justTouchedIDs[0])
		s.MouseX = float64(tx)
		s.MouseY = float64(ty)
		s.MouseJustClicked = true
		// Two-finger tap triggers Dodge Roll on mobile
		if len(s.justTouchedIDs) >= 2 {
			s.DashJustPressed = true
		}
	}
	s.touchIDs = ebiten.AppendTouchIDs(s.touchIDs[:0])
	if len(s.touchIDs) > 0 {
		tx, ty := ebiten.TouchPosition(s.touchIDs[0])
		s.MouseX = float64(tx)
		s.MouseY = float64(ty)
		s.MouseHeld = true
	}

	if allowMouseSteer && s.MouseHeld && s.MoveX == 0 && s.MoveY == 0 {
		dx := s.MouseX - heroX
		dy := s.MouseY - heroY
		if math.Hypot(dx, dy) > 18.0 {
			s.MoveX = dx
			s.MoveY = dy
		}
	}

	// Gamepad polling
	s.gamepadIDs = ebiten.AppendGamepadIDs(s.gamepadIDs[:0])
	for _, id := range s.gamepadIDs {
		gx := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		gy := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
		if math.Hypot(gx, gy) > 0.22 {
			s.MoveX += gx
			s.MoveY += gy
		}
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftLeft) {
			s.MoveX -= 1
		}
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftRight) {
			s.MoveX += 1
		}
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftTop) {
			s.MoveY -= 1
		}
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftBottom) {
			s.MoveY += 1
		}
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftLeft) {
			s.LeftJustPressed = true
		}
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftRight) {
			s.RightJustPressed = true
		}
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) {
			s.ConfirmJustPressed = true
		}
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightRight) {
			s.DashJustPressed = true
		}
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonCenterRight) {
			s.PauseJustPressed = true
		}
	}

	// Normalize diagonal movement vector
	mag := math.Hypot(s.MoveX, s.MoveY)
	if mag > 1e-5 {
		s.MoveX /= mag
		s.MoveY /= mag
		s.AnyActivity = true
	}

	if s.DashJustPressed || s.PauseJustPressed || s.ConfirmJustPressed || s.MuteJustPressed ||
		s.Key1JustPressed || s.Key2JustPressed || s.Key3JustPressed || s.MouseJustClicked ||
		s.LeftJustPressed || s.RightJustPressed {
		s.AnyActivity = true
	}
}
