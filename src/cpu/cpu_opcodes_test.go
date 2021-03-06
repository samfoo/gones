package cpu

import "testing"

func branchTest(t *testing.T, op Opcode, branch func(*CPU), nobranch func(*CPU)) {
    var p = NewCPU()
    p.Reset()
    branch(p)
    p.execute(op, []byte{0x02})
    if p.PC != 0x03 {
        t.Errorf("%#02x did't branch when it should have", op)
        t.Errorf("Expected %#02x, got %#02x", 0x03, p.PC)
    }

    p = NewCPU()
    p.Reset()
    nobranch(p)
    p.execute(op, []byte{0x02})
    if p.PC != 0x01 {
        t.Errorf("%#02x branched when it shouldn't have", op)
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.PC)
    }
}

func testOp(t * testing.T, name string, run func(*CPU), assertion func(*CPU) bool) {
    var p = NewCPU()
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)
    p.Reset()
    run(p)

    if !assertion(p) {
        t.Errorf("%s failed", name)
    }
}

func (p *CPU) execute(op Opcode, arguments []byte) (*CPU) {
    for i := range arguments {
        p.Memory.Write(arguments[i], Address(i))
    }

    p.Execute(p.Operations()[op])

    return p
}

func testSetFlag(t *testing.T, name string, flag byte, opcode Opcode) {
    testOp(t, name, func(p *CPU) {
            p.Flags = flag
            p.execute(opcode, []byte{})
        }, func(p *CPU) bool { return p.Flags & flag == flag })
}

func testClearFlag(t *testing.T, name string, flag byte, opcode Opcode) {
    testOp(t, name, func(p *CPU) {
            p.Flags = flag
            p.execute(opcode, []byte{})
        }, func(p *CPU) bool { return p.Flags & flag == 0x00 })
}

func TestIncrementRegistersOpcodes(t *testing.T) {
    testOp(t, "Inx",
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xe8, []byte{})
            },
            func(p *CPU) bool {
                return p.X == 0x02
            })

    testOp(t, "Iny",
            func(p *CPU) {
                p.Y = 0x01
                p.execute(0xc8, []byte{})
            },
            func(p *CPU) bool {
                return p.Y == 0x02
            })
}

func TestIncOpcodes(t *testing.T) {
    successful := func(p *CPU) bool {
        return p.Memory.Read(0x0002) == 0xff
    }

    tests := map[string]func(*CPU) {
        "Inc zero page": func(p *CPU) { p.execute(0xe6, []byte{0x02, 0x00, 0xfe}) },
        "Inc zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xf6, []byte{0x01, 0x00, 0xfe})
            },
        "Inc absolute": func(p *CPU) { p.execute(0xee, []byte{0x02, 0x00, 0xfe}) },
        "Inc absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xfe, []byte{0x01, 0x00, 0xfe})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestDecrementRegistersOpcodes(t *testing.T) {
    testOp(t, "Dex",
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xca, []byte{})
            },
            func(p *CPU) bool {
                return p.X == 0x00
            })

    testOp(t, "Dey",
            func(p *CPU) {
                p.Y = 0x01
                p.execute(0x88, []byte{})
            },
            func(p *CPU) bool {
                return p.Y == 0x00
            })
}

