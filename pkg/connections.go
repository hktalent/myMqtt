package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnectByMQTT(config Config, opts *mqtt.ClientOptions) mqtt.Client {
	if nil == opts {
		opts = mqtt.NewClientOptions()
	}
	broker := fmt.Sprintf("tcp://%s:%d", config.Host, config.Port)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func ConnectByMQTTS(config Config, opts *mqtt.ClientOptions) mqtt.Client {
	var tlsConfig tls.Config
	if config.Tls && config.CaCert == "" {
		log.Fatalln("TLS field in config is required")
	}
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.CaCert)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	tlsConfig.RootCAs = certpool

	if nil == opts {
		opts = mqtt.NewClientOptions()
	}
	broker := fmt.Sprintf("ssl://%s:%d", config.Host, config.Port)
	println(broker)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetTLSConfig(&tlsConfig)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func ConnectByWS(config Config, opts *mqtt.ClientOptions) mqtt.Client {
	if nil == opts {
		opts = mqtt.NewClientOptions()
	}
	broker := fmt.Sprintf("ws://%s:%d/mqtt", config.Host, config.Port)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func ConnectByWSS(config Config, opts *mqtt.ClientOptions) mqtt.Client {
	var tlsConfig tls.Config
	if config.Tls && config.CaCert == "" {
		log.Fatalln("TLS field in config is required")
	}
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.CaCert)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	tlsConfig.RootCAs = certpool

	if nil == opts {
		opts = mqtt.NewClientOptions()
	}
	broker := fmt.Sprintf("wss://%s:%d/mqtt", config.Host, config.Port)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetTLSConfig(&tlsConfig)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}
