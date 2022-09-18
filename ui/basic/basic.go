package basic

import "image/color"

var (
	Background = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	Red        = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	Green      = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	Blue       = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
)

func Black(f float32) color.NRGBA {
	return White(1 - f)
}

func White(f float32) color.NRGBA {
	return color.NRGBA{R: uint8(f * 255), G: uint8(f * 255), B: uint8(f * 255), A: 0xFF}
}

//// ColorBox creates a widget with the specified dimensions and color.
//func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
//	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
//	paint.ColorOp{Color: color}.Add(gtx.Ops)
//	paint.PaintOp{}.Add(gtx.Ops)
//	return layout.Dimensions{Size: size}
//}

//	func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
//		ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
//		return layout.Center.Layout(gtx, material.H3(th, text).Layout)
//	}
