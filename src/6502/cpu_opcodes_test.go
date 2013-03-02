package cpu

import "testing"

func (p *CPU) execute(op Opcode, arguments []byte) (*CPU) {
    for i := range arguments {
        p.Memory[i] = arguments[i]
    }

    p.Op(op)()

    return p
}

func TestAndOpcodes(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.execute(0x29, []byte{0x0f})
    if p.A != 0x01 {
        t.Errorf("And immediate failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.execute(0x25, []byte{0x01, 0x0f})
    if p.A != 0x01 {
        t.Errorf("And zero page failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.X = 0x01
    p.execute(0x35, []byte{0x01, 0x00, 0x0f})
    if p.A != 0x01 {
        t.Errorf("And zero page X failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.execute(0x2d, []byte{0x02, 0x00, 0x0f})
    if p.A != 0x01 {
        t.Errorf("And absolute failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.X = 0x01
    p.execute(0x3d, []byte{0x02, 0x00, 0x00, 0x0f})
    if p.A != 0x01 {
        t.Errorf("And absolute X failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.Y = 0x01
    p.execute(0x39, []byte{0x02, 0x00, 0x00, 0x0f})
    if p.A != 0x01 {
        t.Errorf("And absolute Y failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.X = 0x01
    p.execute(0x21, []byte{0x00, 0x03, 0x00, 0x0f})
    if p.A != 0x01 {
        t.Errorf("And indexed indirect failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x41
    p.Y = 0x01
    p.execute(0x31, []byte{0x01, 0x00, 0x0f})
    if p.A != 0x01 {
        t.Errorf("And indirect indexed failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }
}

func TestAdcOpcodes(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.execute(0x69, []byte{0x11})
    if p.A != 0x11 {
        t.Errorf("Adc immediate failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.execute(0x65, []byte{0x01, 0x11})
    if p.A != 0x11 {
        t.Errorf("Adc zero page failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.X = 0x01
    p.execute(0x75, []byte{0x01, 0x00, 0x11})
    if p.A != 0x11 {
        t.Errorf("Adc zero page X failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.execute(0x6d, []byte{0x02, 0x00, 0x11})
    if p.A != 0x11 {
        t.Errorf("Adc absolute failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.X = 0x02
    p.execute(0x7d, []byte{0x00, 0x00, 0x11})
    if p.A != 0x11 {
        t.Errorf("Adc absolute X failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.Y = 0x02
    p.execute(0x79, []byte{0x00, 0x00, 0x11})
    if p.A != 0x11 {
        t.Errorf("Adc absolute Y failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.X = 0x01
    p.execute(0x61, []byte{0x00, 0x03, 0x00, 0x11})
    if p.A != 0x11 {
        t.Errorf("Adc indexed indirect failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.Y = 0x01
    p.execute(0x71, []byte{0x01, 0x00, 0x11})
    if p.A != 0x11 {
        t.Errorf("Adc indirect indexed failed")
        t.Errorf("Expected %#02x, got %#02x", 0x11, p.A)
    }
}
