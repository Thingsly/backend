package simulationpublish

import (
	"net"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func PublishMessage(host string, port string, topic string, payload string, username string, password string, clientId string) error {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(net.JoinHostPort(host, port))
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)

	opts.SetCleanSession(true)

	opts.SetOrderMatters(false)

	opts.SetResumeSubs(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Println("simulation mqtt connect success")
	})
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logrus.Error("Simulation MQTT Broker connection failed:", token.Error())
		return token.Error()
	}
	defer c.Disconnect(250)
	logrus.Debug("username:", username)
	logrus.Debug("password:", password)
	logrus.Debug("clientId:", clientId)
	logrus.Debug("host:", host)
	logrus.Debug("port:", port)
	logrus.Debug("Topic:", topic)
	logrus.Debug("Payload:", payload)
	token := c.Publish(topic, 0, false, []byte(payload))
	if token.Wait() && token.Error() != nil {
		logrus.Error("Simulation MQTT Broker connection failed:", token.Error())
		return token.Error()
	}
	return nil
}
