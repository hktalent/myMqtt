package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hktalent/mymqtt/pkg"
	"io/ioutil"
	"log"
	"sync"
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

// https://www.emqx.com/en/blog/how-to-use-mqtt-in-golang
func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	// Import client certificate/key pair
	clientKeyPair, err := tls.LoadX509KeyPair("client-crt.pem", "client-key.pem")
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{clientKeyPair},
	}
}

type SubData struct {
	SubTopic string
	Callback mqtt.MessageHandler
}

func Sub(client mqtt.Client, a ...SubData) {
	var wg sync.WaitGroup
	for _, x := range a {
		wg.Add(1)
		go func(x1 *SubData) {
			defer wg.Done()
			topic := x1.SubTopic
			token := client.Subscribe(topic, 1, x1.Callback)
			token.Wait()
			fmt.Printf("Subscribed to topic %s\n", topic)
		}(&x)
	}
	wg.Wait()
}

type PublishData struct {
	SubTopic string
	Data     interface{}
}

func Publish(client mqtt.Client, a ...PublishData) {
	var wg sync.WaitGroup
	for _, x := range a {
		wg.Add(1)
		go func(x1 *PublishData) {
			defer wg.Done()
			token := client.Publish(x1.SubTopic, 0, false, x1.Data)
			token.Wait()
		}(&x)
	}
	wg.Wait()
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
