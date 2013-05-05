package main

import (
    "nes"
    "video"
    "os"
    "log"
)

func main() {
    path := os.Args[1]

    var file *os.File
    var err error
    if file, err = os.Open(path); err != nil {
        log.Fatal(err)
        return
    }

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    machine := nes.NewMachine()
    machine.Insert(rom)

    machine.CPU.Debug = false
    machine.CPU.Reset()

    machine.CPU.Cycle = func() {
        for i:=0; i < 3; i++ {
            machine.PPU.Step()
        }
    }

    screen := video.NewVideo()
    screen.Init(640, 600)

    go func() {
        for {
            machine.CPU.Step()

            frame := new(video.Frame)
            frame.Data = machine.PPU.Display
            frame.Width = 256
            frame.Height = 240

            screen.Frames <- frame
        }
    }()

    screen.Loop()
}

