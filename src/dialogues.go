package main

import (
	"fmt"
	"log"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
	"gopkg.in/yaml.v2"
)

const (
	// MouseDoublePress default duration of mouse double press
	MouseDoublePress = 500
)

var dialogues = make(map[string]Dialogue)

type dialogueData struct {
	texts                *Dialogue
	currentText          *Dialogue
	selectedChoice       int
	extraTick            bool
	mouseDoublePressTime int32
}

var dialogue dialogueData

// Dialogue defines connversation flow
type Dialogue struct {
	Name       string    `yaml:"name"`
	AvatarFile string    `yaml:"avatar"`
	Text       string    `yaml:"text"`
	Choices    []*Choice `yaml:"choices"`
	Event      string    `yaml:"event"`
	EventArgs  string    `yaml:"eventArgs"`
	SkipPrompt bool      `yaml:"skipPrompt"`
	Next       *Dialogue `yaml:"next"`
	avatar     *rl.Texture2D
}

// Choice is a selection from dialogue branches
type Choice struct {
	Text string    `yaml:"text"`
	Next *Dialogue `yaml:"next"`
}

// InitText initializes the dialogue's text
func InitText(t *Dialogue) {
	if t.AvatarFile != "" {
		t.avatar = system.GetTexture("gfx/" + t.AvatarFile)
	}

	if t.Next != nil {
		InitText(t.Next)
	}

	if t.Choices != nil {
		for _, ch := range t.Choices {
			if ch.Next != nil {
				InitText(ch.Next)
			}
		}
	}
}

// GetDialogue retrieves dialogue.texts for a dialogue
func GetDialogue(name string) *Dialogue {
	dia, ok := dialogues[name]

	if ok {
		return &dia
	}

	data := system.GetFile(fmt.Sprintf("texts/%s", name), false)
	err := yaml.Unmarshal(data, &dia)

	if err != nil {
		log.Printf("Dialogue '%s' is broken!\n", name)
		return &Dialogue{}
	}

	dialogues[name] = dia
	return &dia
}

// InitDialogue initializes a dialogue
func InitDialogue(name string) {
	if dialogue.extraTick {
		return
	}

	log.Printf("Initializing dialogue '%s' ...\n", name)
	dialogue.texts = GetDialogue(name)
	dialogue.currentText = dialogue.texts
	dialogue.extraTick = false
	InitText(dialogue.currentText)
}

func updateDialogue() {
	if core.CurrentMap == nil {
		dialogue = dialogueData{}
		return
	}

	if dialogue.texts == nil {
		if dialogue.extraTick {
			dialogue.extraTick = system.IsKeyDown("use")
		}
		return
	}

	if !dialogue.extraTick {
		dialogue.extraTick = system.IsKeyReleased("use")
		return
	}

	core.CanSave = core.BitsSet(core.CanSave, core.IsInDialogue)

	if dialogue.mouseDoublePressTime > 0 {
		dialogue.mouseDoublePressTime -= int32(1000 * (system.FrameTime * float32(core.TimeScale)))
	} else if dialogue.mouseDoublePressTime < 0 {
		dialogue.mouseDoublePressTime = 0
	}

	if len(dialogue.currentText.Choices) > 0 {
		if system.IsKeyPressed("up") {
			dialogue.selectedChoice--

			if dialogue.selectedChoice < 0 {
				dialogue.selectedChoice = len(dialogue.currentText.Choices) - 1
			}
		}

		if system.IsKeyPressed("down") {
			dialogue.selectedChoice++

			if dialogue.selectedChoice >= len(dialogue.currentText.Choices) {
				dialogue.selectedChoice = 0
			}
		}
	}

	if system.IsKeyPressed("use") || (rl.IsMouseButtonReleased(rl.MouseLeftButton) && dialogue.mouseDoublePressTime > 0) {
		if dialogue.mouseDoublePressTime > 0 {
			dialogue.mouseDoublePressTime = 0
		}

		evnt := dialogue.currentText.Event
		evntArglist := dialogue.currentText.EventArgs
		evntArgs := core.CompileEventArgs(evntArglist)

		if len(dialogue.currentText.Choices) > 0 {
			dialogue.currentText = dialogue.currentText.Choices[dialogue.selectedChoice].Next
		} else {
			dialogue.currentText = dialogue.currentText.Next
		}

		if dialogue.currentText != nil && dialogue.currentText.SkipPrompt {
			evnt = dialogue.currentText.Event
			evntArglist = dialogue.currentText.EventArgs
			evntArgs = []string{evntArglist}

			dialogue.currentText = nil
		}

		if dialogue.currentText == nil {
			dialogue.texts = nil
			dialogue.extraTick = true
			core.CanSave = core.BitsClear(core.CanSave, core.IsInDialogue)
		}

		if evnt != "" {
			core.FireEvent(evnt, evntArgs)
		}
	}
}

func drawDialogue() {
	if dialogue.texts == nil {
		return
	}

	var height int32 = 120
	width := system.WindowWidth
	start := system.ScreenHeight - height

	rl.DrawRectangle(0, start, width, height, rl.NewColor(46, 46, 84, 255))
	rl.DrawRectangle(5, start+5, 32, 32, rl.NewColor(53, 64, 59, 255))
	rl.DrawRectangleLines(4, start+4, 34, 34, rl.NewColor(55, 148, 110, 255))

	ot := dialogue.currentText

	// Pos X: 5, Y: 5
	// Scale W: 34, 35
	if ot.AvatarFile != "" {
		rl.DrawTexturePro(
			*ot.avatar,
			rl.NewRectangle(0, 0, float32(ot.avatar.Width), float32(ot.avatar.Height)),
			rl.NewRectangle(5, float32(start)+5, 32, 32),
			rl.Vector2{},
			0,
			rl.White,
		)
	}

	rl.DrawText(
		ot.Name,
		45,
		start+16,
		10,
		rl.Orange,
	)

	rl.DrawText(
		ot.Text,
		5,
		start+45,
		10,
		rl.White,
	)

	// choices
	chsX := system.ScreenWidth - 220
	chsY := start + 16

	if len(ot.Choices) > 0 {
		for idx, ch := range ot.Choices {
			ypos := chsY + int32(idx)*15 - 2
			if idx == dialogue.selectedChoice {
				rl.DrawRectangle(chsX, ypos, 200, 15, rl.DarkPurple)
			}

			rl.DrawText(
				fmt.Sprintf("%d. %s", idx+1, ch.Text),
				chsX+5,
				chsY+int32(idx)*15,
				10,
				rl.White,
			)

			if core.IsMouseInRectangle(chsX, ypos, 200, 15) {
				if rl.IsMouseButtonDown(rl.MouseLeftButton) {
					rl.DrawRectangleLines(chsX, ypos, 200, 15, rl.Pink)
				} else {
					rl.DrawRectangleLines(chsX, ypos, 200, 15, rl.Purple)
				}

				if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
					dialogue.selectedChoice = idx

					dialogue.mouseDoublePressTime = MouseDoublePress
				}
			}
		}
	} else {
		rl.DrawRectangle(chsX, chsY-2, 200, 15, rl.DarkPurple)
		rl.DrawText(
			"Press E to continue...",
			chsX+5,
			chsY,
			10,
			rl.White,
		)
	}
}
