package main

import (
	"flag"
	"log"
	"net/url"

	mqttclient "github.com/automatedhome/flow-meter/pkg/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var topicsMap map[string]string

func onMessage(client mqtt.Client, message mqtt.Message) {
	for in, out := range topicsMap {
		if in == message.Topic() {
			client.Publish(out, 0, false, message.Payload())
			return
		}
	}
}

func init() {
	topicsMap = map[string]string{
		"evok/temp/28FF0A9171150270/value": "solar/temperature/in",
		"evok/temp/28FF1A181515019F/value": "solar/temperature/out",
		"evok/temp/28FF4C30041503A7/value": "tank/temperature/up",
		"evok/temp/28FF4D15151501C6/value": "heater/temperature/in",
		"evok/temp/28FF5AF502150270/value": "heater/temperature/out",
		"evok/temp/287CECBF060000DA/value": "climate/temperature/outside",
		"evok/temp/28FF89DB06000034/value": "climate/temperature/inside",
		// "solar/actuators/flow": "evok/ao/1/set",
		// "solar/actuators/pump": "evok/relay/3/set",
		// "solar/actuators/switch": "evok/relay/2/set",
		// "heater/actuators/burner": "evok/relay/5/set",
		// "heater/actuators/switch": "evok/relay/1/set"
	}
}

func main() {
	broker := flag.String("broker", "tcp://127.0.0.1:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	clientID := flag.String("clientid", "proxy", "A clientid for the connection")
	flag.Parse()

	brokerURL, _ := url.Parse(*broker)
	var topics []string

	for key := range topicsMap {
		topics = append(topics, key)
	}

	mqttclient.New(*clientID, brokerURL, topics, onMessage)
	log.Printf("Connected to %s as %s and waiting for messages\n", *broker, *clientID)

	// wait forever
	select {}
}
