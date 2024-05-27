package scenes

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster-control/components"
	"github.com/brunoga/robomaster-control/entities"
	"github.com/brunoga/robomaster-control/systems"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Robomaster struct {
	Client *robomaster.Client
}

func (r *Robomaster) Preload() {
	// Do nothing.
}

func (r *Robomaster) Setup(u engo.Updater) {
	engo.Input.RegisterAxis("Left/Right",
		engo.AxisKeyPair{Min: engo.KeyA, Max: engo.KeyD})
	engo.Input.RegisterAxis("Forward/Backward",
		engo.AxisKeyPair{Min: engo.KeyW, Max: engo.KeyS})

	engo.Input.RegisterAxis("MouseXAxis",
		engo.NewAxisMouse(engo.AxisMouseHori))
	engo.Input.RegisterAxis("MouseYAxis",
		engo.NewAxisMouse(engo.AxisMouseVert))

	engo.Input.RegisterButton("exit", engo.KeyEscape)

	controller := &components.Controller{
		Controller: r.Client.Controller(),
	}

	controllerBasicEntity := ecs.NewBasic()

	controllerEntity := entities.Controller{
		BasicEntity: &controllerBasicEntity,
		Controller:  controller,
	}

	gunComponent := &components.Gun{
		Gun: r.Client.Gun(),
	}

	gunBasicEntity := ecs.NewBasic()

	gunEntity := entities.Gun{
		BasicEntity: &gunBasicEntity,
		Gun:         gunComponent,
	}

	// Disable cursor.
	if engo.CurrentBackEnd == engo.BackEndGLFW ||
		engo.CurrentBackEnd == engo.BackEndVulkan {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode,
			glfw.CursorDisabled)
	} else {
		panic("Backend does not seem to support mouse capture.")
	}

	w, _ := u.(*ecs.World)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&systems.Video{
		Camera: r.Client.Camera(),
	})
	w.AddSystem(&systems.Controller{})
	w.AddSystem(&systems.Gun{})
	w.AddSystem(&systems.Information{
		Robot:      r.Client.Robot(),
		Connection: r.Client.Connection(),
	})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *systems.Controller:
			sys.Add(controllerEntity.BasicEntity,
				controller)
		case *systems.Gun:
			sys.Add(gunEntity.BasicEntity, gunComponent)
		}
	}
}

func (r *Robomaster) Type() string {
	return "Robomaster"
}

func (r *Robomaster) Exit() {
	fmt.Println("Exiting...")
}
