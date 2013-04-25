package main

import (
    "nes"
    "os"
    "log"
    "fmt"
    "video"
    "strconv"
)

func RenderTile(tile []byte, x int, y int, frame []byte) {
    first := tile[:8]
    second := tile[8:]

    for row:=0; row < 8; row++ {
        f, s := first[row], second[row]
        for col:=0; col < 8; col++ {
            mask := byte(0x80 >> uint(col))

            offset := ((y+row) * 320 + x + col) * 3
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
    frame := make([]byte, 320 * 80 * 3)

    for i:=0; i < 0x1000; i+=16 {
        tile := table[i:i+16]

        x_index := (i / 16) % 32
        y_index := (i / 16 / 32)

        RenderTile(tile, 1 + x_index * 8 + x_index * 2, 1 + y_index * 8 + y_index * 2, frame)
    }

    return frame
}

func main() {
    path := os.Args[1]

    var err error

    var bank int
    if bank, err = strconv.Atoi(os.Args[2]); err != nil {
        log.Fatal(err)
        return
    }

    var file *os.File
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

    fmt.Printf("num chr banks: %d\n", len(rom.ChrBanks))

    screen := video.NewVideo()
    screen.Init(640, 160)

    go func() {
        buffer := RenderPatternTables(rom.ChrBanks[bank])

        frame := new(video.Frame)
        frame.Data = buffer
        frame.Width = 320
        frame.Height = 80

        screen.Frames <- frame
    }()

    screen.Loop()
}

