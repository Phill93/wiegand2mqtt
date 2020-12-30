package mqtt

import (
  "context"
  "fmt"
  "github.com/Phill93/wiegand2mqtt/src/config"
  "github.com/Phill93/wiegand2mqtt/src/log"
  mqtt "github.com/eclipse/paho.mqtt.golang"
  "github.com/mustafaturan/bus"
  "reflect"
  "time"
)

var client mqtt.Client
var b1 bus.Bus
var ctx context.Context
var timeout time.Duration

type message struct {
  Topic string
  Payload string
  Qos byte
  Retained bool
}

func NewMessage(t string, p string) *message{
  var m message
  m.Topic = t
  m.Payload = p
  m.Qos = 0
  m.Retained = false
  return &m
}

func (m *message) SetQos(q byte) {
  m.Qos = q
}

func (m *message) SetRetained(r bool) {
  m.Retained = r
}

func Init(b *bus.Bus) {
  b1 = *b
  b1.RegisterTopics("mqtt.receive")
  b1.RegisterTopics("mqtt.publish")

  ctx = context.WithValue(context.Background(), bus.CtxKeyTxID, "mqtt")

  opts := loadConf()
  opts.SetConnectionLostHandler(HandleConnectionLost)
  opts.SetOnConnectHandler(HandleConnection)
  opts.SetDefaultPublishHandler(HandleIncomingMessage)
  client = mqtt.NewClient(&opts)
  if token := client.Connect(); token.Wait() && token.Error() != nil {
    log.Panic(token.Error())
  }
  handler := bus.Handler{
    Handle: func(e *bus.Event) {
      if reflect.TypeOf(e.Data).String() == "*mqtt.message" {
        m := e.Data.(*message)
        go publish(m)
      } else {
        log.Error("Invalid data received!")
      }
    },
    Matcher: "mqtt.publish",
  }
  b.RegisterHandler("mqtt.publish", &handler)
}

func Subscribe(t string, q byte, c mqtt.MessageHandler) {
  log.Debugf("Try to subscribe to topic %s", t)
  token := client.Subscribe(t, q, c)
  if token.WaitTimeout(timeout) {
    log.Infof("Successful subscribed to topic %s", t)
  } else {
    if token.Error() == nil {
      log.Errorf("Timeout subscribing to topic %s", t)
    } else {
      log.Error(token.Error())
    }
  }
}

func makeUri(broker string, port int, protocol string) string {
  return fmt.Sprintf("%s://%s:%d", protocol, broker, port)
}

func loadConf() mqtt.ClientOptions {
  var port int
  var broker string
  var method string
  var user string
  var password string

  cfg := config.Config()

  conf := cfg.GetStringMap("mqtt")

  if conf["publishTimeout"] == nil {
    log.Warn("No publish timeout set using default (2 Seconds)")
    timeout, _ = time.ParseDuration(fmt.Sprintf("%ds", 2))
  } else {
    timeout, _ = time.ParseDuration(fmt.Sprintf("%ds", conf["publishTimeout"].(int)))
  }

  if conf["broker"] == nil {
    log.Panic("MQTT broker missing, please check your configuration")
  } else {
    broker = conf["broker"].(string)
  }

  if conf["port"] == nil {
    log.Warn("MQTT broker port missing, using default port")
    port = 1883
  } else {
    port = conf["port"].(int)
  }

  if conf["method"] == nil {
    log.Warn("Connection method missing, trying to figure it out")
    if port == 8883 {
      log.Debug("Detected mqtt over tls")
      method = "ssl"
    } else if port == 1883 {
      log.Debug("Detected plain mqtt")
      method = "tcp"
    } else if port == 80 {
      log.Debug("Detected mqtt over websocket")
      method = "ws"
    } else if port == 443 {
      log.Debug("Detected mqtt over secure websocket")
      method = "wss"
    } else {
      log.Panic("Could not recognize connection method!")
    }
  }

  if conf["user"] == nil {
    log.Warn("No mqtt user provided, disabling authentication")
  } else {
    if conf["password"] == nil {
      log.Panic("mqtt user is provided but no password")
    } else {
      user = conf["user"].(string)
      password = conf["password"].(string)
    }
  }

  if method == "ssl" || method == "ws" || method == "wss" {
    log.Panic("TLS and Websocket is not implemented jet")
  }

  opts := mqtt.NewClientOptions()
  opts.AddBroker(makeUri(broker, port, method))
  opts.SetClientID("wiegand2mqtt")
  if user != "" && password != "" {
    opts.SetUsername(user)
    opts.SetPassword(password)
  }

  return *opts
}

var HandleConnectionLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
  log.Warnf("MQTT Connection Lost, %v", err)
}

var HandleConnection mqtt.OnConnectHandler = func(client mqtt.Client) {
  log.Info("MQTT Connection established")
}

var HandleIncomingMessage mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
  log.Debugf("Received a message on topic %s with payload %s", message.Topic(), message.Payload())
  _, err := b1.Emit(ctx, "mqtt.receive", message)
  if err != nil {
    log.Error(err)
  }
}

func publish(m *message) {
  log.Debugf("Try to send a message to topic %s with payload %s", m.Topic, m.Payload)
  t := client.Publish(m.Topic, m.Qos, m.Retained, m.Payload)
  if t.WaitTimeout(timeout) {
    log.Infof("Message send successful to topic %s", m.Topic)
  } else {
    if t.Error() == nil {
      log.Error("Timeout sending message to topic %s", m.Topic)
    } else {
      log.Error(t.Error())
    }
  }
}
