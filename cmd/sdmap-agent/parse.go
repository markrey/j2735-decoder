package main

import (
	"flag"
	"os"
	"strconv"
)

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
	params.format, _ = strconv.Atoi(getEnv("FORMAT", strconv.Itoa(params.format)))
	params.pubFreq, _ = strconv.Atoi(getEnv("PUBFREQ", strconv.Itoa(params.pubFreq)))
	params.expiry, _ = strconv.Atoi(getEnv("EXPIRY", strconv.Itoa(params.format)))
	// reverse bearing temp fix
	params.revBearing, _ = strconv.ParseBool(getEnv("REVBEARING", "0"))

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
	params.format = *flag.Int("format", params.format, "Decoding format of message")
	params.pubFreq = *flag.Int("pubFreq", params.pubFreq, "Publish frequency in 100ms increments")
	params.pubFreq = *flag.Int("expiry", params.expiry, "Defines when a message is considered stale")
	flag.Parse()
}
