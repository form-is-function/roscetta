package main

import (
  "log"

  "github.com/brutella/can"
  "github.com/hypebeast/go-osc/osc"
)

const CanInterfaceName = "can0"

var canFrameLock = can.Frame{
  ID:     0x077,
  Length: 1,
  Flags:  0,
  Res0:   0,
  Res1:   0,
  Data:   [8]uint8{0x05},
}


var canFrameUnlock = can.Frame{
  ID:     0x077,
  Length: 1,
  Flags:  0,
  Res0:   0,
  Res1:   0,
  Data:   [8]uint8{0x0A},
}

func main() {
  log.Println("Hello! I'm roscetta")
  defer func() {
    log.Println("Goodbye!")
  }()

  bus, _ := can.NewBusForInterfaceWithName(CanInterfaceName)
  err := bus.ConnectAndPublish()
  if err != nil {
    log.Fatalln(err)
  }

  addr := "127.0.0.1:9000"
  server := &osc.Server{Addr: addr}

  err = server.Handle("/door/unlock", func(msg *osc.Message) {
    err := bus.Publish(canFrameUnlock)
    if err != nil {
      log.Println(err)
    }
  })
  if err != nil {
    log.Fatalln("unable to setup handler:", err)
  }

  err = server.Handle("/door/lock", func(msg *osc.Message) {
    err := bus.Publish(canFrameLock)
    if err != nil {
      log.Println(err)
    }
  })
  if err != nil {
    log.Fatalln("unable to setup handler:", err)
  }

  bus.SubscribeFunc(func(frame can.Frame) {
    client := osc.NewClient("10.42.42.26", 8000)
    msg := osc.NewMessage("/any")
    msg.Append(int32(frame.ID))
    msg.Append(true)
    msg.Append(frame.Data[:frame.Length])
    err := client.Send(msg)
    if err != nil {
      log.Println(err)
    }
  })

  log.Fatalln(server.ListenAndServe())
}
