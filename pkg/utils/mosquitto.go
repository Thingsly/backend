package utils

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-basic/uuid"
)

// mosquitto_pub -h xx.xx.xx.xx -p 1883 -t "devices/telemetry" -m "{\"tems\":112}" -u "xxxxx" -P "xxxxx" -i "0"
func BuildMosquittoPubCommand(host string, port string, username string, password string, topic string, payload string, clientId string) string {
	var sb strings.Builder
	sb.WriteString("mosquitto_pub")
	sb.WriteString(fmt.Sprintf(" -h %s", host))
	sb.WriteString(fmt.Sprintf(" -p %s", port))

	if topic != "" {
		sb.WriteString(fmt.Sprintf(" -t \"%s\"", topic))
	}
	if payload != "" {
		sb.WriteString(fmt.Sprintf(" -m \"%s\"", payload))
	}
	if username != "" {
		sb.WriteString(fmt.Sprintf(" -u \"%s\"", username))

	}
	if password != "" {
		sb.WriteString(fmt.Sprintf(" -P \"%s\"", password))
	}
	if clientId != "" {
		sb.WriteString(fmt.Sprintf(" -i \"%s\"", clientId))
	}
	return sb.String()
}

type MQTTParams struct {
	Host     string
	Port     string
	Username string
	Password string
	Topic    string
	Payload  string
	ClientId string
}

// Parse the mosquitto_pub command
// mosquitto_pub -h xx.xx.xx.xx -p 1883 -t "devices/telemetry" -m "{\"tems\":112}" -u "xxxxx" -P "xxxxx" -i "0"
func ParseMosquittoPubCommand(command string) (*MQTTParams, error) {
	args := strings.Split(command, " ")

	// Check if the command is "mosquitto_pub"
	if args[0] != "mosquitto_pub" {
		return nil, fmt.Errorf("invalid command: %s", args[0])
	}

	// Remove "mosquitto_pub"
	args = args[1:]

	f := flag.NewFlagSet("mqtt", flag.ContinueOnError)

	host := f.String("h", "localhost", "MQTT server address")
	port := f.String("p", "1883", "MQTT server port")
	user := f.String("u", "", "Username")
	password := f.String("P", "", "Password")
	topic := f.String("t", "", "MQTT topic")
	message := f.String("m", "", "Message content to publish")
	clientId := f.String("i", "", "Client ID")

	if *clientId == "" || *clientId == "0" {
		c := "mosquitto_pub_" + uuid.New()[0:8]
		clientId = &c
	}
	err := f.Parse(args)
	if err != nil {
		return nil, err
	}
	// Manually remove the quotes and escape characters around the argument values
	*host = strings.Trim(*host, "\"")
	*port = strings.Trim(*port, "\"")
	*user = strings.Trim(*user, "\"")
	*password = strings.Trim(*password, "\"")
	*topic = strings.Trim(*topic, "\"")
	*message, err = strconv.Unquote("\"" + strings.Trim(*message, "\"") + "\"")
	if err != nil {
		return nil, err
	}
	*clientId = strings.Trim(*clientId, "\"")

	params := &MQTTParams{
		Host:     *host,
		Port:     *port,
		Username: *user,
		Password: *password,
		Topic:    *topic,
		Payload:  *message,
		ClientId: *clientId,
	}

	return params, nil
}
