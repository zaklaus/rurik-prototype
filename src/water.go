package main

import (
	"github.com/solarlune/resolv/resolv"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type water struct{}

// NewWater water
func NewWater(o *core.Object) {
	o.IsCollidable = true
	o.CollisionType = "trigger"
	o.Size = []int32{int32(o.Meta.Width), int32(o.Meta.Height)}
	o.DebugVisible = false
	o.ContainedObjects = []core.TriggerContact{}

	o.Update = func(o *core.Object, dt float32) {
		for _, v := range o.ContainedObjects {
			other := v.Object

			other.Movement.Y = core.ScalarLerp(other.Movement.Y, buoyancy*system.FrameTime, 0.30)
		}
	}

	o.GetAABB = core.GetSolidAABB

	o.HandleCollisionEnter = func(res *resolv.Collision, o, other *core.Object) {
		switch v := other.UserData.(type) {
		case *player:
			v.ctrl.IsInWater = true
		}
	}

	o.HandleCollisionLeave = func(res *resolv.Collision, o, other *core.Object) {
		switch v := other.UserData.(type) {
		case *player:
			v.ctrl.IsInWater = false
		}
	}
}
