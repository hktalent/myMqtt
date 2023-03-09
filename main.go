package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hktalent/mymqtt/pkg"
	"log"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	var broker = "127.0.0.1"
	//var port = 8083 // 1883
	opts := mqtt.NewClientOptions()
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	//client := pkg.ConnectByMQTT(pkg.Config{Host: broker, Port: 1883, Username: "admin", Password: "public___#$%^&*"}, opts)
	// client := pkg.ConnectByWS(pkg.Config{Host: broker, Port: 8083, Username: "admin", Password: "public___#$%^&*"}, opts)
	client := pkg.ConnectByWS(pkg.Config{Host: broker, Port: 8083, Username: "admin", Password: "public___#$%^&*"}, opts)
	//client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Println("client connect is ok")
	}
}
