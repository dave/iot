package button

import (
	"fmt"
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/dave/iot/ui/basic"
)

type Button struct {
	Theme   *material.Theme
	Id      string
	Value   float32
	Changed func(value float32)

	drag bool
}

func (b *Button) Layout(gtx layout.Context) layout.Dimensions {

	cellWidth := float32(gtx.Constraints.Max.X)
	cellHeight := float32(gtx.Constraints.Max.Y)
	buttonWidth := cellWidth * 0.3
	buttonHeight := cellHeight * 0.8
	offsetX := (cellWidth - buttonWidth) * 0.5
	offsetY := (cellHeight - buttonHeight) * 0.5

	for _, ev := range gtx.Events(b) {
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Type {
		case pointer.Drag, pointer.Press:
			previous := b.Value
			b.drag = true
			positionOfset := e.Position.Y - offsetY
			b.Value = 1.0 - (positionOfset / buttonHeight)
			switch {
			case b.Value > 1.0:
				b.Value = 1.0
			case b.Value < 0.0:
				b.Value = 0.0
			}
			if int(b.Value*100) != int(previous*100) {
				// only call changed if count has actually changed by more than 1%
				b.Changed(b.Value)
			}
			b.drag = false
		case pointer.Release:
			b.drag = false
		}
	}

	inner := image.Rect(
		int(offsetX),
		int(offsetY),
		int(offsetX+buttonWidth),
		int(offsetY+buttonHeight),
	)
	area := clip.Rect(inner).Push(gtx.Ops)
	pointer.InputOp{
		Tag:   b,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}.Add(gtx.Ops)
	area.Pop()

	area1 := clip.Rect(inner).Push(gtx.Ops)
	paint.ColorOp{Color: basic.White(0.7)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	area1.Pop()

	innerShaded := image.Rect(
		int(offsetX),
		int(offsetY+buttonHeight*(1-b.Value)),
		int(offsetX+buttonWidth),
		int(offsetY+buttonHeight),
	)
	area2 := clip.Rect(innerShaded).Push(gtx.Ops)
	paint.ColorOp{Color: basic.Black(0.6)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	area2.Pop()

	{
		offset := op.Offset(f32.Pt(offsetX, offsetY)).Push(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(int(buttonWidth), int(buttonHeight)))
		title := material.Body1(b.Theme, fmt.Sprintf("Value: %d%%", int(b.Value*100.0)))
		title.Color = basic.White(1)
		title.Alignment = text.Middle
		title.Layout(gtx)
		offset.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
