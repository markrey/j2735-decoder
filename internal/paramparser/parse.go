package paramparser

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

// Get reads from:
// (1) environement variables
// (2) overrides with command line arguments
// and places the results in a parameter struct
func Get(params *MapParams) {
	// get environment variables
	params.Hostname = getEnv("HOSTNAME", params.Hostname)
	params.SubServer = getEnv("SUBSERVER", params.SubServer)
	params.PubServer = getEnv("PUBSERVER", params.PubServer)
	params.SubTopic = getEnv("SUBTOPIC", params.SubTopic)
	params.PubTopic = getEnv("PUBTOPIC", params.PubTopic)
	params.Qos, _ = strconv.Atoi(getEnv("QOS", strconv.Itoa(params.Qos)))
	params.ClientID = getEnv("CLIENTID", params.ClientID)
	params.Username = getEnv("USERNAME", params.Username)
	params.Password = getEnv("PASSWORD", params.Password)
	params.PubFile = getEnv("PUBFILE", params.PubFile)
	params.Format, _ = strconv.Atoi(getEnv("FORMAT", strconv.Itoa(params.Format)))
	params.PubFreq, _ = strconv.Atoi(getEnv("PUBFREQ", strconv.Itoa(params.PubFreq)))
	params.Expiry, _ = strconv.Atoi(getEnv("EXPIRY", strconv.Itoa(params.Expiry)))
	// reverse bearing temp fix
	// params.RevBearing, _ = strconv.ParseBool(getEnv("REVBEARING", "0"))

	// command line overrides
	params.Hostname = *flag.String("hostname", params.Hostname, "The host machine name")
	params.SubServer = *flag.String("subServer", params.SubServer, "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	params.PubServer = *flag.String("pubServer", params.PubServer, "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	params.SubTopic = *flag.String("subTopic", params.SubTopic, "Topic to subscribe to")
	params.PubTopic = *flag.String("pubTopic", params.PubTopic, "Topic to publish to")
	params.Qos = *flag.Int("qos", params.Qos, "The QoS to publish/subscribe to messages at")
	params.ClientID = *flag.String("clientid", params.ClientID, "A clientid for the connection")
	params.Username = *flag.String("username", params.Username, "A username to authenticate to the MQTT server")
	params.Password = *flag.String("password", params.Password, "Password to match username")
	params.Format = *flag.Int("format", params.Format, "Decoding format of message")
	params.PubFreq = *flag.Int("pubFreq", params.PubFreq, "Publish frequency in 100ms increments")
	params.Expiry = *flag.Int("expiry", params.Expiry, "Defines when a message is considered stale")
	params.PubFile = *flag.String("pubFile", params.PubFile, "Defines a file to playback message (debug only)")
	flag.Parse()
}
