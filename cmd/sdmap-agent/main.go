package main

import (
	"crypto/tls"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/yh742/j2735-decoder/pkg/decoder"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"
	CMap"github.com/orcaman/concurrent-map"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var logger log.Logger
var cmap CMap.ConcurrentMap

type parameters struct {
	hostname  string
	subServer string
	subTopic  string
	pubServer string
	pubTopic  string
	qos       int
	clientid  string
	username  string
	password  string
	format    int
	pubFreq   int
}

func decodeRoutine(client MQTT.Client, message MQTT.Message) {
	decodedMsg := decoder.Decode(message.Payload(), uint(len(message.Payload())), decoder.SDMAP)
	sdData, ok := decodedMsg.(*decoder.SDMap)
	if ok {
		cmap.Set(sdData.ID, sdData)
	}
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	logger.Infof("Received message on topic: %s", message.Topic())
	logger.Infof("Message: %s", message.Payload())
	go decodeRoutine(client, message)
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
	params.subServer = getEnv("SUBSERVER", params.subServer)
	params.pubServer = getEnv("PUBSERVER", params.pubServer)
	params.subTopic = getEnv("SUBTOPIC", params.subTopic)
	params.pubTopic = getEnv("PUBTOPIC", params.pubTopic)
	params.qos, _ = strconv.Atoi(getEnv("QOS", strconv.Itoa(params.qos)))
	params.clientid = getEnv("CLIENTID", params.clientid)
	params.username = getEnv("USERNAME", params.username)
	params.password = getEnv("PASSWORD", params.password)
//	params.format, _ = strconv.Atoi(getEnv("FORMAT", strconv.Itoa(params.format)))
	params.pubFreq, _ = strconv.Atoi(getEnv("PUBFREQ", strconv.Itoa(params.pubFreq)))

	// command line overrides
	params.hostname = *flag.String("hostname", params.hostname, "The host machine name")
	params.subServer = *flag.String("subServer", params.subServer, "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	params.pubServer = *flag.String("pubServer", params.pubServer, "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	params.subTopic = *flag.String("subTopic", params.subTopic, "Topic to subscribe to")
	params.pubTopic = *flag.String("pubTopic", params.pubTopic, "Topic to publish to")
	params.qos = *flag.Int("qos", params.qos, "The QoS to publish/subscribe to messages at")
	params.clientid = *flag.String("clientid", params.clientid, "A clientid for the connection")
	params.username = *flag.String("username", params.username, "A username to authenticate to the MQTT server")
	params.password = *flag.String("password", params.password, "Password to match username")
//	params.format = *flag.Int("format", params.format, "Decoding format of message")
	params.pubFreq = *flag.Int("pubFreq", params.pubFreq, "Publish frequency in 100ms increments")
	flag.Parse()
}

func init() {
	logger = stdlog.GetFromFlags()
	cmap  = CMap.New()
}

func createClient(clientid string, username string, password string, server string, qos int, topic string,
	msgRcd func(client MQTT.Client, message MQTT.Message)) MQTT.Client {
	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientid).SetCleanSession(true)
	// set password
	if username != "" {
		logger.Debug("Username and password specfied")
		connOpts.SetUsername(username)
		if password != "" {
			connOpts.SetPassword(password)
		}
	}
	// set TLS config
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)
	// set connection callback
	connOpts.OnConnect = func(c MQTT.Client) {
		if topic != "" {
			if token := c.Subscribe(topic, byte(qos), msgRcd); 
				token.Wait() && token.Error() != nil {
				logger.Error(token.Error())
				os.Exit(3)
			}
		}
	}
	// create client
	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error(token.Error())
		os.Exit(4)
	}
	return client
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	hostname, _ := os.Hostname()

	// initialize struct
	params := &parameters{
		hostname: hostname,
		subServer:  "",
		pubServer: "",
		subTopic: "#",
		qos:      0,
		clientid: hostname + strconv.Itoa(time.Now().Second()),
		username: "",
		password: "",
		pubFreq: 5,
	}
	// get parameters from (1) environment then (2) command line
	getParameters(params)

	// print out flags
	logger.Debug("Initializing client with following parameters")
	logger.Debug("Hostname: ", params.hostname)
	logger.Debug("Subscribe Server: ", params.subServer)
	logger.Debug("Publish Server: ", params.pubServer)
	logger.Debug("SubTopic: ", params.subTopic)
	logger.Debug("PubTopic: ", params.pubTopic)
	logger.Debug("Clientid: ", params.clientid)
	logger.Debug("Username: ", params.username)
	logger.Debug("Password: ", params.password)
//	logger.Debug("Format: ", params.format)
	logger.Debug("Publish Frequency", params.pubFreq)

	if params.subServer == "" {
		logger.Error("Must specify a server to connect to")
		os.Exit(2)
	}
	subClient := createClient(params.clientid + "s", 
		params.username, params.password, params.subServer, 
		params.qos, params.subTopic, onMessageReceived)
	logger.Debugf("Connected to %s\n", params.subServer)
	
	// create a seperate mqtt broker to publish to
	var pubClient MQTT.Client
	if params.pubServer != "" {
		pubClient = createClient(params.clientid + "p", params.username, 
			params.password, params.pubServer, params.qos, "", nil)
	} else {
		pubClient = subClient
	}

	// frequency at which we publish the shared map
	duration := time.Duration(params.pubFreq * 100 * int(time.Millisecond))
	for {
		select {
		case <- time.After(duration):
			if cmap.Count() != 0 {
				jsonBytes, err := cmap.MarshalJSON()
				if err == nil {
					pubClient.Publish(params.pubTopic, byte(params.qos), false, string(jsonBytes))
				}
			}
		case <-c:
			os.Exit(0)
		}
	}
}
