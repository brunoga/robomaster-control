package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster-control/components"
)

type Chassis struct {
	*ecs.BasicEntity
	*components.Chassis
}
