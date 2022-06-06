package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	keys := make(chan string, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, _ := reader.ReadString('\n')
			keys <- s
		}
	}()

	fmt.Println("Mochi MQTT Server initializing...")

	server := mqtt.New()
	tcp := listeners.NewTCP("t1", ":1883")
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Add OnConnect Event Hook
	server.Events.OnConnect = func(cl events.Client, pk events.Packet) {
		fmt.Printf("<< OnConnect client connected %s: %+v\n", cl.ID, pk)
	}

	// Add OnDisconnect Event Hook
	server.Events.OnDisconnect = func(cl events.Client, err error) {
		fmt.Printf("<< OnDisconnect client dicconnected %s: %v\n", cl.ID, err)
	}

	// Add OnMessage Event Hook
	server.Events.OnMessage = func(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
		pkx = pk

		if pk.TopicName == "dimmer1" {
			fmt.Println(pk.TopicName, string(pk.Payload))
		}
		//if string(pk.Payload) == "hello" {
		//	pkx.Payload = []byte("hello world")
		//	fmt.Printf("< OnMessage modified message from client %s: %s\n", cl.ID, string(pkx.Payload))
		//} else {
		//	fmt.Printf("< OnMessage received message from client %s: %s: %s\n", cl.ID, pkx.TopicName, string(pkx.Payload))
		//}

		if pkx.TopicName == "accel_x" {
			fmt.Println(string(pkx.Payload))
		}

		// Example of using AllowClients to selectively deliver/drop messages.
		// Only a client with the id of `allowed-client` will received messages on the topic.
		if pkx.TopicName == "a/b/restricted" {
			pkx.AllowClients = []string{"allowed-client"} // slice of known client ids
		}

		return pkx, nil
	}

	go func() {
		for {
			<-keys
			for i := 0; i < 1000; i++ {
				server.Publish("dimmer1", []byte("5000"), false)
				//time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	// Demonstration of directly publishing messages to a topic via the
	// `server.Publish` method. Subscribe to `direct/publish` using your
	// MQTT client to see the messages.
	go func() {
		//for range time.Tick(time.Second * 10) {
		//server.Publish("light_on", []byte("1"), false)
		//time.Sleep(time.Second)
		//server.Publish("light_off", []byte("1"), false)
		//}
	}()

	fmt.Println("  Started!  ")

	<-done
	fmt.Println("  Caught Signal  ")

	server.Close()
	fmt.Println("  Finished  ")

}
