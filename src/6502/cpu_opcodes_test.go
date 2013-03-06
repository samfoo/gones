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

func testClearFlag(t *testing.T, name string, flag byte, opcode Opcode) {
    testOp(t, name, func(p *CPU) {
            p.Flags = flag
            p.execute(opcode, []byte{})
        }, func(p *CPU) bool { return p.Flags & flag == 0x00 })
}

func TestClearOpcodes(t *testing.T) {
    flags := map[string]byte {
        "Clc": 0x01,
        "Cli": 0x04,
        "Clv": 0x40,
    }

    tests := map[string]Opcode {
        "Clc": 0x18,
        "Cli": 0x58,
        "Clv": 0xb8,
    }

    for name, opcode := range tests {
        testClearFlag(t, name, flags[name], opcode)
    }
}

func TestBrkOpcode(t *testing.T) {
    testOp(t, "Brk implied", func(p *CPU) {
            p.Memory[0xfffe] = 0xff
            p.Memory[0xffff] = 0xee
            p.execute(0x00, []byte{})
        }, func(p *CPU) bool { return p.PC == 0xeeff })
}

func TestBitOpcodes(t *testing.T) {
    successful := func(p *CPU) bool {
        return p.Zero() && p.Overflow() && p.Negative()
    }

    tests := map[string]func(*CPU) {
        "Bit zero page":
            func(p *CPU) {
                p.A = 0x00
                p.execute(0x24, []byte{0x01, 0xff})
            },
        "Bit absolute":
            func(p *CPU) {
                p.A = 0x00
                p.execute(0x2c, []byte{0x02, 0x00, 0xff})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
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

func testOp(t * testing.T, name string, run func(*CPU), assertion func(*CPU) bool) {
    var p = new(CPU)
    p.Reset()
    run(p)

    if !assertion(p) {
        t.Errorf("%s failed", name)
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
    successful := func(p *CPU) bool { return p.A == 0x01 }

    tests := map[string]func(*CPU) {
        "And immediate":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x29, []byte{0x0f})
            },
        "And zero page":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x25, []byte{0x01, 0x0f})
            },
        "And zero page X":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x35, []byte{0x01, 0x00, 0x0f})
            },
        "And absolute":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x2d, []byte{0x02, 0x00, 0x0f})
            },
        "And absolute X":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x3d, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "And absolute Y":
            func(p *CPU) {
                p.A = 0x41
                p.Y = 0x01
                p.execute(0x39, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "And indexed indirect":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x21, []byte{0x00, 0x03, 0x00, 0x0f})
            },
        "And indirect indexed":
            func(p *CPU) {
                p.A = 0x41
                p.Y = 0x01
                p.execute(0x31, []byte{0x01, 0x00, 0x0f})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestAdcOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.A == 0x11 }

    tests := map[string]func(*CPU) {
        "Adc immediate":
            func(p *CPU) { p.execute(0x69, []byte{0x11}) },
        "Adc zero page":
            func(p *CPU) { p.execute(0x65, []byte{0x01, 0x11}) },
        "Adc zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x75, []byte{0x01, 0x00, 0x11})
            },
        "Adc absolute":
            func(p *CPU) { p.execute(0x6d, []byte{0x02, 0x00, 0x11}) },
        "Adc absolute X":
            func(p *CPU) {
                p.X = 0x02
                p.execute(0x7d, []byte{0x00, 0x00, 0x11})
            },
        "Adc absolute Y":
            func(p *CPU) {
                p.Y = 0x02
                p.execute(0x79, []byte{0x00, 0x00, 0x11})
            },
        "Adc indexed indirect":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x61, []byte{0x00, 0x03, 0x00, 0x11})
            },
        "Adc indirect indexed":
            func(p *CPU) {
                p.Y = 0x01
                p.execute(0x71, []byte{0x01, 0x00, 0x11})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}
