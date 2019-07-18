package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

`	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"`
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var logger log.Logger

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	logger.Infof("Received message on topic: %s", message.Topic())
	logger.Infof("Message: %s", message.Payload())
}

func init() {
	logger = stdlog.GetFromFlags()
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	hostname, _ := os.Hostname()
	server := flag.String("server", "tcp://127.0.0.1:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	topic := flag.String("topic", "#", "Topic to subscribe to")
	qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	// print out flags
	logger.Debug("Server: ", *server)
	logger.Debug("Topic: ", *topic)
	logger.Debug("Clientid: ", *clientid)
	logger.Debug("Username: ", *username)
	logger.Debug("Password: ", *password)

	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *server)
	}

	<-c
}
