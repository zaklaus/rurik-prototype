package main

import (
	"encoding/gob"

	"github.com/zaklaus/rurik/src/core"

	"github.com/solarlune/resolv/resolv"
	rl "github.com/zaklaus/raylib-go/raylib"
	ry "github.com/zaklaus/raylib-go/raymath"
	"github.com/zaklaus/rurik/src/system"
)

const ()

type player struct {
	ctrl *characterController
}

func (p *player) Serialize(enc *gob.Encoder) {
	enc.Encode(p)
}

func (p *player) Deserialize(dec *gob.Decoder) {
	dec.Decode(p)
}

// NewPlayer player
func NewPlayer(p *core.Object) {
	aseData := system.GetAnimData("gfx/player.json")
	p.Ase = &aseData
	p.Texture = system.GetTexture("gfx/player.png")
	p.Size = []int32{p.Ase.FrameWidth, p.Ase.FrameHeight}
	p.Update = updatePlayer
	p.Draw = drawPlayer
	p.GetAABB = core.GetSpriteAABB
	p.HandleCollision = handlePlayerCollision
	p.Facing = rl.NewVector2(1, 0)
	p.IsCollidable = true
	p.Finish = finishPlayer
	p.InsideArea = func(o, area *core.Object) bool {
		return system.IsKeyPressed("use")
	}
	p.CanTrigger = true

	core.LocalPlayer = p

	core.PlayAnim(p, "StandE")
	plr := &player{
		ctrl: &characterController{
			Object: p,
		},
	}

	p.UserData = plr
}

func finishPlayer(p *core.Object) {}

func updatePlayer(p *core.Object, dt float32) {
	p.Ase.Update(dt)

	ctrl := p.UserData.(*player).ctrl

	if core.CanSave == 0 || core.BitsHas(core.CanSave, core.IsInChallenge) {
		factor := system.GetAxis("horizontal")
		ctrl.move(factor)

		if system.IsKeyDown("jump") || system.GetAxis("vertical") < 0 {
			ctrl.jump()
		}

		if system.GetAxis("vertical") > 0 {
			ctrl.down()
		}

		if core.DebugMode && rl.IsKeyDown(rl.KeyU) {
			pushWaterParticle(p.GetWorld(), p.Position)
		}
	} else {
		return
	}

	var tag string

	if ry.Vector2Length(p.Movement) > 0 {
		p.Facing.X = core.SignFloat(p.Movement.X)

		tag = "Walk"
	} else {
		tag = "Stand"
	}

	if p.Facing.X > 0 {
		tag += "E"
	} else if p.Facing.X < 0 {
		tag += "W"
	}

	core.PlayAnim(p, tag)
	ctrl.update()
}

func drawPlayer(p *core.Object) {
	source := core.GetSpriteRectangle(p)
	dest := core.GetSpriteOrigin(p)

	if core.DebugMode && p.DebugVisible {
		c := core.GetSpriteAABB(p)
		rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
		core.DrawTextCentered(p.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
	}

	rl.DrawTexturePro(*p.Texture, source, dest, rl.Vector2{}, 0, rl.White)
}

func handlePlayerCollision(res *resolv.Collision, p, other *core.Object) {}
