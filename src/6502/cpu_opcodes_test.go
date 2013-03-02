package cpu

import "testing"

func (p *CPU) execute(op Opcode, arguments []byte) (*CPU) {
    for i := range arguments {
        p.Memory[i] = arguments[i]
    }

    p.Op(op)()

    return p
}

func TestBrkOpcode(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.Memory[0xfffe] = 0xff
    p.Memory[0xffff] = 0xee
    p.execute(0x00, []byte{})
    if p.PC != 0xeeff {
        t.Errorf("Brk didn't load the interrupt vector into PC")
        t.Errorf("Expected %#04x, got %#04x", 0xeeff, p.PC)
    }
}

func TestBitOpcodes(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.A = 0x00
    p.execute(0x24, []byte{0x01, 0xff})
    if !p.Zero() {
        t.Errorf("Zero flag not set right")
    }
    if !p.Overflow() || !p.Negative() {
        t.Errorf("Either the overflow or negative flag wasn't set right")
    }

    p = new(CPU)
    p.Reset()
    p.A = 0x00
    p.execute(0x2c, []byte{0x02, 0x00, 0xff})
    if !p.Zero() {
        t.Errorf("Zero flag not set right")
    }
    if !p.Overflow() || !p.Negative() {
        t.Errorf("Either the overflow or negative flag wasn't set right")
    }
}

func TestBmiOpcode(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.Flags = 0x80
    p.execute(0x30, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("Bmi with zero set didn't branch correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }

    p = new(CPU)
    p.Reset()
    p.Flags = 0x00
    p.execute(0x30, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("Bmi withput zero set branched")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }
}

func TestBplOpcode(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.Flags = 0x80
    p.execute(0x10, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("Bpl with zero set branched")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }

    p = new(CPU)
    p.Reset()
    p.Flags = 0x00
    p.execute(0x10, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("Bpl withput zero didn't branch correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }
}

func TestBneOpcode(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.Flags = 0x02
    p.execute(0xd0, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("Bne with zero set branched")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }

    p = new(CPU)
    p.Reset()
    p.Flags = 0x00
    p.execute(0xd0, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("Bne withput zero set didn't branch correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }
}

func TestBeqOpcode(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.Flags = 0x02
    p.execute(0xf0, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("Beq with zero set didn't branch correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }

    p = new(CPU)
    p.Reset()
    p.Flags = 0x00
    p.execute(0xf0, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("Beq withput zero branched")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }
}

func TestBcsOpcode(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.Flags = 0x01
    p.execute(0xb0, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("Bcs with carry set didn't branch correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }

    p = new(CPU)
    p.Reset()
    p.Flags = 0x00
    p.execute(0xb0, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("Bcs without carry branched")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }
}

func TestBccOpcode(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.Flags = 0x01
    p.execute(0x90, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("Bcc with carry set branched")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }

    p = new(CPU)
    p.Reset()
    p.Flags = 0x00
    p.execute(0x90, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("Bcc without carry set didn't branch correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }

    p = new(CPU)
    p.Reset()
    p.Flags = 0x00
    p.execute(0x90, []byte{0xff})
    if p.PC != 0x00 {
        t.Errorf("Bcc with negative offset didn't branch correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x00, p.PC)
    }
}

func TestAslOpcodes(t *testing.T) {
    var p = new(CPU)
    p.Reset()
    p.A = 0x0f
    p.execute(0x0a, []byte{})
    if p.A != 0x1e {
        t.Errorf("Asl accumulator failed")
        t.Errorf("Expected %#02x, got %#02x", 0x1e, p.A)
    }

    p = new(CPU)
    p.Reset()
    p.execute(0x06, []byte{0x01, 0x0f})
    if p.Memory[1] != 0x1e {
        t.Errorf("Asl zero page failed")
        t.Errorf("Expected %#02x, got %#02x", 0x1e, p.Memory[1])
    }

    p = new(CPU)
    p.Reset()
    p.X = 0x01
    p.execute(0x16, []byte{0x01, 0x00, 0x0f})
    if p.Memory[2] != 0x1e {
        t.Errorf("Asl zero page X failed")
        t.Errorf("Expected %#02x, got %#02x", 0x1e, p.Memory[2])
    }

    p = new(CPU)
    p.Reset()
    p.execute(0x0e, []byte{0x02, 0x00, 0x0f})
    if p.Memory[2] != 0x1e {
        t.Errorf("Asl absolute failed")
        t.Errorf("Expected %#02x, got %#02x", 0x1e, p.Memory[2])
    }

    p = new(CPU)
    p.Reset()
    p.X = 0x01
    p.execute(0x1e, []byte{0x02, 0x00, 0x00, 0x0f})
    if p.Memory[3] != 0x1e {
        t.Errorf("Asl absolute X failed")
        t.Errorf("Expected %#02x, got %#02x", 0x1e, p.Memory[3])
    }
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
