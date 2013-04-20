package ppu

type Pixel struct {
    Color uint32
    Transparent bool
}

type Display struct {
    Buffer []Pixel
}

func NewDisplay() *Display {
    d := new(Display)

    d.Buffer = make([]Pixel, 0xF000)

    return d
}
