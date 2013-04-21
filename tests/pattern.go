package main

import (
    "nes"
    "os"
    "log"
    "fmt"
    "video"
    "time"
    "runtime"
)

func RenderTile(tile []byte, x int, y int, frame []byte) {
    first := tile[:8]
    second := tile[8:]

    for i:=0; i < 8; i++ {
        f, s := first[i], second[i]
        for j:=0; j < 8; j++ {
            mask := byte(0x80 >> uint(j))

            offset := ((y+i) * 256 + x + j) * 3
            switch {
                case f & s & mask == mask:
                    frame[offset] = 0x00
                    frame[offset+1] = 0x00
                    frame[offset+2] = 0xff
                case f & mask == mask:
                    frame[offset] = 0xff
                    frame[offset+1] = 0x00
                    frame[offset+2] = 0x00
                case s & mask == mask:
                    frame[offset] = 0x00
                    frame[offset+1] = 0xff
                    frame[offset+2] = 0x00
                default:
                    frame[offset] = 0xff
                    frame[offset+1] = 0xff
                    frame[offset+2] = 0xff
            }
        }
    }
}

func RenderPatternTables(table []byte) []byte {
    frame := make([]byte, 256 * 256 * 3)

    for i:=0; i < 0x1000; i+=16 {
        tile := table[i:i+16]

        x_index := (i / 16) % 32
        y_index := (i / 16 / 32)

        RenderTile(tile, x_index * 8, y_index * 8, frame)
    }

    return frame
}

func main() {
    var file *os.File
    var err error
    if file, err = os.Open("assets/donkeykong.nes"); err != nil {
        log.Fatal(err)
        return
    }

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    fmt.Printf("num banks: %d\n", len(rom.ChrBanks))
    fmt.Printf("first bank: %d\n", len(rom.ChrBanks[0]))
    fmt.Printf("second bank: %d\n", len(rom.ChrBanks[1]))

    screen := video.NewVideo()
    screen.Init(1024, 1024)

    go func() {
        for i:=0;;i++ {
            var buffer []byte
            if i % 2 == 0 {
                buffer = RenderPatternTables(rom.ChrBanks[0])
            } else {
                buffer = RenderPatternTables(rom.ChrBanks[1])
            }

            frame := new(video.Frame)
            frame.Data = buffer
            frame.Width = 256
            frame.Height = 256

            screen.Frames <- frame
            time.Sleep(1000 * time.Millisecond)
        }
    }()

    runtime.LockOSThread()
    screen.Loop()
}

