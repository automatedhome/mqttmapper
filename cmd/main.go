package main

import (
	"flag"
	"log"
	"net/url"

	mqttclient "github.com/automatedhome/flow-meter/pkg/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Pipe struct {
	In  string
	Out string
}

var pipesList = []Pipe{
	Pipe{
		In:  "evok/temp/28FF0A9171150270/value",
		Out: "solar/temperature/in",
	},
	Pipe{
		In:  "evok/temp/28FF1A181515019F/value",
		Out: "solar/temperature/out",
	},
	Pipe{
		In:  "evok/temp/28FF4C30041503A7/value",
		Out: "tank/temperature/up",
	},
	Pipe{
		In:  "evok/temp/28FF4D15151501C6/value",
		Out: "heater/temperature/in",
	},
	Pipe{
		In:  "evok/temp/28FF5AF502150270/value",
		Out: "heater/temperature/out",
	},
	Pipe{
		In:  "evok/temp/287CECBF060000DA/value",
		Out: "outside/temperature",
	},
	Pipe{
		In:  "evok/temp/28FF89DB06000034/value",
		Out: "inside/temperature",
	},
}

func onMessage(client mqtt.Client, message mqtt.Message) {
	for _, pipe := range pipesList {
		if pipe.In == message.Topic() {
			client.Publish(pipe.Out, 0, false, message.Payload())
			return
		}
	}
}

func main() {
	broker := flag.String("broker", "tcp://127.0.0.1:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	clientID := flag.String("clientid", "proxy", "A clientid for the connection")
	flag.Parse()

	brokerURL, _ := url.Parse(*broker)
	var topics []string

	for _, pipe := range pipesList {
		topics = append(topics, pipe.In)
	}

	mqttclient.New(*clientID, brokerURL, topics, onMessage)
	log.Printf("Connected to %s as %s and waiting for messages\n", *broker, *clientID)

	// wait forever
	select {}
}
