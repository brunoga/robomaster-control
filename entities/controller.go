package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster-control/components"
)

type Controller struct {
	*ecs.BasicEntity
	*components.Controller
}
