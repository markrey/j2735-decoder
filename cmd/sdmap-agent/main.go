package main

import (
	"crypto/tls"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/yh742/j2735-decoder/internal/paramparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	CMap "github.com/orcaman/concurrent-map"
)

var logger log.Logger
var cmap *RWMap

// RWMap is concurrent map with read/write lock
type RWMap struct {
	mapInst CMap.ConcurrentMap
	expTbl  CMap.ConcurrentMap
	mapLock sync.RWMutex
	cleared bool
}

func addEntryToMap(id string, obj interface{}, expiryCnt int) {
	cmap.mapLock.RLock()
	cmap.mapInst.Set(id, obj)
	cmap.expTbl.Set(id, expiryCnt)
	cmap.mapLock.RUnlock()
}

// fixBearing is only used as temporary adjustment for camera feeds
func onMessageReceived(format int, client MQTT.Client, message MQTT.Message, expiryCnt int) { // , fixBearing bool) {
	logger.Infof("Received message on topic: %s", message.Topic())
	logger.Infof("Message: %X", message.Payload())
	// for decoding hex
	// data, err := hex.DecodeString(hexString)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }

	decodedMsg, err := decoder.DecodeMapAgt(message.Payload(),
		uint(len(message.Payload())),
		decoder.MapAgentFormatType(format))
	if err != nil {
		logger.Error(err)
	} else if decoder.MapAgentFormatType(format) == decoder.MAPSPaT {
		field, ok := decodedMsg.(*decoder.MapAgtSPaT)
		if !ok {
			logger.Error(err)
		}
		for _, intersectionstate := range field.IntersectionStateList {
			addEntryToMap(strconv.FormatUint(intersectionstate.ID, 10), intersectionstate, expiryCnt)
		}
	} else {
		addEntryToMap(decodedMsg.GetID(), decodedMsg, expiryCnt)
		logger.Debugf("Msg ID: %s, Data: %+v", decodedMsg.GetID(), decodedMsg)
	}
}

func init() {
	logger = stdlog.GetFromFlags()
	cmap = &RWMap{
		mapInst: CMap.New(),
		expTbl:  CMap.New(),
		cleared: false,
	}
}

func createClient(ClientID string, Username string, Password string, server string, Qos int, topic string,
	msgRcd func(client MQTT.Client, message MQTT.Message)) MQTT.Client {
	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(ClientID).SetCleanSession(true)
	// set Password
	if Username != "" {
		logger.Debug("Username and Password specfied")
		connOpts.SetUsername(Username)
		if Password != "" {
			connOpts.SetPassword(Password)
		}
	}
	// set TLS config
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)
	// set connection callback
	connOpts.OnConnect = func(c MQTT.Client) {
		if topic != "" {
			if token := c.Subscribe(topic, byte(Qos), msgRcd); token.Wait() && token.Error() != nil {
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
		//		RevBearing: false,
	}
	// get parameters from (1) environment then (2) command line
	paramparser.Get(params)

	// print out flags
	logger.Debug("Initializing client with following parameters")
	logger.Debug("Hostname: ", params.Hostname)
	logger.Debug("Subscribe Server: ", params.SubServer)
	logger.Debug("Publish Server: ", params.PubServer)
	logger.Debug("SubTopic: ", params.SubTopic)
	logger.Debug("PubTopic: ", params.PubTopic)
	logger.Debug("ClientID: ", params.ClientID)
	logger.Debug("Username: ", params.Username)
	logger.Debug("Password: ", params.Password)
	logger.Debug("Format: ", params.Format)
	logger.Debug("Publish Frequency: ", params.PubFreq)
	logger.Debug("Expiry: ", params.Expiry)
	// reversing bearing fix
	// logger.Debug("Reversing Bearing", params.RevBearing)

	if params.SubServer == "" {
		logger.Error("Must specify a server to connect to")
		os.Exit(2)
	}
	subClient := createClient(params.ClientID+"s",
		params.Username, params.Password, params.SubServer,
		params.Qos, params.SubTopic, func(client MQTT.Client, message MQTT.Message) {
			onMessageReceived(params.Format, client, message, params.Expiry) // , params.RevBearing)
		})
	logger.Debugf("Connected to %s\n", params.SubServer)

	// create a seperate mqtt broker to publish to
	var pubClient MQTT.Client
	if params.PubServer != "" {
		pubClient = createClient(params.ClientID+"p", params.Username,
			params.Password, params.PubServer, params.Qos, "", nil)
	} else {
		pubClient = subClient
	}

	// frequency at which we publish the shared map
	duration := time.Duration(params.PubFreq * int(time.Millisecond))
	for {
		select {
		case <-time.After(duration):
			if !cmap.mapInst.IsEmpty() {
				cmap.mapLock.Lock()
				cmap.cleared = false
				jsonBytes, err := cmap.mapInst.MarshalJSON()
				for _, key := range cmap.mapInst.Keys() {
					if tmp, ok := cmap.expTbl.Get(key); ok {
						expCnt := tmp.(int)
						if expCnt == 0 {
							cmap.expTbl.Remove(key)
							cmap.mapInst.Remove(key)
						} else {
							expCnt--
							cmap.expTbl.Set(key, expCnt)
						}
					} else {
						logger.Debugf("Key %s was not found in expiry table", key)
						cmap.mapInst.Remove(key)
					}
				}
				cmap.mapLock.Unlock()
				// publish
				go func() {
					if err == nil {
						logger.Debugf("Publishing: %s", string(jsonBytes))
						pubClient.Publish(params.PubTopic, byte(params.Qos), false, string(jsonBytes))
					}
				}()
			} else if !cmap.cleared {
				logger.Debug("Sending one more to clear map")
				cmap.mapLock.Lock()
				pubClient.Publish(params.PubTopic, byte(params.Qos), false, "{}")
				cmap.cleared = true
				cmap.mapLock.Unlock()
			}
		case <-c:
			os.Exit(0)
		}
	}
}
