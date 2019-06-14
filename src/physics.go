package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
)

const (
	gravity                 float32 = 20
	buoyancy                float32 = 80
	upwardWaterForce        float32 = 80
	ladderClimbSpeed        float32 = 120
	maxFallSpeed            float32 = 840
	jumpForce               float32 = 290
	movementSpeed           float32 = 144
	movementSmoothingFactor float32 = 0.22
	movementFallSpeedFactor float32 = 0.25

	waterVertexCount      int32   = 12
	waterVertexWindFactor float32 = 1
)

type physicsProps struct {
	IsGrounded        bool
	IsFalling         bool
	IsInWater         bool
	IsOnLadder        bool
	IsGettingOnLadder bool
}

func calculateContactResponse(props *physicsProps, resolve int32) (float32, float32) {
	res := core.RoundInt32ToFloat(resolve)

	return res, 0
}

func getFixedSpriteAABB(o *core.Object) rl.RectangleInt32 {
	if o.Ase == nil {
		return rl.RectangleInt32{
			X:      int32(o.Position.X),
			Y:      int32(o.Position.Y - 32),
			Width:  32,
			Height: 32,
		}
	}

	return rl.RectangleInt32{
		X:      int32(o.Position.X) - int32(float32(o.Ase.FrameWidth/2)) + int32(float32(o.Ase.FrameWidth/4)),
		Y:      int32(o.Position.Y) - int32(float32(o.Ase.FrameHeight/2)),
		Width:  o.Ase.FrameWidth / 2,
		Height: o.Ase.FrameHeight,
	}
}
