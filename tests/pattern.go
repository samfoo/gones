package main

import (
    "nes"
    "os"
    "log"
    "fmt"
    "github.com/banthar/gl"
    "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

func InitVideo() *sdl.Surface {
    sdl.Init(sdl.INIT_VIDEO)
    screen := sdl.SetVideoMode(512, 512, 32, sdl.OPENGL)
    if screen == nil {
        log.Fatal(sdl.GetError())
    }

    if gl.Init() != 0 {
        log.Fatal(sdl.GetError())
    }

    gl.Enable(gl.TEXTURE_2D)

    gl.Viewport(0, 0, 512, 512)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(-1, 1, -1, 1, -1, 1)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    gl.Disable(gl.DEPTH_TEST)

    return screen
}

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

    InitVideo()
    texture := gl.GenTexture()

    /*frame := RenderPatternTables(rom.ChrBanks[0])*/
    frame := RenderPatternTables(rom.ChrBanks[1])

    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    texture.Bind(gl.TEXTURE_2D)

    gl.TexImage2D(gl.TEXTURE_2D, 0, 3, 256, 256, 0, gl.RGB, gl.UNSIGNED_BYTE, frame)

    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

    gl.Begin(gl.QUADS)
    gl.TexCoord2f(0.0, 1.0)
    gl.Vertex3f(-1.0, -1.0, 0.0)
    gl.TexCoord2f(1.0, 1.0)
    gl.Vertex3f(1.0, -1.0, 0.0)
    gl.TexCoord2f(1.0, 0.0)
    gl.Vertex3f(1.0, 1.0, 0.0)
    gl.TexCoord2f(0.0, 0.0)
    gl.Vertex3f(-1.0, 1.0, 0.0)
    gl.End()

    sdl.GL_SwapBuffers()

    for {
    }
}

