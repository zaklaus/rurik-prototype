package main

import (
	"encoding/gob"

	"github.com/zaklaus/raylib-go/raymath"

	"github.com/solarlune/resolv/resolv"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
)

type ball struct {
	physicsProps
}

func (b *ball) Serialize(enc *gob.Encoder)   {}
func (b *ball) Deserialize(dec *gob.Decoder) {}

// NewBall test ball
func NewBall(o *core.Object) {
	o.IsCollidable = true
	o.Size = []int32{int32(o.Meta.Width), int32(o.Meta.Height)}
	o.GetAABB = core.GetSolidAABB
	o.CollisionType = core.CollisionRigid
	o.DebugVisible = true
	o.UserData = &ball{}

	o.Draw = func(o *core.Object) {
		rect := o.GetAABB(o)

		rl.DrawCircle(
			rect.X+rect.Width/2,
			rect.Y+rect.Height/2,
			float32(rect.Width),
			rl.Blue,
		)
	}

	o.Update = func(o *core.Object, dt float32) {
		v := o.UserData.(*ball)

		if v.IsInWater {
			o.Movement.X = core.ScalarLerp(o.Movement.X, 0, 0.189)
		} else {
			o.Movement.Y += gravity * dt
			o.Movement.X = core.ScalarLerp(o.Movement.X, 0, 0.12)
		}

		dx := core.RoundFloat(o.Movement.X)
		dy := core.RoundFloat(o.Movement.Y)

		if res, _ := core.CheckForCollisionEx([]uint32{core.CollisionSolid, core.CollisionTrigger, collisionPawn}, o, core.RoundFloatToInt32(dx), 0); res.Colliding() && !res.Teleporting {
			dx = float32(res.ResolveX)
			o.Movement.X = 0
		}

		core.CheckForCollisionEx([]uint32{core.CollisionTrigger}, o, 0, 4)

		if res, _ := core.CheckForCollisionEx([]uint32{core.CollisionSolid}, o, 0, core.RoundFloatToInt32(dy)); res.Colliding() && !res.Teleporting {
			diff := float32(res.ResolveY)
			dy = -diff
			o.Movement.Y = 0
		}

		o.Position.X += dx
		o.Position.Y += dy
	}

	o.HandleCollision = func(res *resolv.Collision, a, b *core.Object) {
		v := raymath.Vector2Subtract(b.Position, a.Position)
		raymath.Vector2Normalize(&v)

		a.Movement.X += -v.X * absFloat(b.Movement.X)
		a.Movement.Y += -v.Y * absFloat(b.Movement.Y)
	}
}
