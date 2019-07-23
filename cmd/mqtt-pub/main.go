package main

import (
	"crypto/tls"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"io/ioutil"

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
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	logger.Infof("Received message on topic: %s", message.Topic())
	logger.Infof("Message: %s", message.Payload())
}

func getEnv(key string, def string) string {
	variable := os.Getenv(key)
	if variable != "" {
		return variable
	}
	return def
}

func getParameters(params *parameters) {
	// get environment variables
	params.hostname = getEnv("HOSTNAME", params.hostname)
	params.server = getEnv("SERVER", params.server)
	params.filename = getEnv("FILENAME", params.filename)
	params.pubTopic = getEnv("PUBTOPIC", params.pubTopic)
	params.subTopic = getEnv("SUBTOPIC", params.subTopic)
	params.qos, _ = strconv.Atoi(getEnv("QOS", strconv.Itoa(params.qos)))
	params.clientid = getEnv("CLIENTID", params.clientid)
	params.username = getEnv("USERNAME", params.username)
	params.password = getEnv("PASSWORD", params.password)
	params.format, _ = strconv.Atoi(getEnv("FORMAT", strconv.Itoa(params.format)))
}

func init() {
	logger = stdlog.GetFromFlags()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	hostname, _ := os.Hostname()

	// initialize struct
	params = parameters{
		hostname: hostname,
		server:   "",
		subTopic: "#",
		filename: "",
		qos:      0,
		clientid: hostname + strconv.Itoa(time.Now().Second()),
		username: "",
		password: "",
	}
	// get parameters from (1) environment then (2) command line
	getParameters(&params)

	// print out flags
	logger.Debug("Initializing client with following parameters")
	logger.Debug("Hostname: ", params.hostname)
	logger.Debug("Server: ", params.server)
	logger.Debug("Filename: ", params.filename)
	logger.Debug("SubTopic: ", params.subTopic)
	logger.Debug("PubTopic: ", params.pubTopic)
	logger.Debug("Clientid: ", params.clientid)
	logger.Debug("Username: ", params.username)
	logger.Debug("Password: ", params.password)

	if params.server == "" {
		logger.Error("Must specify a server to connect to")
		os.Exit(2)
	}

	connOpts := MQTT.NewClientOptions().AddBroker(params.server).SetClientID(params.clientid).SetCleanSession(true)
	if params.username != "" {
		logger.Debug("Username and password specfied")
		connOpts.SetUsername(params.username)
		if params.password != "" {
			connOpts.SetPassword(params.password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(params.subTopic, byte(params.qos), onMessageReceived); token.Wait() && token.Error() != nil {
			logger.Error(token.Error())
			os.Exit(3)
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error(token.Error())
		os.Exit(4)
	}
	logger.Debugf("Connected to %s\n", params.server)

	dat, err := ioutil.ReadFile(params.filename)
	if err != nil {
		logger.Error(err)
		os.Exit(5)
	}
	for true {
		time.Sleep(2 * time.Second)
		client.Publish(params.pubTopic, byte(params.qos), false, dat)
	}
	// wait for control-c signal here
	<-c
}
