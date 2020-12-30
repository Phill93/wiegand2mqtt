package main

import (
  "github.com/Phill93/wiegand2mqtt/src/com"
  "github.com/Phill93/wiegand2mqtt/src/log"
  "github.com/Phill93/wiegand2mqtt/src/mqtt"
  "github.com/Phill93/wiegand2mqtt/src/wiegand"
  "time"
)

func main() {
  log.Infof("wiegand2mqtt starting")
  log.Infof("starting internal communication")
  b := *com.GetBus()
  mqtt.Init(&b)
  wiegand.NewKeypad(&b)
  for {
    time.Sleep(1)
  }
}
