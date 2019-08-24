package main

import (
	"time"

	"github.com/zaklaus/rurik/src/system"

	rl "github.com/zaklaus/raylib-go/raylib"
)

const (
	pdaLayoutWidth  float32 = 960 * 0.55
	pdaLayoutHeight float32 = 540 * 0.65
	pdaLayoutX      float32 = 960*0.45 - (pdaLayoutWidth / 2)
	pdaLayoutY      float32 = 540*0.35 - (pdaLayoutHeight / 2)

	pdaScreenWidth  float32 = 463
	pdaScreenHeight float32 = 289
	pdaScreenX      float32 = 32
	pdaScreenY      float32 = 31
)

type pdaSystem struct {
	frameTexture *rl.Texture2D

	currentTimeAndDate time.Time

	installedApps []pdaApp
	activeApp     *pdaApp
}

type pdaApp interface {
	on()
	off()
	update()
	render()
}

type pdaAppBase struct {
	icon  [2]int32
	title string
	state int32
}

func makePDA() pdaSystem {
	return pdaSystem{
		frameTexture:       system.GetTexture("gfx/pda.png"),
		currentTimeAndDate: time.Now(),
		installedApps:      []pdaApp{},
		activeApp:          nil,
	}
}

func drawPDA(g *gameMode) {
	p := g.pda

	// draw a frame
	rl.DrawTexturePro(
		*p.frameTexture,
		rl.NewRectangle(0, 0, pdaLayoutWidth, pdaLayoutHeight),
		rl.NewRectangle(pdaLayoutX, pdaLayoutY, pdaLayoutWidth, pdaLayoutHeight),
		rl.Vector2{},
		0,
		rl.White,
	)
}

func updatePDA(g *gameMode) {

}
