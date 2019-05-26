package main

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
)

type physicsProps struct {
	IsGrounded        bool
	IsFalling         bool
	IsInWater         bool
	IsOnLadder        bool
	IsGettingOnLadder bool
}
