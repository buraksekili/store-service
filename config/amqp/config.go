package amqp

import "os"

// AMQPConf represents needed configuration variables
type AMQPConf struct {
	Addr     string
	Exchange string
	Queue    string
}

// ExtractAMQPConfigs reads environment variables related to
// AMQP messaging. These environment variables can be read or
// updated from ./docker/.env file.
func ExtractAMQPConfigs() (ac AMQPConf) {
	if u := os.Getenv("S_AMQP_ADDRESS"); u != "" {
		ac.Addr = u
	}
	if en := os.Getenv("S_AMQP_EXCHANGE_NAME"); en != "" {
		ac.Exchange = en
	}
	if qn := os.Getenv("S_AMQP_QUEUE_NAME"); qn != "" {
		ac.Queue = qn
	}
	return ac
}
