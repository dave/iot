package grid

import "gioui.org/layout"

type Cell[A App, C Widget] struct {
	Contents C
	Init     func(A, C)
}

func (c *Cell[A, C]) init(a A) {
	c.Init(a, c.Contents)
}

func (c *Cell[A, C]) Layout(gtx layout.Context) layout.Dimensions {
	return c.Contents.Layout(gtx)
}
