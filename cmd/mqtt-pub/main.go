package main

import (
	"bufio"
	"crypto/tls"
	"encoding/hex"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/yh742/j2735-decoder/internal/paramparser"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var logger log.Logger
var params parameters

type parameters struct {
	hostname string
	server   string
	subTopic string
	pubTopic string
	filename string
	qos      int
	clientid string
	username string
	password string
	format   int
	pubFreq  int
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	logger.Debugf("Received message on topic: %s", message.Topic())
	logger.Debugf("Message: %s", message.Payload())
}

func init() {
	logger = stdlog.GetFromFlags()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	hostname, _ := os.Hostname()

	// initialize struct
	params := &paramparser.MapParams{
		Hostname:  hostname,
		SubServer: "",
		PubServer: "",
		SubTopic:  "#",
		Qos:       0,
		ClientID:  hostname + strconv.Itoa(time.Now().Second()),
		Username:  "",
		Password:  "",
		PubFreq:   5,
		Expiry:    1,
		PubFile:   "",
		Format:    0,
	}
	// get parameters from (1) environment then (2) command line
	paramparser.Get(params)

	// print out flags
	logger.Debug("Initializing client with following parameters")
	logger.Debug("Hostname: ", params.Hostname)
	logger.Debug("Publish Server: ", params.PubServer)
	logger.Debug("SubTopic: ", params.SubTopic)
	logger.Debug("PubTopic: ", params.PubTopic)
	logger.Debug("ClientID: ", params.ClientID)
	logger.Debug("Username: ", params.Username)
	logger.Debug("Password: ", params.Password)
	logger.Debug("Publish Frequency: ", params.PubFreq)
	logger.Debug("Publish File: ", params.PubFile)

	if params.PubServer == "" {
		logger.Error("Must specify a server to connect to")
		os.Exit(2)
	}

	connOpts := MQTT.NewClientOptions().AddBroker(params.PubServer).SetClientID(params.ClientID).SetCleanSession(true)
	if params.Username != "" {
		logger.Debug("Username and password specfied")
		connOpts.SetUsername(params.Username)
		if params.Password != "" {
			connOpts.SetPassword(params.Password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(params.SubTopic, byte(params.Qos), onMessageReceived); token.Wait() && token.Error() != nil {
			logger.Error(token.Error())
			os.Exit(3)
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error(token.Error())
		os.Exit(4)
	}
	logger.Debugf("Connected to %s\n", params.PubServer)

	file, err := os.Open(params.PubFile)
	defer file.Close()
	if err != nil {
		logger.Error(err)
		os.Exit(5)
	}
	reader := bufio.NewReader(file)
	lineCnt := 0
	go func() {
		for true {
			time.Sleep(time.Duration(params.PubFreq * int(time.Millisecond)))
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				logger.Error("Something bad happened ....")
				continue
			}
			if err == io.EOF {
				logger.Debug("EOF reached resetting ...")
				file.Seek(0, 0)
				lineCnt = 0
				continue
			}
			splits := strings.Split(line, ":")
			hexString := strings.TrimSpace(splits[len(splits)-1])
			data, err := hex.DecodeString(hexString)
			if err != nil {
				logger.Debugf("publishing raw line %d: %s", lineCnt, line)
				client.Publish(params.PubTopic, byte(params.Qos), false, line)
			} else {
				logger.Debugf("line %d: %s", lineCnt, hexString)
				client.Publish(params.PubTopic, byte(params.Qos), false, data)
			}
			lineCnt++
		}
	}()
	// wait for control-c signal here
	<-c
}
