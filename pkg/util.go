package pkg

import (
	"crypto/tls"
	"crypto/x509"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"sync"
)

// https://www.emqx.com/en/blog/how-to-use-mqtt-in-golang
func NewTlsConfig(caPem, clientCrtPem, clientKeyPem string) *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caPem)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	// Import client certificate/key pair
	clientKeyPair, err := tls.LoadX509KeyPair(clientCrtPem, clientKeyPem)
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
			log.Printf("Subscribed to topic %s\n", topic)
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
