package main

import (
	"encoding/gob"

	"github.com/zaklaus/rurik/src/core"

	rl "github.com/zaklaus/raylib-go/raylib"
	ry "github.com/zaklaus/raylib-go/raymath"
	"github.com/zaklaus/resolv/resolv"
	"github.com/zaklaus/rurik/src/system"
)

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
	p.GetAABB = getFixedSpriteAABB
	p.HandleCollision = handlePlayerCollision
	p.Facing = rl.NewVector2(1, 0)
	p.CollisionType = collisionPawn
	p.IsCollidable = true
	p.Finish = finishPlayer
	p.InsideArea = func(o, area *core.Object) bool {
		return system.IsKeyPressed("use")
	}
	p.CanTrigger = true
	p.DebugVisible = true

	core.LocalPlayer = p

	core.PlayAnim(p, "StandE")
	plr := &player{
		ctrl: &characterController{
			Object: p,
		},
	}

	p.UserData = plr

	p.HandleCollisionEnter = func(res *resolv.Collision, o, other *core.Object) {
		switch other.Class {
		case "water":
			o.UserData.(*player).ctrl.IsInWater = true
		case "ladder":
			o.UserData.(*player).ctrl.IsOnLadder = true
		}
	}

	p.HandleCollisionLeave = func(res *resolv.Collision, o, other *core.Object) {
		switch other.Class {
		case "water":
			o.UserData.(*player).ctrl.IsInWater = false
		case "ladder":
			o.UserData.(*player).ctrl.IsOnLadder = false
		}
	}
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

	tag := "Stand"

	if ry.Vector2Length(p.Movement) > 0 {
		p.Facing.X = core.SignFloat(p.Movement.X)
	}

	core.PlayAnim(p, tag)
	ctrl.update()
}

func drawPlayer(p *core.Object) {
	source := core.GetSpriteRectangle(p)
	dest := core.GetSpriteOrigin(p)

	if p.Facing.X == -1 {
		source.Width *= -1
	}

	rl.DrawTexturePro(*p.Texture, source, dest, rl.Vector2{}, 0, rl.White)
}

func handlePlayerCollision(res *resolv.Collision, p, other *core.Object) {}
