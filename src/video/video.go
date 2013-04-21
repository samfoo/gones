package video

import (
    "log"
    "github.com/go-gl/gl"
    "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

type Frame struct {
    Data []byte
    Width int
    Height int
}

type Video struct {
    Width int
    Height int

    Screen *sdl.Surface
    Texture gl.Texture

    Frames chan *Frame
}

func NewVideo() *Video {
    v := new(Video)

    v.Frames = make(chan *Frame)

    return v
}

func (v *Video) HandleResize() {
    for {
        select {
            case e := <-sdl.Events:
                switch data := e.(type) {
                    case sdl.ResizeEvent:
                        v.Resize(int(data.W), int(data.H))
                }
        }
    }
}

func (v *Video) Resize(w int, h int) {
    v.Screen = sdl.SetVideoMode(w, h, 32, sdl.OPENGL)

    gl.Viewport(0, 0, w, h)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(-1, 1, -1, 1, -1, 1)
}

func (v *Video) Render(frame []byte, frame_w int, frame_h int) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    v.Texture.Bind(gl.TEXTURE_2D)

    gl.TexImage2D(gl.TEXTURE_2D, 0, 3, frame_w, frame_h, 0, gl.RGB, gl.UNSIGNED_BYTE, frame)

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
}

func (v *Video) Loop() {
    for {
        select {
            case frame := <-v.Frames:
                v.Render(frame.Data, frame.Width, frame.Height)
        }
    }
}

func (v *Video) Init(w int, h int) {
    sdl.Init(sdl.INIT_VIDEO)
    v.Screen = sdl.SetVideoMode(w, h, 32, sdl.OPENGL)
    if v.Screen == nil {
        log.Fatal(sdl.GetError())
    }

    if gl.Init() != 0 {
        log.Fatal(sdl.GetError())
    }

    gl.Enable(gl.TEXTURE_2D)

    v.Resize(w, h)

    v.Texture = gl.GenTexture()
}
