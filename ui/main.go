package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const MQTT_DELAY = time.Millisecond * 250

type ValueMessage struct {
	Value  float32 `json:"value"`
	Source string  `json:"source"`
}

var grid = &Grid{
	Cells: [][]GridCellInterface{
		{
			&GridCell[*Button]{
				Contents: &Button{},
				Init: func(b *Button, ui *Ui) {
					b.th = ui.theme
					b.changed = func(value float32) {
						ui.Publish("dimmer1/value", ValueMessage{Value: value, Source: "button"})
					}
					ui.mqtt.Subscribe("dimmer1/value", 0, func(client mqtt.Client, msg mqtt.Message) {
						var v ValueMessage
						if err := json.Unmarshal(msg.Payload(), &v); err != nil {
							fmt.Println(err)
							return
						}
						if v.Source == ui.clinetId {
							// don't react to our own messages
							return
						}
						b.value = v.Value
						b.count = int(b.value * 1024)
						ui.window.Invalidate()
					})
				},
			},
			&GridCell[*Button]{
				Contents: &Button{},
				Init: func(b *Button, ui *Ui) {
					b.th = ui.theme
					b.changed = func(value float32) {
						//MqttClient.Publish("dimmer2/value", 0, false, fmt.Sprint(value))
					}
				},
			},
		},
		{
			&GridCell[*Button]{
				Contents: &Button{},
				Init: func(b *Button, ui *Ui) {
					b.th = ui.theme
					b.changed = func(value float32) {
						//MqttClient.Publish("dimmer3/value", 0, false, fmt.Sprint(value))
					}
				},
			},
		},
	},
}

var defaultPublishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type Ui struct {
	theme     *material.Theme
	window    *app.Window
	clinetId  string
	mqtt      mqtt.Client
	pubs      map[string]any
	pubsMutex *sync.Mutex
}

func (ui *Ui) Publish(topic string, message any) {
	ui.pubsMutex.Lock()
	defer ui.pubsMutex.Unlock()

	_, found := ui.pubs[topic]
	if found {
		// we published recently to this topic, don't publish right now, but store the message
		ui.pubs[topic] = message
		return
	}

	// publish the message immediately
	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	ui.mqtt.Publish(topic, 0, true, b)
	ui.pubs[topic] = nil
	go func() {
		// start a goroutine to wait for MQTT_DELAY before publishing the next message
		time.Sleep(MQTT_DELAY)
		ui.pubsMutex.Lock()
		defer ui.pubsMutex.Unlock()
		if message, found := ui.pubs[topic]; found {
			if message == nil {
				delete(ui.pubs, topic)
				return
			}
			b, err := json.Marshal(message)
			if err != nil {
				delete(ui.pubs, topic)
				fmt.Println(err)
				return
			}
			ui.mqtt.Publish(topic, 0, true, b)
			delete(ui.pubs, topic)
		}
	}()
}

func main() {

	ui := &Ui{
		pubs:      make(map[string]any),
		pubsMutex: new(sync.Mutex),
	}

	clientIdFlag := flag.String("id", "", "The client id to use (each client must use a unique id)")
	flag.Parse()
	clientId := *clientIdFlag
	if clientId == "" {
		log.Fatal("client id must be set")
	}
	ui.clinetId = clientId

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://pi1.lan:1883").SetClientID(clientId)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(defaultPublishHandler)
	opts.SetPingTimeout(1 * time.Second)
	ui.mqtt = mqtt.NewClient(opts)
	if token := ui.mqtt.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	go func() {
		ui.window = app.NewWindow(app.Size(unit.Dp(1920), unit.Dp(1080)))
		if err := ui.run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (ui *Ui) run() error {
	ui.theme = material.NewTheme(gofont.Collection())
	grid.Init(ui)

	var ops op.Ops

	for {
		e := <-ui.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			grid.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

type Grid struct {
	Theme         *material.Theme
	Rows, Columns int
	Cells         [][]GridCellInterface
}

func (g *Grid) Init(ui *Ui) {
	var columns int
	for _, row := range g.Cells {
		if len(row) > columns {
			columns = len(row)
		}
		for _, cell := range row {
			cell.init(ui)
		}
	}
	g.Rows = len(g.Cells)
	g.Columns = columns
}

type GridCell[C Widget] struct {
	Contents C
	Init     func(C, *Ui)
}

func (c *GridCell[C]) init(ui *Ui) {
	c.Init(c.Contents, ui)
}

func (c *GridCell[C]) Layout(gtx layout.Context) layout.Dimensions {
	return c.Contents.Layout(gtx)
}

type Widget interface {
	Layout(gtx layout.Context) layout.Dimensions
}

type GridCellInterface interface {
	init(*Ui)
	Layout(gtx layout.Context) layout.Dimensions
}

func (g *Grid) Layout(gtx layout.Context) layout.Dimensions {
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

type Button struct {
	value   float32
	count   int
	th      *material.Theme
	drag    bool
	changed func(value float32)
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
			prevCount := b.count
			b.drag = true
			positionOfset := e.Position.Y - offsetY
			b.value = 1.0 - (positionOfset / buttonHeight)
			switch {
			case b.value > 1.0:
				b.value = 1.0
			case b.value < 0.0:
				b.value = 0.0
			}
			b.count = int(b.value * 1024)
			if prevCount != b.count {
				// only call changed if count has actually changed
				b.changed(b.value)
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
	paint.ColorOp{Color: White(0.7)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	area1.Pop()

	innerShaded := image.Rect(
		int(offsetX),
		int(offsetY+buttonHeight*(1-b.value)),
		int(offsetX+buttonWidth),
		int(offsetY+buttonHeight),
	)
	area2 := clip.Rect(innerShaded).Push(gtx.Ops)
	paint.ColorOp{Color: Black(0.6)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	area2.Pop()

	{
		offset := op.Offset(f32.Pt(offsetX, offsetY)).Push(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(int(buttonWidth), int(buttonHeight)))
		title := material.Body1(b.th, fmt.Sprintf("Value: %d%%", int(b.value*100.0)))
		title.Color = White(1)
		title.Alignment = text.Middle
		title.Layout(gtx)
		offset.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

//	func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
//		ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
//		return layout.Center.Layout(gtx, material.H3(th, text).Layout)
//	}
var (
	background = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	red        = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	green      = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	blue       = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
)

func Black(f float32) color.NRGBA {
	return White(1 - f)
}

func White(f float32) color.NRGBA {
	return color.NRGBA{R: uint8(f * 255), G: uint8(f * 255), B: uint8(f * 255), A: 0xFF}
}

//
//// ColorBox creates a widget with the specified dimensions and color.
//func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
//	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
//	paint.ColorOp{Color: color}.Add(gtx.Ops)
//	paint.PaintOp{}.Add(gtx.Ops)
//	return layout.Dimensions{Size: size}
//}
