package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/unit"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const MQTT_DELAY = time.Millisecond * 250

func main() {

	application := &App{
		pubs:      make(map[string]any),
		pubsMutex: new(sync.Mutex),
	}

	clientIdFlag := flag.String("id", "", "The client id to use (each client must use a unique id)")
	flag.Parse()
	clientId := *clientIdFlag
	if clientId == "" {
		log.Fatal("client id must be set")
	}
	application.clientId = clientId

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://pi1.lan:1883").SetClientID(clientId)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	})
	opts.SetPingTimeout(1 * time.Second)
	application.mqtt = mqtt.NewClient(opts)
	if token := application.mqtt.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	go func() {
		application.window = app.NewWindow(app.Size(unit.Dp(1920), unit.Dp(1080)))
		if err := application.run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
