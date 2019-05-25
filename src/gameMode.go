package main

import (
	"encoding/gob"
	"math"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type gameMode struct {
	playState      int
	textWave       int32
	showHelpScreen bool
}

const (
	stateTitleScreen = iota
	statePlay
	statePaused
	stateLevelSelection
)

func (g *gameMode) Init() {
	initLevels()

	g.playState = stateLevelSelection
}

func (g *gameMode) Shutdown() {}

func (g *gameMode) Update() {
	switch g.playState {
	case statePaused:
		if rl.IsKeyPressed(rl.KeyEscape) {
			g.playState = statePlay
		}

		if system.IsKeyPressed("use") {
			core.FlushMaps()
			g.playState = stateLevelSelection
			levelSelection.selectedChoice = 0
			return
		}

	case stateTitleScreen:
		g.textWave = int32(math.Round(math.Sin(float64(rl.GetTime()) * 10)))

		if system.IsKeyPressed("use") {
			g.playState = stateLevelSelection
		}

		if rl.IsKeyPressed(rl.KeyEscape) {
			core.CloseGame()
			return
		}

	case stateLevelSelection:
		g.updateLevelSelection()

		if rl.IsKeyPressed(rl.KeyEscape) {
			g.playState = stateTitleScreen
		}

	case statePlay:
		core.UpdateMaps()
		updateDialogue()
		updateNotifications()

		/* particle systems */
		updateWaterParticles()

		if rl.IsKeyPressed(rl.KeyEscape) && core.CurrentMap.Name != "start" {
			g.playState = statePaused
		}

		if rl.IsKeyPressed(rl.KeyF1) {
			g.showHelpScreen = !g.showHelpScreen
		}
	}
}

func (g *gameMode) Serialize(enc *gob.Encoder) {
	data := demoGameSaveData{}

	enc.Encode(data)
}

func (g *gameMode) Deserialize(dec *gob.Decoder) {
	var saveData demoGameSaveData
	dec.Decode(&saveData)
}

type demoGameSaveData struct {
	ObjectCounter int
}

func (g *gameMode) Draw() {
	rl.BeginMode2D(core.RenderCamera)
	{
		core.DrawMap(false)
		drawWaterParticles()
	}
	rl.EndMode2D()
}

func (g *gameMode) DrawUI() {
	switch g.playState {
	case stateTitleScreen:
		core.DrawTextCentered("Darkorbia", system.ScreenWidth/2, system.ScreenHeight/2-20+g.textWave, 24, rl.RayWhite)
		core.DrawTextCentered("Press E/ENTER to continue", system.ScreenWidth/2, system.ScreenHeight/2+5+g.textWave, 14, rl.White)

	case statePaused:
		rl.DrawRectangle(0, 0, system.ScreenWidth, system.ScreenHeight, rl.Fade(rl.Black, 0.8))
		core.DrawTextCentered("Darkorbia", system.ScreenWidth/2, system.ScreenHeight/2-20+g.textWave, 24, rl.RayWhite)
		core.DrawTextCentered("Press ESC to unpause or E/ENTER to return to the menu", system.ScreenWidth/2, system.ScreenHeight/2+5+g.textWave, 14, rl.White)

	case stateLevelSelection:
		core.DrawTextCentered("Darkorbia", system.ScreenWidth/2, system.ScreenHeight/2-20+g.textWave, 24, rl.RayWhite)
		g.drawLevelSelection()

	case statePlay:
		core.DrawMapUI()
		drawDialogue()
		drawNotifications()
	}
}

func (g *gameMode) PostDraw() {

	switch g.playState {
	case stateTitleScreen:

	case statePaused:
		fallthrough

	case statePlay:
		// Generates and applies the lightmaps
		core.UpdateLightingSolution()
	}

}
