package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/dave/iot/ui/tool"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type App struct {
	theme     *material.Theme
	window    *app.Window
	clientId  string
	mqtt      mqtt.Client
	pubs      map[string]any
	pubsMutex *sync.Mutex
}

func (a *App) run() error {
	a.theme = material.NewTheme(gofont.Collection())
	root.Init(a)

	var ops op.Ops

	for {
		e := <-a.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			root.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) handle(b *tool.Slider) {
	b.Theme = a.theme
	b.Changed = func(value float32) {
		a.publish(fmt.Sprintf("%s/value", b.Id), ValueMessage{Value: value, Client: a.clientId})
	}
	a.mqtt.Subscribe(fmt.Sprintf("%s/value", b.Id), 0, func(client mqtt.Client, msg mqtt.Message) {
		var v ValueMessage
		if err := json.Unmarshal(msg.Payload(), &v); err != nil {
			fmt.Println(err)
			return
		}
		if v.Client == a.clientId {
			// ignore messages sent by this client
			return
		}
		b.Value = v.Value
		a.window.Invalidate()
	})
}

func (a *App) publish(topic string, message any) {
	a.pubsMutex.Lock()
	defer a.pubsMutex.Unlock()

	_, found := a.pubs[topic]
	if found {
		// we published recently to this topic, don't publish right now, but store the message
		a.pubs[topic] = message
		return
	}

	// publish the message immediately
	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.mqtt.Publish(topic, 0, true, b)
	a.pubs[topic] = nil
	go func() {
		// start a goroutine to wait for MQTT_DELAY before publishing the next message
		time.Sleep(MQTT_DELAY)
		a.pubsMutex.Lock()
		defer a.pubsMutex.Unlock()
		if message, found := a.pubs[topic]; found {
			if message == nil {
				delete(a.pubs, topic)
				return
			}
			b, err := json.Marshal(message)
			if err != nil {
				delete(a.pubs, topic)
				fmt.Println(err)
				return
			}
			a.mqtt.Publish(topic, 0, true, b)
			delete(a.pubs, topic)
		}
	}()
}
