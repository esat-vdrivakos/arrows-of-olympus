package main

import (
	"log"

	"github.com/drivakos/arrows-of-olympus/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(game.VirtualWidth, game.VirtualHeight)
	ebiten.SetWindowTitle("Arrows of Olympus — 12 Gods Survival Arena")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)

	g, err := game.New()
	if err != nil {
		log.Fatalf("arrows-of-olympus boot failed: %v", err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("arrows-of-olympus runtime error: %v", err)
	}
}
