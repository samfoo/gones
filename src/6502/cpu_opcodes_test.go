package cpu

import "testing"

func branchTest(t *testing.T, op Opcode, branch func(*CPU), nobranch func(*CPU)) {
    var p = new(CPU)
    p.Reset()
    branch(p)
    p.execute(op, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("%#02x did't branch when it should have", op)
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }

    p = new(CPU)
    p.Reset()
    nobranch(p)
    p.execute(op, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("%#02x branched when it shouldn't have", op)
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }
}

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

func TestBvsOpcode(t *testing.T) {
    branchTest(t, 0x70,
        func(p *CPU) { p.Flags = 0x40 },
        func(p *CPU) { p.Flags = 0x00 },
    )
}

func TestBvcOpcode(t *testing.T) {
    branchTest(t, 0x50,
        func(p *CPU) { p.Flags = 0x00 },
        func(p *CPU) { p.Flags = 0x40 },
    )
}

func TestBmiOpcode(t *testing.T) {
    branchTest(t, 0x30,
        func(p *CPU) { p.Flags = 0x80 },
        func(p *CPU) { p.Flags = 0x00 },
    )
}

func TestBplOpcode(t *testing.T) {
    branchTest(t, 0x10,
        func(p *CPU) { p.Flags = 0x00 },
        func(p *CPU) { p.Flags = 0x80 },
    )
}

func TestBneOpcode(t *testing.T) {
    branchTest(t, 0xd0,
        func(p *CPU) { p.Flags = 0x00 },
        func(p *CPU) { p.Flags = 0x02 },
    )
}

func TestBeqOpcode(t *testing.T) {
    branchTest(t, 0xf0,
        func(p *CPU) { p.Flags = 0x02 },
        func(p *CPU) { p.Flags = 0x00 },
    )
}

func TestBcsOpcode(t *testing.T) {
    branchTest(t, 0xb0,
        func(p *CPU) { p.Flags = 0x01 },
        func(p *CPU) { p.Flags = 0x00 },
    )
}

func TestBccOpcode(t *testing.T) {
    branchTest(t, 0x90,
        func(p *CPU) { p.Flags = 0x00 },
        func(p *CPU) { p.Flags = 0x01 },
    )
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
