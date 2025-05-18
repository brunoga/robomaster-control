package systems

import (
	"fmt"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster-control/components"
	"github.com/brunoga/robomaster-control/entities"
	"github.com/brunoga/robomaster/module/controller"
)

var (
	zeroStickPosition = controller.StickPosition{}
)

type Controller struct {
	controllerEntityMap map[uint64]*entities.Controller

	previousChassisStickPosition controller.StickPosition
	previousGimbalStickPosition  controller.StickPosition

	previousMoveTime time.Time
}

func (c *Controller) New(w *ecs.World) {
	c.controllerEntityMap = make(map[uint64]*entities.Controller)
}

func (c *Controller) Add(basicEntity *ecs.BasicEntity,
	controllerComponent *components.Controller, client *robomaster.Client) {
	_, ok := c.controllerEntityMap[basicEntity.ID()]
	if ok {
		return
	}

	c.controllerEntityMap[basicEntity.ID()] = &entities.Controller{
		BasicEntity: basicEntity,
		Controller:  controllerComponent,
		Client:      client,
	}
}

func (c *Controller) Remove(basicEntity ecs.BasicEntity) {
	delete(c.controllerEntityMap, basicEntity.ID())
}

func (c *Controller) Update(dt float32) {
	if btn := engo.Input.Button("exit"); btn.JustPressed() {
		engo.Exit()
	}

	if btn := engo.Input.Button("StartStop"); btn.JustPressed() {
		fmt.Println("Start/Stop")
		for _, controllerEntity := range c.controllerEntityMap {
			client := controllerEntity.Client
			if client.Connection().Connected() {
				err := client.Stop()
				if err != nil {
					panic(fmt.Sprintln("Error stopping client:", err))
				}
			} else {
				err := client.Start()
				if err != nil {
					panic(fmt.Sprintln("Error starting client:", err))
				}
			}
		}
		return
	}

	currentChassisStickPosition := controller.StickPosition{
		X: float64(engo.Input.Axis("Left/Right").Value()),
		Y: float64(engo.Input.Axis("Forward/Backward").Value()),
	}

	currentGimbalStickPosition := controller.StickPosition{
		X: float64(clampValueTo(engo.Input.Axis("MouseXAxis").Value(), 100) / 100),
		Y: float64(clampValueTo(engo.Input.Axis("MouseYAxis").Value(), 100) / 100),
	}

	currentMoveTime := time.Now()

	// Check if our move status changed.
	if currentChassisStickPosition == c.previousChassisStickPosition &&
		currentGimbalStickPosition == c.previousGimbalStickPosition {
		// Apparently not. Check if we are completelly stationary.
		if c.previousChassisStickPosition == zeroStickPosition &&
			c.previousGimbalStickPosition == zeroStickPosition {
			// We are completelly stationary. Nothing to do.
			return
		} else {
			// We are not completelly stationary. Maybe we should force a move.
			forceMove := time.Since(c.previousMoveTime) > time.Millisecond*900
			if !forceMove {
				// Nope. Nothing to do.
				return
			}
		}
	}

	// Update previous values to the current ones.
	c.previousChassisStickPosition = currentChassisStickPosition
	c.previousGimbalStickPosition = currentGimbalStickPosition
	c.previousMoveTime = currentMoveTime

	for _, controllerEntity := range c.controllerEntityMap {
		cec := controllerEntity.Controller

		cec.Controller.Move(&currentChassisStickPosition, &currentGimbalStickPosition,
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
