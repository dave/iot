package grid

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type App any

type Widget interface {
	Layout(gtx layout.Context) layout.Dimensions
}

type CellInterface[A App] interface {
	init(A)
	Layout(gtx layout.Context) layout.Dimensions
}

type Grid[A App] struct {
	Theme         *material.Theme
	Rows, Columns int
	Cells         [][]CellInterface[A]
}

func (g *Grid[A]) Init(a A) {
	var columns int
	for _, row := range g.Cells {
		if len(row) > columns {
			columns = len(row)
		}
		for _, cell := range row {
			cell.init(a)
		}
	}
	g.Rows = len(g.Cells)
	g.Columns = columns
}

func (g *Grid[A]) Layout(gtx layout.Context) layout.Dimensions {
	gridWidth := gtx.Constraints.Max.X
	gridHeight := gtx.Constraints.Max.Y
	cellWidth := gridWidth / g.Columns
	cellHeight := gridHeight / g.Rows

	for rowIndex := 0; rowIndex < g.Rows; rowIndex++ {
		for colIndex := 0; colIndex < g.Columns; colIndex++ {
			xOffset := colIndex * cellWidth
			yOffset := rowIndex * cellHeight

			trans := op.Offset(f32.Pt(float32(xOffset), float32(yOffset))).Push(gtx.Ops)
			gtx := gtx
			gtx.Constraints = layout.Exact(image.Pt(cellWidth, cellHeight))
			if len(g.Cells[rowIndex]) > colIndex {
				g.Cells[rowIndex][colIndex].Layout(gtx)
			}
			trans.Pop()
		}
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
