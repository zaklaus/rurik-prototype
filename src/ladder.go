package main

import (
	"github.com/solarlune/resolv/resolv"
	"github.com/zaklaus/rurik/src/core"
)

type ladder struct{}

// NewLadder ladder
func NewLadder(o *core.Object) {
	o.IsCollidable = true
	o.CollisionType = "trigger"
	o.Size = []int32{int32(o.Meta.Width), int32(o.Meta.Height)}
	o.DebugVisible = false

	o.GetAABB = core.GetSolidAABB

	o.HandleCollisionEnter = func(res *resolv.Collision, o, other *core.Object) {
		other.Movement.Y = 0

		switch v := other.UserData.(type) {
		case *player:
			v.ctrl.IsOnLadder = true
		}
	}

	o.HandleCollisionLeave = func(res *resolv.Collision, o, other *core.Object) {
		switch v := other.UserData.(type) {
		case *player:
			v.ctrl.IsOnLadder = false
		}
	}
}
