package systems

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster-control/components"
	"github.com/brunoga/robomaster-control/entities"
	"github.com/brunoga/robomaster/module/chassis/controller"
)

type Chassis struct {
	controllerEntityMap map[uint64]*entities.Chassis

	previousLeftRight       float32
	previousForwardBackward float32
	previousMouseXDelta     float32
	previousMouseYDelta     float32

	lastMove time.Time
}

func (c *Chassis) New(w *ecs.World) {
	c.controllerEntityMap = make(map[uint64]*entities.Chassis)
}

func (c *Chassis) Add(basicEntity *ecs.BasicEntity,
	controllerComponent *components.Chassis) {
	_, ok := c.controllerEntityMap[basicEntity.ID()]
	if ok {
		return
	}

	c.controllerEntityMap[basicEntity.ID()] = &entities.Chassis{
		BasicEntity: basicEntity,
		Chassis:     controllerComponent,
	}
}

func (c *Chassis) Remove(basicEntity ecs.BasicEntity) {
	delete(c.controllerEntityMap, basicEntity.ID())
}

func (c *Chassis) Update(dt float32) {
	if btn := engo.Input.Button("exit"); btn.JustPressed() {
		engo.Exit()
	}

	currentLeftRight := engo.Input.Axis("Left/Right").Value()
	currentForwardBackward := engo.Input.Axis("Forward/Backward").Value()

	currentMouseXDelta := clampValueTo(engo.Input.Axis("MouseXAxis").Value(), 100)
	currentMouseYDelta := clampValueTo(engo.Input.Axis("MouseYAxis").Value(), 100)

	// Force a move if we are close to 1 second since the last move. A single move
	// gets the robot moving for 1 second.
	forceMove := time.Since(c.lastMove) > time.Millisecond*900

	// Check if any movenet happened, if not, just return. We do this because
	// it is wasteful to send requests to the robot 60 times per second.
	if !forceMove && currentLeftRight == c.previousLeftRight &&
		currentForwardBackward == c.previousForwardBackward &&
		currentMouseXDelta == c.previousMouseXDelta &&
		currentMouseYDelta == c.previousMouseYDelta {
		return
	}

	// Update previous values to the current ones.
	c.previousLeftRight = currentLeftRight
	c.previousForwardBackward = currentForwardBackward
	c.previousMouseXDelta = currentMouseXDelta
	c.previousMouseYDelta = currentMouseYDelta

	c.lastMove = time.Now()

	for _, controllerEntity := range c.controllerEntityMap {
		cec := controllerEntity.Chassis

		chassisStickPosition := &controller.StickPosition{
			X: float64(currentLeftRight),
			Y: float64(currentForwardBackward),
		}

		gimbalStickPosition := &controller.StickPosition{
			X: float64(currentMouseXDelta) / float64(100),
			Y: float64(currentMouseYDelta) / float64(100),
		}

		cec.Chassis.Move(chassisStickPosition, gimbalStickPosition,
			controller.ModeFPV)
	}
}

func clampValueTo(value, clamp float32) float32 {
	if value > clamp {
		return clamp
	} else if value < -clamp {
		return -clamp
	}
	return value
}
