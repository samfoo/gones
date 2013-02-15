package cpu

import "testing"

func pha(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = value
    p.SP = 0xff
    p.Pha()

    return p
}

func php(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.Flags = value
    p.SP = 0xff
    p.Php()

    return p
}

func TestPushingValues(t *testing.T) {
    validate := func(p *CPU) {
        if p.SP != 0xfe {
            t.Errorf("Pushing didn't decrement the stack pointer correctly")
            t.FailNow()
        }

        head := 0x0100 + Address(p.SP) + 0x0001

        if p.Memory[head] != 0xbe {
            t.Errorf("Push didn't push to the top of the stack correctly")
            t.Errorf("Expected 0xbe, got %#02x at %#04x", p.Memory[head], head)
        }
    }

    validate(pha(0xbe))
    validate(php(0xbe))
}

func TestPushingMultipleValues(t *testing.T) {
    var p *CPU = new(CPU)
    p.SP = 0xff

    for i := 0; i < 5; i++ {
        p.A = byte(i)
        p.Pha()
    }

    if p.SP != 0xff - 0x05 {
        t.Errorf("Pushing didn't decrement the stack pointer correctly")
    }

    values := p.Memory[0x01fb:0x0200]

    for i := range values {
        if values[i] != 0x04 - byte(i) {
            t.Errorf("Invalid stack value at %#04x", 0x01fb - i)
            t.Errorf("Expected %#02x, got %#02x", 5 - i, values[i])
        }
    }
}

func TestPlaSetsZeroFlag(t *testing.T) {
    var p = pha(0x00)
    p.Pla()

    if !p.Zero() {
        t.Errorf("Pulling a zero accumulator didn't set the zero flag")
    }
}

func TestPlaSetsNegativeFlag(t *testing.T) {
    var p = pha(0x80)
    p.Pla()

    if !p.Negative() {
        t.Errorf("Pulling a negative accumulator didn't set the negative flag")
    }
}

func TestPullingValues(t *testing.T) {
    var p = pha(0xbe)
    p.A = 0x00
    p.Pla()

    if p.SP != 0xff {
        t.Errorf("Pulling didn't increment the stack pointer correctly")
    }

    if p.A != 0xbe {
        t.Errorf("Failed to pull from the stack into the accumulator.")
        t.Errorf("Expected %#02x, got %#02x", 0xbe, p.A)
    }

    p = php(0xfa)
    p.Flags = 0x00
    p.Plp()

    if p.SP != 0xff {
        t.Errorf("Pulling didn't increment the stack pointer correctly")
    }

    if p.Flags != 0xfa {
        t.Errorf("Failed to pull from the stack into flags.")
        t.Errorf("Expected %#02x, got %#02x", 0xbe, p.Flags)
    }
}
