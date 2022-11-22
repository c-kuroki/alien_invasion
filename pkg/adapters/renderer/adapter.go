package renderer

import (
	"context"
	"io"

	"github.com/c-kuroki/alien_invasion/pkg/model"
)

//go:generate mockery --name Renderer

// Renderer Adapter interface to render world map
type Adapter interface {
	Render(ctx context.Context, cities []*model.City, aliens map[int]map[int]*model.Alien, w io.Writer) error
}
