package main

import (
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type characterController struct {
	Object     *core.Object
	IsGrounded bool
	IsFalling  bool
	IsInWater  bool
	IsOnLadder bool
}

func (c *characterController) move(factor float32) {
	c.Object.Movement.X = factor * movementSpeed * system.FrameTime
}

func (c *characterController) jump() {
	if c.IsGrounded {
		c.Object.Movement.Y = -jumpForce * system.FrameTime
	}

	if c.IsInWater {
		c.Object.Movement.Y = -upwardWaterForce * system.FrameTime
	}

	if c.IsOnLadder {
		c.Object.Movement.Y = -ladderClimbSpeed * system.FrameTime
	}
}

func (c *characterController) down() {
	if c.IsOnLadder {
		c.Object.Movement.Y = ladderClimbSpeed * system.FrameTime
	}

	if c.IsInWater {
		c.Object.Movement.Y = (upwardWaterForce + buoyancy/2) * system.FrameTime
	}
}

func (c *characterController) update() {
	// Handle free fall
	{
		down, _ := core.CheckForCollisionEx("*", c.Object, 0, 4)
		c.IsGrounded = down.Colliding()
		if !c.IsGrounded && !c.IsInWater && !c.IsOnLadder {
			g := gravity

			c.Object.Movement.Y += g * system.FrameTime

			if c.Object.Movement.Y > maxFallSpeed {
				c.Object.Movement.Y = maxFallSpeed
			}
		}
	}

	x := core.RoundFloat(c.Object.Movement.X)
	y := core.RoundFloat(c.Object.Movement.Y)

	// Handle collision
	{
		// Handle slope movement
		if res, _ := core.CheckForCollisionEx("slope", c.Object, core.RoundFloatToInt32(x), core.RoundFloatToInt32(y)+4); res.Colliding() && !res.Teleporting {
			y = core.RoundInt32ToFloat(res.ResolveY)
			c.Object.Movement.Y = 0
		}

		// Handle solid+trigger collisions
		if res, _ := core.CheckForCollisionEx("solid+trigger", c.Object, core.RoundFloatToInt32(x), 0); res.Colliding() && !res.Teleporting {
			x = core.RoundInt32ToFloat(res.ResolveX)
			c.Object.Movement.X = 0
		}

		if res, _ := core.CheckForCollisionEx("solid+trigger", c.Object, 0, core.RoundFloatToInt32(y)); res.Colliding() && !res.Teleporting {
			y = core.RoundInt32ToFloat(res.ResolveY)
			c.Object.Movement.Y = 0
		}

		// Handle ceiling solid+trigger
		if res, _ := core.CheckForCollisionEx("solid+trigger", c.Object, 0, -4); res.Colliding() && !res.Teleporting {
			y = core.RoundInt32ToFloat(res.ResolveY)
			c.Object.Movement.Y = 0
		}

		// Apply motion
		c.Object.Position.X += x
		c.Object.Position.Y += y
	}

	c.IsFalling = c.Object.Movement.Y > 0 && !c.IsInWater && !c.IsOnLadder

	if c.IsOnLadder {
		c.Object.Movement.Y = 0
	}
}
