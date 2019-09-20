package main

import (
	"crypto/tls"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/yh742/j2735-decoder/pkg/decoder"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	CMap "github.com/orcaman/concurrent-map"
)

var logger log.Logger
var cmap *RWMap

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
	// reverse bearing temporary fix for cameras
	revBearing bool
}

// RWMap is concurrent map with read/write lock
type RWMap struct {
	mapInst CMap.ConcurrentMap
	mapLock sync.RWMutex
}

func addEntryToMap(id string, obj interface{}) {
	cmap.mapLock.RLock()
	cmap.mapInst.Set(id, obj)
	cmap.mapLock.RUnlock()
}

// fixBearing is only used as temporary adjustment for camera feeds
func onMessageReceived(format int, client MQTT.Client, message MQTT.Message, fixBearing bool) {
	logger.Infof("Received message on topic: %s", message.Topic())
	logger.Infof("Message: %s", message.Payload())
	decodedMsg, err := decoder.DecodeMapAgt(message.Payload(),
		uint(len(message.Payload())),
		decoder.MapAgentFormatType(format))
	if err != nil {
		if fixBearing {
			decodedMsg.SetHeading(-1*decodedMsg.GetHeading() + 28800)
		}
		addEntryToMap(decodedMsg.GetID(), decodedMsg)
		logger.Debugf("Msg ID: %s, Data: %+v", decodedMsg.GetID(), decodedMsg)
	}
}

func init() {
	logger = stdlog.GetFromFlags()
	cmap = &RWMap{
		mapInst: CMap.New(),
	}
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
			if token := c.Subscribe(topic, byte(qos), msgRcd); token.Wait() && token.Error() != nil {
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
		hostname:   hostname,
		subServer:  "",
		pubServer:  "",
		subTopic:   "#",
		qos:        0,
		clientid:   hostname + strconv.Itoa(time.Now().Second()),
		username:   "",
		password:   "",
		pubFreq:    5,
		revBearing: false,
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
	logger.Debug("Format: ", params.format)
	logger.Debug("Publish Frequency", params.pubFreq)
	// reversing bearing fix
	logger.Debug("Reversing Bearing", params.revBearing)

	if params.subServer == "" {
		logger.Error("Must specify a server to connect to")
		os.Exit(2)
	}
	subClient := createClient(params.clientid+"s",
		params.username, params.password, params.subServer,
		params.qos, params.subTopic, func(client MQTT.Client, message MQTT.Message) {
			onMessageReceived(params.format, client, message, params.revBearing)
		})
	logger.Debugf("Connected to %s\n", params.subServer)

	// create a seperate mqtt broker to publish to
	var pubClient MQTT.Client
	if params.pubServer != "" {
		pubClient = createClient(params.clientid+"p", params.username,
			params.password, params.pubServer, params.qos, "", nil)
	} else {
		pubClient = subClient
	}

	// frequency at which we publish the shared map
	duration := time.Duration(params.pubFreq * 100 * int(time.Millisecond))
	for {
		select {
		case <-time.After(duration):
			if !cmap.mapInst.IsEmpty() {
				cmap.mapLock.Lock()
				jsonBytes, err := cmap.mapInst.MarshalJSON()
				for _, key := range cmap.mapInst.Keys() {
					cmap.mapInst.Remove(key)
				}
				cmap.mapLock.Unlock()
				// publish
				go func() {
					if err == nil {
						pubClient.Publish(params.pubTopic, byte(params.qos), false, string(jsonBytes))
					}
				}()
			}
		case <-c:
			os.Exit(0)
		}
	}
}
