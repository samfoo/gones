package main

import (
    "nes"
    "os"
    "log"
    "fmt"
    "unsafe"
    "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

func InitVideo() *sdl.Surface {
    sdl.Init(sdl.INIT_VIDEO)
    screen := sdl.SetVideoMode(256, 256, 32, sdl.SWSURFACE)
    if screen == nil {
      log.Fatal(sdl.GetError())
    }

    return screen
}

func DrawPoint(x int32, y int32, value uint32, screen *sdl.Surface) {
  var pixel = uintptr(screen.Pixels)
  pixel += (uintptr)((y*screen.W)+x) * unsafe.Sizeof(value)

  var pu = unsafe.Pointer(pixel)
  var pp *uint32
  pp = (*uint32)(pu)
  *pp = value
}

const (
    WHITE = uint32(0xffffff)
    RED = uint32(0xff0000)
    GREEN = uint32(0x00ff00)
    BLUE = uint32(0x0000ff)
)


func DrawTile(tile []byte, x int, y int, screen *sdl.Surface) {
    first := tile[:8]
    second := tile[8:]

    for i:=0; i < 8; i++ {
        f, s := first[i], second[i]
        for j:=0; j < 8; j++ {
            mask := byte(0x80 >> uint(j))

            var color = WHITE
            switch {
                case f & s & mask == mask:
                    fmt.Printf("3")
                    color = BLUE
                case f & mask == mask:
                    fmt.Printf("1")
                    color = RED
                case s & mask == mask:
                    fmt.Printf("2")
                    color = GREEN
                default:
                    fmt.Printf(".")
            }

            DrawPoint(int32(x+j-1), int32(y+i), color, screen)
        }
        fmt.Printf("\n")
    }
    fmt.Printf("========")
}

func DrawPatternTables(table []byte, screen *sdl.Surface) {
    for i:=0; i < 0x1000; i+=16 {
        tile := table[i:i+16]

        fmt.Printf("===> %#0x\n", i)
        x_index := (i / 16) % 32
        y_index := (i / 16 / 32)

        DrawTile(tile, x_index * 8, y_index * 8, screen)
    }
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

    screen := InitVideo()

    /*DrawPatternTables(rom.ChrBanks[0], screen)*/
    DrawPatternTables(rom.ChrBanks[1], screen)
    screen.Flip()

    // Next step until the tests are finished.
    for {
    }
}

