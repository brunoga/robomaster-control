package main

import (
	"flag"
	"log/slog"

	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster-control/scenes"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
)

var (
	fullscreen = flag.Bool("fullscreen", false, "Run in fullscreen mode.")
)

func main() {
	flag.Parse()

	l := logger.New(slog.LevelError)

	c, err := robomaster.NewWithModules(l, 0, module.TypeAllButGamePad)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = c.Stop()
		if err != nil {
			panic(err)
		}
	}()

	err = c.Robot().EnableFunction(robot.FunctionTypeMovementControl, true)
	if err != nil {
		panic(err)
	}

	// Cache and restore on exit the current speed level.
	speedLevel, err := c.Robot().ChassisSpeedLevel()
	if err != nil {
		panic(err)
	}
	defer func(speedLevel robot.ChassisSpeedLevelType) {
		err := c.Robot().SetChassisSpeedLevel(speedLevel)
		if err != nil {
			panic(err)
		}
	}(speedLevel)

	// Set the speed level to slow.
	err = c.Robot().SetChassisSpeedLevel(robot.ChassisSpeedLevelFast)
	if err != nil {
		panic(err)
	}

	opts := engo.RunOptions{
		Title:         "Robomaster",
		Width:         1280,
		Height:        720,
		VSync:         true,
		ScaleOnResize: true,
		FPSLimit:      60,
		Fullscreen:    *fullscreen,
	}

	engo.Run(opts, &scenes.Robomaster{
		Client: c,
	})
}
