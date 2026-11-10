// Command assetcheck compiles and validates every embedded sprite, SFX, and music asset.
// Usage: go run ./cmd/assetcheck [-out dir]   (-out writes preview PNGs and WAVs)
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/drivakos/arrows-of-olympus/assets"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procmusic"
	"github.com/drivakos/arrows-of-olympus/internal/audio/procsfx"
	"github.com/drivakos/arrows-of-olympus/internal/sprite"
)

func main() {
	out := flag.String("out", "", "optional directory for preview PNG/WAV exports")
	flag.Parse()
	if *out != "" {
		_ = os.MkdirAll(*out, 0o755)
	}
	failed := 0

	sprites, _ := fs.Glob(assets.FS, "sprites/*.sprite")
	for _, p := range sprites {
		raw, _ := assets.FS.ReadFile(p)
		cs, err := sprite.CompileSpriteDSL(string(raw))
		if err != nil {
			fmt.Printf("FAIL sprite %s: %v\n", p, err)
			failed++
			continue
		}
		rep := sprite.ValidateSprite(cs)
		tags := make([]string, 0, len(cs.Tags))
		for _, t := range cs.Tags {
			tags = append(tags, fmt.Sprintf("%s[%d-%d]", t.Name, t.FromFrame, t.ToFrame))
		}
		status := "ok  "
		if !rep.Passed {
			status = "WARN"
		}
		fmt.Printf("%s sprite %-36s %dx%d frames=%d tags=%s\n", status, p, cs.FrameW, cs.FrameH, len(cs.Frames), strings.Join(tags, ","))
		for _, is := range rep.Issues {
			fmt.Printf("       - %s\n", is)
		}
		if *out != "" {
			name := strings.TrimSuffix(filepath.Base(p), ".sprite")
			_ = sprite.ExportPreviewPNG(filepath.Join(*out, name+".png"), 2, cs)
		}
	}

	sfx, _ := fs.Glob(assets.FS, "audio/*.sfx")
	for _, p := range sfx {
		raw, _ := assets.FS.ReadFile(p)
		pcm, err := procsfx.CompileDSL(string(raw))
		if err != nil {
			fmt.Printf("FAIL sfx %s: %v\n", p, err)
			failed++
			continue
		}
		rep, _ := procsfx.ValidateProductionSFX(pcm)
		status := "ok  "
		if !rep.Passed {
			status = "WARN"
		}
		fmt.Printf("%s sfx    %-36s %.2fs peak=%.1fdBFS rms=%.1fdBFS\n", status, p, rep.DurationSec, rep.PeakDBFS, rep.RMSDBFS)
		for _, is := range rep.Issues {
			fmt.Printf("       - %s\n", is)
		}
		if *out != "" {
			_ = procsfx.ExportWAV(filepath.Join(*out, strings.TrimSuffix(filepath.Base(p), ".sfx")+".wav"), pcm)
		}
	}

	music, _ := fs.Glob(assets.FS, "audio/*.music")
	for _, p := range music {
		raw, _ := assets.FS.ReadFile(p)
		score, err := procmusic.ParseDSL(string(raw))
		if err != nil {
			fmt.Printf("FAIL music %s: parse: %v\n", p, err)
			failed++
			continue
		}
		if err := procmusic.ValidateProductionBGM(score, 4, 20); err != nil {
			fmt.Printf("WARN music %s: %v\n", p, err)
		}
		pcm, err := procmusic.CompileDSL(string(raw))
		if err != nil {
			fmt.Printf("FAIL music %s: compile: %v\n", p, err)
			failed++
			continue
		}
		peak, clipped := pcmPeak(pcm)
		status := "ok  "
		if clipped > 0 || peak < 0.05 {
			status = "WARN"
		}
		fmt.Printf("%s music  %-36s tracks=%d %.1fs peak=%.2f clipped=%d\n", status, p, len(score.Tracks), float64(len(pcm))/4/44100, peak, clipped)
		if *out != "" {
			_ = procmusic.ExportWAV(filepath.Join(*out, strings.TrimSuffix(filepath.Base(p), ".music")+".wav"), pcm)
		}
	}

	if failed > 0 {
		fmt.Printf("\n%d asset(s) failed\n", failed)
		os.Exit(1)
	}
}

// pcmPeak returns the linear peak and the number of full-scale (clipped) samples in 16-bit stereo PCM.
func pcmPeak(pcm []byte) (float64, int) {
	var peak float64
	clipped := 0
	for i := 0; i+1 < len(pcm); i += 2 {
		v := int16(uint16(pcm[i]) | uint16(pcm[i+1])<<8)
		if v >= 32766 || v <= -32767 {
			clipped++
		}
		a := float64(v) / 32768.0
		if a < 0 {
			a = -a
		}
		if a > peak {
			peak = a
		}
	}
	return peak, clipped
}
