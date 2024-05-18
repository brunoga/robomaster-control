package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster-control/components"
)

type Gun struct {
	*ecs.BasicEntity
	*components.Gun
}
