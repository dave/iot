package main

import (
	"github.com/dave/iot/ui/button"
	"github.com/dave/iot/ui/grid"
)

var root = &grid.Grid[*App]{
	Cells: [][]grid.CellInterface[*App]{
		{
			&grid.Cell[*App, *button.Button]{
				Contents: &button.Button{Id: "dimmer1"},
				Init:     func(b *button.Button, a *App) { a.handle(b) },
			},
			&grid.Cell[*App, *button.Button]{
				Contents: &button.Button{Id: "dimmer2"},
				Init:     func(b *button.Button, a *App) { a.handle(b) },
			},
		},
		{
			&grid.Cell[*App, *button.Button]{
				Contents: &button.Button{Id: "dimmer3"},
				Init:     func(b *button.Button, a *App) { a.handle(b) },
			},
		},
	},
}
