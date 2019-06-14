package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"

	"gopkg.in/yaml.v2"

	mqttclient "github.com/automatedhome/common/pkg/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type topicsMapping struct {
	Topics []struct {
		In       string `yaml:"in"`
		Out      string `yaml:"out"`
		Retained bool   `yaml:"retained,omitempty"`
	} `yaml:"topics"`
}

var mappings topicsMapping

func onMessage(client mqtt.Client, message mqtt.Message) {
	for _, entry := range mappings.Topics {
		if entry.In == message.Topic() {
			client.Publish(entry.Out, 0, entry.Retained, message.Payload())
			return
		}
	}
}

func main() {
	broker := flag.String("broker", "tcp://127.0.0.1:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	clientID := flag.String("clientid", "proxy", "A clientid for the connection")
	configFile := flag.String("config", "/config.yaml", "Path to a configuration file")
	flag.Parse()

	brokerURL, _ := url.Parse(*broker)

	log.Printf("Reading file from %s", *configFile)
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("File reading error: %v", err)
		return
	}

	err = yaml.Unmarshal(data, &mappings)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var topics []string
	for _, entry := range mappings.Topics {
		log.Printf("subscribe: %s\n", entry.In)
		topics = append(topics, entry.In)
	}

	mqttclient.New(*clientID, brokerURL, topics, onMessage)
	log.Printf("Connected to %s as %s and waiting for messages\n", *broker, *clientID)

	// wait forever
	select {}
}
