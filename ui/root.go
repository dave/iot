package main

import (
	"github.com/dave/iot/ui/grid"
	"github.com/dave/iot/ui/tool"
)

var root = &grid.Grid[*App]{
	Cells: [][]grid.CellInterface[*App]{
		{
			&grid.Cell[*App, *tool.Slider]{
				Contents: &tool.Slider{Name: "Dimmer 1", Id: "dimmer1"},
				Init:     func(a *App, b *tool.Slider) { a.handle(b) },
			},
			&grid.Cell[*App, *tool.Slider]{
				Contents: &tool.Slider{Name: "Dimmer 2", Id: "dimmer2"},
				Init:     func(a *App, b *tool.Slider) { a.handle(b) },
			},
		},
		{
			&grid.Cell[*App, *tool.Slider]{
				Contents: &tool.Slider{Name: "Dimmer 3", Id: "dimmer3"},
				Init:     func(a *App, b *tool.Slider) { a.handle(b) },
			},
		},
	},
}