func TestDecOpcodes(t *testing.T) {
    successful := func(p *CPU) bool {
        return p.Memory.Read(0x0002) == 0xfe
    }

    tests := map[string]func(*CPU) {
        "Dec zero page": func(p *CPU) { p.execute(0xc6, []byte{0x02, 0x00, 0xff}) },
        "Dec zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xd6, []byte{0x01, 0x00, 0xff})
            },
        "Dec absolute": func(p *CPU) { p.execute(0xce, []byte{0x02, 0x00, 0xff}) },
        "Dec absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xde, []byte{0x01, 0x00, 0xff})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestCompareOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.Carry() && p.Zero() }

    tests := map[string]func(*CPU) {
        "Cmp immediate":
            func(p *CPU) {
                p.A = 0x11
                p.execute(0xc9, []byte{0x11})
            },
        "Cmp zero page":
            func(p *CPU) {
                p.A = 0x11
                p.execute(0xc5, []byte{0x01, 0x11})
            },
        "Cmp zero page X":
            func(p *CPU) {
                p.A = 0x11
                p.X = 0x01
                p.execute(0xd5, []byte{0x01, 0x00, 0x11})
            },
        "Cmp absolute":
            func(p *CPU) {
                p.A = 0x11
                p.execute(0xcd, []byte{0x02, 0x00, 0x11})
            },
        "Cmp absolute X":
            func(p *CPU) {
                p.A = 0x11
                p.X = 0x02
                p.execute(0xdd, []byte{0x00, 0x00, 0x11})
            },
        "Cmp absolute Y":
            func(p *CPU) {
                p.A = 0x11
                p.Y = 0x02
                p.execute(0xd9, []byte{0x00, 0x00, 0x11})
            },
        "Cmp indexed indirect":
            func(p *CPU) {
                p.A = 0x11
                p.X = 0x01
                p.execute(0xc1, []byte{0x00, 0x03, 0x00, 0x11})
            },
        "Cmp indirect indexed":
            func(p *CPU) {
                p.A = 0x11
                p.Y = 0x01
                p.execute(0xd1, []byte{0x01, 0x02, 0x00, 0x11})
            },
        "Cpx immediate":
            func(p *CPU) {
                p.X = 0x11
                p.execute(0xe0, []byte{0x11})
            },
        "Cpx zero page":
            func(p *CPU) {
                p.X = 0x11
                p.execute(0xe4, []byte{0x01, 0x11})
            },
        "Cpx absolute":
            func(p *CPU) {
                p.X = 0x11
                p.execute(0xec, []byte{0x02, 0x00, 0x11})
            },
        "Cpy immediate":
            func(p *CPU) {
                p.Y = 0x11
                p.execute(0xc0, []byte{0x11})
            },
        "Cpy zero page":
            func(p *CPU) {
                p.Y = 0x11
                p.execute(0xc4, []byte{0x01, 0x11})
            },
        "Cpy absolute":
            func(p *CPU) {
                p.Y = 0x11
                p.execute(0xcc, []byte{0x02, 0x00, 0x11})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestSetOpcodes(t *testing.T) {
    flags := map[string]byte {
        "Sec": 0x01,
        "Sei": 0x04,
        "Sed": 0x08,
    }

    tests := map[string]Opcode {
        "Sec": 0x38,
        "Sei": 0x78,
        "Sed": 0xf8,
    }

    for name, opcode := range tests {
        testSetFlag(t, name, flags[name], opcode)
    }
}

func TestClearOpcodes(t *testing.T) {
    flags := map[string]byte {
        "Clc": 0x01,
        "Cli": 0x04,
        "Cld": 0x08,
        "Clv": 0x40,
    }

    tests := map[string]Opcode {
        "Clc": 0x18,
        "Cli": 0x58,
        "Clv": 0xb8,
        "Cld": 0xd8,
    }

    for name, opcode := range tests {
        testClearFlag(t, name, flags[name], opcode)
    }
}

func TestBrkOpcode(t *testing.T) {
    testOp(t, "Brk implied", func(p *CPU) {
            p.Memory.Write(0xff, 0xfffe)
            p.Memory.Write(0xee, 0xffff)
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

func TestShiftAndRotateOpcodes(t *testing.T) {
    testOp(t, "Asl accumulator", func(p *CPU) {
            p.A = 0x0f
            p.execute(0x0a, []byte{})
        }, func(p *CPU) bool { return p.A == 0x1e })

    testOp(t, "Lsr accumulator", func(p *CPU) {
            p.A = 0x3c
            p.execute(0x4a, []byte{})
        }, func(p *CPU) bool { return p.A == 0x1e })

    testOp(t, "Rol accumulator", func(p *CPU) {
            p.A = 0x0f
            p.execute(0x2a, []byte{})
        }, func(p *CPU) bool { return p.A == 0x1e })

    testOp(t, "Ror accumulator", func(p *CPU) {
            p.A = 0x3c
            p.execute(0x6a, []byte{})
        }, func(p *CPU) bool { return p.A == 0x1e })

    successful := func(p *CPU) bool {
        return p.Memory.Read(2) == 0x1e
    }


    tests := map[string]func(*CPU) {
        "Asl zero page":
            func(p *CPU) { p.execute(0x06, []byte{0x02, 0x00, 0x0f}) },
        "Asl zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x16, []byte{0x01, 0x00, 0x0f})
            },
        "Asl absolute":
            func(p *CPU) { p.execute(0x0e, []byte{0x02, 0x00, 0x0f}) },
        "Asl absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x1e, []byte{0x01, 0x00, 0x0f})
            },
        "Lsr zero page":
            func(p *CPU) { p.execute(0x46, []byte{0x02, 0x00, 0x3c}) },
        "Lsr zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x56, []byte{0x01, 0x00, 0x3c})
            },
        "Lsr absolute":
            func(p *CPU) { p.execute(0x4e, []byte{0x02, 0x00, 0x3c}) },
        "Lsr absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x5e, []byte{0x01, 0x00, 0x3c})
            },
        "Rol zero page":
            func(p *CPU) { p.execute(0x26, []byte{0x02, 0x00, 0x0f}) },
        "Rol zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x36, []byte{0x01, 0x00, 0x0f})
            },
        "Rol absolute":
            func(p *CPU) { p.execute(0x2e, []byte{0x02, 0x00, 0x0f}) },
        "Rol absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x3e, []byte{0x01, 0x00, 0x0f})
            },
        "Ror zero page":
            func(p *CPU) { p.execute(0x66, []byte{0x02, 0x00, 0x3c}) },
        "Ror zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x76, []byte{0x01, 0x00, 0x3c})
            },
        "Ror absolute":
            func(p *CPU) { p.execute(0x6e, []byte{0x02, 0x00, 0x3c}) },
        "Ror absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0x7e, []byte{0x01, 0x00, 0x3c})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestJumpOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.PC == 0xbeef }

    tests := map[string]func(*CPU) {
        "Jmp absolute":
            func(p *CPU) {
                p.execute(0x4c, []byte{0xef, 0xbe})
            },
        "Jmp indirect":
            func(p *CPU) {
                p.execute(0x6c, []byte{0x02, 0x00, 0xef, 0xbe})
            },
        "Jsr absolute":
            func(p *CPU) {
                p.execute(0x20, []byte{0xef, 0xbe})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestStoreOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.Memory.Read(4) == 0x0f }

    tests := map[string]func(*CPU) {
        "Sta zero page":
            func(p *CPU) {
                p.A = 0x0f
                p.execute(0x85, []byte{0x04})
            },
        "Sta zero page X":
            func(p *CPU) {
                p.A = 0x0f
                p.X = 0x01
                p.execute(0x95, []byte{0x03})
            },
        "Sta absolute":
            func(p *CPU) {
                p.A = 0x0f
                p.execute(0x8d, []byte{0x04, 0x00})
            },
        "Sta absolute X":
            func(p *CPU) {
                p.A = 0x0f
                p.X = 0x01
                p.execute(0x9d, []byte{0x03, 0x00})
            },
        "Sta absolute Y":
            func(p *CPU) {
                p.A = 0x0f
                p.Y = 0x01
                p.execute(0x99, []byte{0x03, 0x00})
            },
        "Sta indexed indirect":
            func(p *CPU) {
                p.A = 0x0f
                p.X = 0x01
                p.execute(0x81, []byte{0x00, 0x04})
            },
        "Sta indirect indexed":
            func(p *CPU) {
                p.A = 0x0f
                p.Y = 0x01
                p.execute(0x91, []byte{0x01, 0x03, 0x00})
            },
        "Stx zero page":
            func(p *CPU) {
                p.X = 0x0f
                p.execute(0x86, []byte{0x04})
            },
        "Stx zero page Y":
            func(p *CPU) {
                p.X = 0x0f
                p.Y = 0x01
                p.execute(0x96, []byte{0x03})
            },
        "Stx absolute":
            func(p *CPU) {
                p.X = 0x0f
                p.execute(0x8e, []byte{0x04, 0x00})
            },
        "Sty zero page":
            func(p *CPU) {
                p.Y = 0x0f
                p.execute(0x84, []byte{0x04})
            },
        "Sty zero page X":
            func(p *CPU) {
                p.Y = 0x0f
                p.X = 0x01
                p.execute(0x94, []byte{0x03})
            },
        "Sty absolute":
            func(p *CPU) {
                p.Y = 0x0f
                p.execute(0x8c, []byte{0x04, 0x00})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestLdyOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.Y == 0x0f }

    tests := map[string]func(*CPU) {
        "Ldy immediate":
            func(p *CPU) { p.execute(0xa0, []byte{0x0f}) },
        "Ldy zero page":
            func(p *CPU) { p.execute(0xa4, []byte{0x01, 0x0f}) },
        "Ldy zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xb4, []byte{0x01, 0x00, 0x0f})
            },
        "Ldy absolute":
            func(p *CPU) { p.execute(0xac, []byte{0x02, 0x00, 0x0f}) },
        "Ldy absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xbc, []byte{0x02, 0x00, 0x00, 0x0f})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestLdxOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.X == 0x0f }

    tests := map[string]func(*CPU) {
        "Ldx immediate":
            func(p *CPU) { p.execute(0xa2, []byte{0x0f}) },
        "Ldx zero page":
            func(p *CPU) { p.execute(0xa6, []byte{0x01, 0x0f}) },
        "Ldx zero page Y":
            func(p *CPU) {
                p.Y = 0x01
                p.execute(0xb6, []byte{0x01, 0x00, 0x0f})
            },
        "Ldx absolute":
            func(p *CPU) { p.execute(0xae, []byte{0x02, 0x00, 0x0f}) },
        "Ldx absolute Y":
            func(p *CPU) {
                p.Y = 0x01
                p.execute(0xbe, []byte{0x02, 0x00, 0x00, 0x0f})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestLdaOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.A == 0x0f }

    tests := map[string]func(*CPU) {
        "Lda immediate":
            func(p *CPU) { p.execute(0xa9, []byte{0x0f}) },
        "Lda zero page":
            func(p *CPU) { p.execute(0xa5, []byte{0x01, 0x0f}) },
        "Lda zero page X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xb5, []byte{0x01, 0x00, 0x0f})
            },
        "Lda absolute":
            func(p *CPU) { p.execute(0xad, []byte{0x02, 0x00, 0x0f}) },
        "Lda absolute X":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xbd, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "Lda absolute Y":
            func(p *CPU) {
                p.Y = 0x01
                p.execute(0xb9, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "Lda indexed indirect":
            func(p *CPU) {
                p.X = 0x01
                p.execute(0xa1, []byte{0x00, 0x03, 0x00, 0x0f})
            },
        "Lda indirect indexed":
            func(p *CPU) {
                p.Y = 0x01
                p.execute(0xb1, []byte{0x01, 0x02, 0x00, 0x0f})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestOraOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.A == 0x4f }

    tests := map[string]func(*CPU) {
        "Ora immediate":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x09, []byte{0x0f})
            },
        "Ora zero page":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x05, []byte{0x01, 0x0f})
            },
        "Ora zero page X":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x15, []byte{0x01, 0x00, 0x0f})
            },
        "Ora absolute":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x0d, []byte{0x02, 0x00, 0x0f})
            },
        "Ora absolute X":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x1d, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "Ora absolute Y":
            func(p *CPU) {
                p.A = 0x41
                p.Y = 0x01
                p.execute(0x19, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "Ora indexed indirect":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x01, []byte{0x00, 0x03, 0x00, 0x0f})
            },
        "Ora indirect indexed":
            func(p *CPU) {
                p.A = 0x41
                p.Y = 0x01
                p.execute(0x11, []byte{0x01, 0x02, 0x00, 0x0f})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestEorOpcodes(t *testing.T) {
    successful := func(p *CPU) bool { return p.A == 0x4e }

    tests := map[string]func(*CPU) {
        "Eor immediate":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x49, []byte{0x0f})
            },
        "Eor zero page":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x45, []byte{0x01, 0x0f})
            },
        "Eor zero page X":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x55, []byte{0x01, 0x00, 0x0f})
            },
        "Eor absolute":
            func(p *CPU) {
                p.A = 0x41
                p.execute(0x4d, []byte{0x02, 0x00, 0x0f})
            },
        "Eor absolute X":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x5d, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "Eor absolute Y":
            func(p *CPU) {
                p.A = 0x41
                p.Y = 0x01
                p.execute(0x59, []byte{0x02, 0x00, 0x00, 0x0f})
            },
        "Eor indexed indirect":
            func(p *CPU) {
                p.A = 0x41
                p.X = 0x01
                p.execute(0x41, []byte{0x00, 0x03, 0x00, 0x0f})
            },
        "Eor indirect indexed":
            func(p *CPU) {
                p.A = 0x41
                p.Y = 0x01
                p.execute(0x51, []byte{0x01, 0x02, 0x00, 0x0f})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
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
                p.execute(0x31, []byte{0x01, 0x02, 0x00, 0x0f})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}

