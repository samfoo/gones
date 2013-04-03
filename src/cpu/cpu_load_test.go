package cpu

import "testing"

func lda(value byte) (*CPU) {
    var p *CPU = NewCPU()

    p.Memory.buffer[0] = value
    p.Lda(0)

    return p
}

func ldx(value byte) (*CPU) {
    var p *CPU = NewCPU()

    p.Memory.buffer[0] = value
    p.Ldx(0)

    return p
}

func ldy(value byte) (*CPU) {
    var p *CPU = NewCPU()

    p.Memory.buffer[0] = value
    p.Ldy(0)

    return p
}

func TestLoadSetsRegisterValue(t *testing.T) {
    var p = lda(0xbe)
    if p.A != 0xbe {
        t.Errorf("LDA failed")
        t.Errorf("Expected 0xbe, got %#02x", p.A)
    }

    p = ldx(0xbe)
    if p.X != 0xbe {
        t.Errorf("LDX failed")
        t.Errorf("Expected 0xbe, got %#02x", p.X)
    }

    p = ldy(0xbe)
    if p.Y != 0xbe {
        t.Errorf("LDY failed")
        t.Errorf("Expected 0xbe, got %#02x", p.Y)
    }
}

func TestLoadZero(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Zero() {
            t.Errorf("Zero flag should be set when decrement is zero (flags: %08b)", p.Flags)
        }
    }

    validate(lda(0x00))
    validate(ldx(0x00))
    validate(ldy(0x00))
}

func TestLoadNegative(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Negative() {
            t.Errorf("Negative flag should be set when decrement is negative (flags: %08b)", p.Flags)
        }
    }

    validate(lda(0x80))
    validate(ldx(0x80))
    validate(ldy(0x80))
}
