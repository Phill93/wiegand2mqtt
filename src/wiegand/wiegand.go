package wiegand

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/Phill93/wiegand2mqtt/src/config"
  "github.com/Phill93/wiegand2mqtt/src/log"
  mqtt "github.com/Phill93/wiegand2mqtt/src/mqtt"
  "github.com/mustafaturan/bus"
  "time"
)

var low uint32 = 0x0
var high uint32 = 0x0
var carry uint32 = 0
var lowWd chan bool
var highWd chan bool
var end = false
var b bus.Bus
var ctx context.Context

type keypad struct {
	listeners map[string][]chan string
	code string
	codeLastUpdate time.Time
	codeTimeout time.Duration
	topic string
	gpioBeep int
	gpioLed int
	gpioLow int
	gpioHigh int
}

type keypadCode struct {
  Code string    `json:"code"`
  Time int64 `json:"time"`
}

type keypadCard struct {
  Card string `json:"card"`
  Time int64 `json:"time"`
}

func NewKeypad(b1 *bus.Bus) {
  var k keypad
  cfg := config.Config()
  conf := cfg.GetStringMap("keypad")
  if conf["timeout"] == nil {
    log.Warn("No keypad timeout set using default (10 sec)")
    k.codeTimeout, _ = time.ParseDuration(fmt.Sprintf("%ds", 10))
  } else {
    k.codeTimeout, _ = time.ParseDuration(fmt.Sprintf("%ds", conf["timeout"]))
  }

  if conf["topic"] == nil {
    log.Warn("No keypad topic set using default (keypad)")
    k.topic = "keypad"
  } else {
    k.topic = conf["topic"].(string)
  }

  loadPlatformConf(&k, conf)

  ctx = context.WithValue(context.Background(), bus.CtxKeyTxID, "wiegand")
  b = *b1

  mqtt.Subscribe(fmt.Sprintf("cmd/%s/beep", k.topic), 0, k.beep())
  mqtt.Subscribe(fmt.Sprintf("cmd/%s/led", k.topic), 0, k.led())

  InitReader(&k)
}

func (k *keypad) key(ke string) {
  if k.codeLastUpdate.IsZero() {
    k.codeLastUpdate = time.Now()
  }
  if checkTimeout(k.codeLastUpdate, k.codeTimeout.Seconds()) {
    log.Info("Timout for code entry reached, clearing old code!")
    k.clearCode()
  }
  switch ke {
  case "ENT":
    k.sendCode(k.code)
    k.clearCode()
  case "ESC":
    k.clearCode()
  default:
    k.code += ke
    k.codeLastUpdate = time.Now()
  }
}

func (k *keypad) clearCode()  {
  k.code = ""
  k.codeLastUpdate = time.Time{}
}

func (k *keypad) sendCode(c string){
  log.Infof("Got Code %s, try to send", c)
  data := keypadCode{Code: c, Time: time.Now().Unix()}
  dataJson, _ := json.Marshal(data)
  message := mqtt.NewMessage(fmt.Sprintf("state/%s/code", k.topic), string(dataJson))
  _, err := b.Emit(ctx, "mqtt.publish", message)
  if err != nil {
    log.Error(err)
  }
}

func (k *keypad) sendCard(c string){
  log.Debugf("Got Card %s, try to send", c)
  data := keypadCard{Card: c, Time: time.Now().Unix()}
  dataJson, _ := json.Marshal(data)
  message := mqtt.NewMessage(fmt.Sprintf("state/%s/card", k.topic), string(dataJson))
  _, err := b.Emit(ctx, "mqtt.publish", message)
  if err != nil {
    log.Error(err)
  }
}

func checkTimeout(timestamp time.Time, timeout float64) bool {
  now := time.Now()
  offset := now.Sub(timestamp)
  if offset.Seconds() > timeout {
    log.Debugf("Timeout reached! %f > %f", offset.Seconds(), timeout)
    return true
  } else {
    log.Debugf("Timeout not reached! %f !> %f", offset.Seconds(), timeout)
    return false
  }
}
