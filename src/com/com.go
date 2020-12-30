package com

import (
  "github.com/mustafaturan/bus"
  "github.com/mustafaturan/monoton"
  "github.com/mustafaturan/monoton/sequencer"
  "log"
  "time"
)

var b *bus.Bus

func GetBus() *bus.Bus {
  if b == nil {
    b = startBus()
  }
  return b
}

func startBus() *bus.Bus {
  node := uint64(1)
  initialTime := uint64(time.Now().Unix())
  m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime)
  if err != nil {
    log.Panic(err)
  }

  var idGenerator bus.Next = (*m).Next

  b, err := bus.NewBus(idGenerator)

  if err != nil {
    log.Panic(err)
  }

  return b
}

