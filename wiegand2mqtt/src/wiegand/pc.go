// +build !rpi

package wiegand

import (
  "github.com/Phill93/wiegand2mqtt/src/log"
  mqtt "github.com/eclipse/paho.mqtt.golang"
  "github.com/eiannone/keyboard"
	"os"
)

func loadPlatformConf(k *keypad, c map[string]interface{}) {
  log.Info("Loaded PC driver!")
}

func InitReader(pad *keypad) {
	keyEvents, err := keyboard.GetKeys(1)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	for {
		event := <-keyEvents
		if event.Err != nil {
			panic(event.Err)
		}

		switch string(event.Rune) {
		case "1":
		  log.Debug("Key 1 detected")
			pad.key("1")
		case "2":
      log.Debug("Key 2 detected")
			pad.key("2")
		case "3":
      log.Debug("Key 3 detected")
			pad.key("3")
		case "4":
      log.Debug("Key 4 detected")
			pad.key("4")
		case "5":
      log.Debug("Key 5 detected")
			pad.key("5")
		case "6":
      log.Debug("Key 6 detected")
			pad.key("6")
		case "7":
      log.Debug("Key 7 detected")
			pad.key("7")
		case "8":
      log.Debug("Key 8 detected")
			pad.key("8")
		case "9":
      log.Debug("Key 9 detected")
			pad.key("9")
		case "0":
      log.Debug("Key 0 detected")
			pad.key("0")
		}

		switch event.Key {
		case keyboard.KeyEsc:
      log.Debug("Key Esc detected")
			pad.key("ESC")
		case keyboard.KeyEnter:
      log.Debug("Key Ent detected")
			pad.key("ENT")
		case keyboard.KeyF12:
      log.Debug("Key F12 detected")
			os.Exit(0)
		}
	}
}

func (k *keypad) beep() mqtt.MessageHandler {
  return func(client mqtt.Client, message mqtt.Message) {
    log.Debugf("Received a beep message on topic %s, with payload %s", message.Topic(), message.Payload())
    log.Info("Beep!")
  }
}

func (k *keypad) led() mqtt.MessageHandler {
  return func(client mqtt.Client, message mqtt.Message) {
    log.Debugf("Received a led message on topic %s, with payload %s", message.Topic(), message.Payload())
    log.Info("Led!")
  }
}
