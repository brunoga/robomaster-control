package systems

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/robot"

	"golang.org/x/image/font/gofont/gomonobold"
)

// Information is a system that displays battery percentage.
type Information struct {
	entity struct {
		*ecs.BasicEntity
		*common.RenderComponent
		*common.SpaceComponent
	}
	elapsed    float32
	Font       *common.Font // Font used to display the FPS to the screen, defaults to gomonobold
	Robot      *robot.Robot
	Connection *connection.Connection
}

// New is called when FPSSystem is added to the world
func (f *Information) New(w *ecs.World) {
	if f.Font == nil {
		if err := engo.Files.LoadReaderData("gomonobold_fps.ttf", bytes.NewReader(gomonobold.TTF)); err != nil {
			panic("unable to load gomonobold.ttf for the fps system! Error was: " + err.Error())
		}

		f.Font = &common.Font{
			URL:  "gomonobold_fps.ttf",
			FG:   color.White,
			BG:   color.Black,
			Size: 32,
		}

		if err := f.Font.CreatePreloaded(); err != nil {
			panic("unable to create gomonobold.ttf for the fps system! Error was: " + err.Error())
		}
	}

	txt := common.Text{
		Font: f.Font,
		Text: f.DisplayString(),
	}
	b := ecs.NewBasic()
	f.entity.BasicEntity = &b
	f.entity.RenderComponent = &common.RenderComponent{
		Drawable: txt,
	}
	f.entity.RenderComponent.SetShader(common.HUDShader)
	f.entity.RenderComponent.SetZIndex(1000)
	f.entity.SpaceComponent = &common.SpaceComponent{}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(f.entity.BasicEntity, f.entity.RenderComponent, f.entity.SpaceComponent)
		}
	}
}

// Add doesn't do anything since New creates the only entity used
func (*Information) Add() {}

// Remove doesn't do anything since New creates the only entity used
func (*Information) Remove(b ecs.BasicEntity) {}

// Update changes the dipslayed text and prints to the terminal every second
// to report the FPS
func (f *Information) Update(dt float32) {
	f.elapsed += dt
	text := f.DisplayString()
	if f.elapsed >= 1 {
		f.entity.Drawable = common.Text{
			Font: f.Font,
			Text: text,
		}
		f.elapsed--
	}
}

// DisplayString returns the display string in the format Battery: 60%
func (f *Information) DisplayString() string {
	batteryString := "Battery: N/A"
	if f.Robot != nil {
		batteryString = fmt.Sprintf("Battery: %d%%", f.Robot.BatteryPowerPercent())
	}

	signalString := "Signal: N/A"
	if f.Connection != nil {
		signalString = fmt.Sprintf("Signal: %d/4", f.Connection.SignalQualityBars())
	}

	return fmt.Sprintf("%s %s", batteryString, signalString)
}
