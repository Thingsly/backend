package device

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/HustIoTPlatform/backend/mqtt/subscribe"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DeviceListener struct {
	redis     *redis.Client
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup sync.WaitGroup
}

func NewDeviceListener(redis *redis.Client) *DeviceListener {
	ctx, cancel := context.WithCancel(context.Background())
	return &DeviceListener{
		redis:  redis,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (l *DeviceListener) Start() error {

	if err := l.checkRedisConfig(); err != nil {
		return err
	}

	l.waitGroup.Add(1)
	go l.run()
	return nil
}

func (l *DeviceListener) checkRedisConfig() error {
	config, err := l.redis.ConfigGet(l.ctx, "notify-keyspace-events").Result()
	if err != nil {
		return fmt.Errorf("Failed to get Redis configuration: %v", err)
	}

	configValue := config["notify-keyspace-events"]
	if !strings.Contains(configValue, "Ex") {
		err = l.redis.ConfigSet(l.ctx, "notify-keyspace-events", "Ex").Err()
		if err != nil {
			return fmt.Errorf("Failed to set Redis configuration: %v", err)
		}
		logrus.Info("Redis expiration notification configuration updated")
	}
	return nil
}

func (l *DeviceListener) run() {
	defer l.waitGroup.Done()
	defer logrus.Info("Device listener has exited")

	dbNum := viper.GetInt("db.redis.db1")
	if dbNum == 0 {
		dbNum = 10
	}

	pubsub := l.redis.PSubscribe(l.ctx, fmt.Sprintf("__keyevent@%d__:expired", dbNum))
	defer pubsub.Close()

	if err := pubsub.Ping(l.ctx); err != nil {
		logrus.WithError(err).Error("Failed to subscribe to Redis expiration events")
		return
	}

	logrus.Infof("Device listener started successfully, listening for expired events on db%d", dbNum)

	ch := pubsub.Channel(redis.WithChannelSize(100))
	for {
		select {
		case <-l.ctx.Done():
			logrus.Info("Listener context canceled")
			return
		case msg, ok := <-ch:
			if !ok {
				logrus.Warn("Channel closed")
				return
			}

			if strings.HasPrefix(msg.Payload, "device:") &&
				(strings.HasSuffix(msg.Payload, ":heartbeat") ||
					strings.HasSuffix(msg.Payload, ":timeout")) {
				l.handleExpiredKey(msg)
			}
		}
	}
}

func (l *DeviceListener) handleExpiredKey(msg *redis.Message) {
	if msg == nil {
		return
	}

	keyParts := strings.Split(msg.Payload, ":")
	if len(keyParts) != 3 {
		logrus.WithField("payload", msg.Payload).
			Warn("Invalid key format")
		return
	}

	deviceID := keyParts[1]
	logrus.WithFields(logrus.Fields{
		"deviceID": deviceID,
		"type":     keyParts[2],
	}).Debug("Processing device status update")

	select {
	case <-l.ctx.Done():
		return
	default:
		subscribe.DeviceOnline([]byte("0"), "devices/status/"+deviceID)
	}
}

func (l *DeviceListener) Stop() error {
	logrus.Info("Stopping device listener...")

	l.cancel()

	done := make(chan struct{})
	go func() {
		l.waitGroup.Wait()
		close(done)
	}()

	select {
	case <-done:
		logrus.Info("Device listener stopped successfully")
		return nil
	case <-time.After(3 * time.Second):
		return errors.New("Device listener stop timed out")
	}
}
