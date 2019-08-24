package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/raylib-go/raymath"
)

func questInitMathCommands(q *questManager) {
	q.registerCommand("vec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("vec", qs, qt, len(args), 1)
		}

		vecName := args[0]
		qs.setVector(vecName, rl.Vector2{})

		qs.printf(qt, "vector '%s' was declared!", vecName)
		return true
	})

	q.registerCommand("setvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("setvec", qs, qt, len(args), 3)
		}

		vecName := args[0]

		xI, _ := qs.getNumberOrVariable(args[1])
		x := float64to32(xI)
		yI, _ := qs.getNumberOrVariable(args[2])
		y := float64to32(yI)

		qs.setVector(vecName, rl.NewVector2(x, y))

		qs.printf(qt, "vector '%s' was set to [%f, %f]!", vecName, xI, yI)
		return true
	})

	q.registerCommand("copyvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 2 {
			return questCommandErrorArgCount("copyvec", qs, qt, len(args), 2)
		}

		vecName := args[0]
		rhsVecName := args[1]

		rhs, ok := qs.getVector(rhsVecName)

		if !ok {
			return questCommandErrorThing("copyvec", "vector", qs, qt, rhsVecName)
		}

		qs.setVector(vecName, rhs)
		return true
	})

	q.registerCommand("getvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("getvec", qs, qt, len(args), 3)
		}

		vecName := args[0]
		xName := args[1]
		yName := args[2]

		vec, ok := qs.getVector(vecName)

		if !ok {
			return questCommandErrorThing("getvec", "vector", qs, qt, vecName)
		}

		if xName != "0" {
			qs.setVariable(xName, float64(vec.X))
		}

		if yName != "0" {
			qs.setVariable(yName, float64(vec.Y))
		}

		return true
	})

	q.registerCommand("addvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("addvec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhs, rhsFound := qs.getVector(rhsVecName)

		if !lhsFound {
			return questCommandErrorThing("addvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("addvec", "vector", qs, qt, rhsVecName)
		}

		qs.setVector(destVecName, rl.NewVector2(lhs.X+rhs.X, lhs.Y+rhs.Y))
		return true
	})

	q.registerCommand("addivec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("addivec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhsI, rhsFound := qs.getNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return questCommandErrorThing("addivec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("addivec", "number", qs, qt, args[2])
		}

		qs.setVector(destVecName, rl.NewVector2(lhs.X+rhs, lhs.Y+rhs))
		return true
	})

	q.registerCommand("subvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("subvec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhs, rhsFound := qs.getVector(rhsVecName)

		if !lhsFound {
			return questCommandErrorThing("subvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("subvec", "vector", qs, qt, rhsVecName)
		}

		qs.setVector(destVecName, rl.NewVector2(lhs.X-rhs.X, lhs.Y-rhs.Y))
		return true
	})

	q.registerCommand("subivec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("subivec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhsI, rhsFound := qs.getNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return questCommandErrorThing("subivec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("subivec", "number", qs, qt, args[2])
		}

		qs.setVector(destVecName, rl.NewVector2(lhs.X-rhs, lhs.Y-rhs))
		return true
	})

	q.registerCommand("divivec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("divivec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhsI, rhsFound := qs.getNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return questCommandErrorThing("divivec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("divivec", "number", qs, qt, args[2])
		}

		if rhs == 0 {
			return questCommandErrorDivideByZero("divivec", qs, qt)
		}

		qs.setVector(destVecName, rl.NewVector2(lhs.X/rhs, lhs.Y/rhs))
		return true
	})

	q.registerCommand("mulvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("mulvec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhsI, rhsFound := qs.getNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return questCommandErrorThing("mulvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("mulvec", "number", qs, qt, args[2])
		}

		qs.setVector(destVecName, rl.NewVector2(lhs.X*rhs, lhs.Y*rhs))
		return true
	})

	q.registerCommand("dotvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("dotvec", qs, qt, len(args), 3)
		}

		destName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhs, rhsFound := qs.getVector(rhsVecName)

		if !lhsFound {
			return questCommandErrorThing("dotvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("dotvec", "number", qs, qt, args[2])
		}

		res := raymath.Vector2DotProduct(lhs, rhs)
		qs.setVariable(destName, float64(res))
		return true
	})

	q.registerCommand("crossvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 3 {
			return questCommandErrorArgCount("crossvec", qs, qt, len(args), 3)
		}

		destName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.getVector(lhsVecName)
		rhs, rhsFound := qs.getVector(rhsVecName)

		if !lhsFound {
			return questCommandErrorThing("crossvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return questCommandErrorThing("crossvec", "number", qs, qt, args[2])
		}

		res := raymath.Vector2CrossProduct(lhs, rhs)
		qs.setVariable(destName, float64(res))
		return true
	})

	q.registerCommand("normvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 2 {
			return questCommandErrorArgCount("normvec", qs, qt, len(args), 2)
		}

		destName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.getVector(lhsVecName)

		if !lhsFound {
			return questCommandErrorThing("normvec", "vector", qs, qt, lhsVecName)
		}

		raymath.Vector2Normalize(&lhs)
		qs.setVector(destName, lhs)
		return true
	})

	q.registerCommand("flipvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 2 {
			return questCommandErrorArgCount("flipvec", qs, qt, len(args), 2)
		}

		destName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.getVector(lhsVecName)

		if !lhsFound {
			return questCommandErrorThing("flipvec", "vector", qs, qt, lhsVecName)
		}

		qs.setVector(destName, rl.NewVector2(lhs.Y, -lhs.X))
		return true
	})

	q.registerCommand("lenvec", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 2 {
			return questCommandErrorArgCount("lenvec", qs, qt, len(args), 2)
		}

		destName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.getVector(lhsVecName)

		if !lhsFound {
			return questCommandErrorThing("lenvec", "vector", qs, qt, lhsVecName)
		}

		qs.setVariable(destName, float64(raymath.Vector2Length(lhs)))
		return true
	})
}
