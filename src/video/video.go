package video

import (
    "log"
    "runtime"
    "github.com/go-gl/gl"
    "github.com/go-gl/glfw"
)

type Frame struct {
    Data []byte
    Width int
    Height int
}

type Video struct {
    Width int
    Height int

    Texture gl.Texture

    Frames chan *Frame
    frame *Frame
}

func NewVideo() *Video {
    v := new(Video)

    v.Frames = make(chan *Frame)

    return v
}

func (v *Video) Render(frame []byte, frame_w int, frame_h int) {
    gl.Clear(gl.COLOR_BUFFER_BIT)
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

    glfw.SwapBuffers()
}

func (v *Video) Loop() {
    go func(c chan *Frame) {
        for {
            frame := <-c
            v.frame = frame
        }
    }(v.Frames)

    for glfw.WindowParam(glfw.Opened) == gl.TRUE {
        if v.frame != nil {
            v.Render(v.frame.Data, v.frame.Width, v.frame.Height)
        }
        runtime.Gosched()
    }

    glfw.Terminate()
}

func resize(w int, h int) {
    gl.Viewport(0, 0, w, h)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(-1, 1, -1, 1, -1, 1)
}

func (v *Video) Init(w int, h int) {
    var err error
    if err = glfw.Init(); err != nil {
        log.Fatal(err)
    }

    if err = glfw.OpenWindow(w, h, 8, 8, 8, 0, 24, 0, glfw.Windowed); err != nil {
        log.Fatal(err)
    }

    if gl.Init() != 0 {
        log.Fatal("ummm... hmmm")
    }

    glfw.SetWindowSizeCallback(resize)

    gl.Enable(gl.TEXTURE_2D)

    resize(w, h)

    v.Texture = gl.GenTexture()

}
