package renderer

import (
	"context"
	"fmt"
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/c-kuroki/alien_invasion/pkg/model"
)

// check that interface is implemented
var _ Adapter = (*SVGRenderer)(nil)

const (
	cityColor  = "fill:#283f93"
	alienColor = "fill:#19e822"
	fightColor = "fill:#e81922"
)

// SVGRenderer render an invasion map on a SVG image
type SVGRenderer struct {
	width      int
	height     int
	citySize   int
	cityWidth  int
	connSize   int
	alienWidth int
}

func NewSVGRenderer() *SVGRenderer {
	return &SVGRenderer{
		width:      800,
		height:     800,
		citySize:   120,
		cityWidth:  40,
		connSize:   10,
		alienWidth: 12,
	}
}

func (r *SVGRenderer) Render(ctx context.Context, cities []*model.City, aliens map[int]map[int]*model.Alien, w io.Writer) error {
	canvas := svg.New(w)
	canvas.Start(r.width, r.height)
	for _, city := range cities {
		x := r.citySize * city.X
		y := r.citySize * city.Y
		// render city
		// connections
		if city.North != "" {
			canvas.Rect(x+r.citySize/2-r.connSize/2, y, r.connSize, r.citySize/2, cityColor)
		}
		if city.South != "" {
			canvas.Rect(x+r.citySize/2-r.connSize/2, y+r.citySize/2, r.connSize, r.citySize/2, cityColor)
		}
		if city.East != "" {
			canvas.Rect(x+r.citySize/2, y+r.citySize/2-r.connSize/2, r.citySize/2, r.connSize, cityColor)
		}
		if city.West != "" {
			canvas.Rect(x, y+r.citySize/2-r.connSize/2, r.citySize/2, r.connSize, cityColor)
		}
		// name
		canvas.Circle(x+r.citySize/2, y+r.citySize/2, r.cityWidth, cityColor)
		canvas.Text(x+r.citySize/2, y+r.citySize/2, city.Name, "text-anchor:middle;font-size:16px;font-family:helvetica;fill:white")
		// aliens
		numAliens := len(aliens[city.ID])
		if numAliens > 0 {
			if numAliens == 1 {
				canvas.Circle(x+r.citySize/2, y+r.citySize/2+r.citySize/8, r.alienWidth, alienColor)
			} else {
				canvas.Circle(x+r.citySize/2, y+r.citySize/2+r.citySize/8, r.alienWidth, fightColor)
			}
			canvas.Text(x+r.citySize/2, y+r.citySize/2+r.citySize/8+r.alienWidth/3, fmt.Sprintf("%d", numAliens), "text-anchor:middle;font-size:10px;font-family:helvetica;fill:black")
		}
	}
	canvas.End()
	return nil
}
