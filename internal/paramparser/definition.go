package paramparser

// MapParams represents a struct for env/cmd line options
type MapParams struct {
	Hostname  string
	SubServer string
	SubTopic  string
	PubServer string
	PubTopic  string
	Qos       int
	ClientID  string
	Username  string
	Password  string
	Format    int
	PubFreq   int
	Expiry    int
	PubFile   string
}