func TestAddAndSubOpcodes(t *testing.T) {
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
                p.execute(0x71, []byte{0x01, 0x02, 0x00, 0x11})
            },
        "Sbc immediate":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.execute(0xe9, []byte{0x11})
            },
        "Sbc zero page":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.execute(0xe5, []byte{0x01, 0x11})
            },
        "Sbc zero page X":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.X = 0x01
                p.execute(0xf5, []byte{0x01, 0x00, 0x11})
            },
        "Sbc absolute":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.execute(0xed, []byte{0x02, 0x00, 0x11})
            },
        "Sbc absolute X":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.X = 0x02
                p.execute(0xfd, []byte{0x00, 0x00, 0x11})
            },
        "Sbc absolute Y":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.Y = 0x02
                p.execute(0xf9, []byte{0x00, 0x00, 0x11})
            },
        "Sbc indexed indirect":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.X = 0x01
                p.execute(0xe1, []byte{0x00, 0x03, 0x00, 0x11})
            },
        "Sbc indirect indexed":
            func(p *CPU) {
                p.setCarryFlag(true)
                p.A = 0x22
                p.Y = 0x01
                p.execute(0xf1, []byte{0x01, 0x02, 0x00, 0x11})
            },
    }

    for name, test := range tests {
        testOp(t, name, test, successful)
    }
}
