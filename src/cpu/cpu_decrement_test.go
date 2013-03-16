package cpu

import "testing"

func dec(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.Memory.buffer[0] = value
    p.Dec(0)

    return p
}

func dex(value byte) (*CPU) {
    var p *CPU = new(CPU)
    p.X = value
    p.Dex()
    return p
}

func dey(value byte) (*CPU) {
    var p *CPU = new(CPU)
    p.Y = value
    p.Dey()
    return p
}

func TestDecrementMutatesItem(t *testing.T) {
    var p = dec(0x02)
    if p.Memory.Read(0) != 0x01 {
        t.Errorf("Decrement memory failed")
        t.Errorf("Expected 0x01, got %#02x", p.Memory.Read(0))
    }

    p = dex(0x02)
    if p.X != 0x01 {
        t.Errorf("Decrement register X failed")
        t.Errorf("Expected 0x01, got %#02x", p.X)
    }

    p = dey(0x02)
    if p.Y != 0x01 {
        t.Errorf("Decrement register Y failed")
        t.Errorf("Expected 0x01, got %#02x", p.Y)
    }
}

func TestDecrementNonZero(t *testing.T) {
    var validate = func(p *CPU) {
        if p.Zero() {
            t.Errorf("Zero flag shouldn't be set when decrement isn't zero (flags: %08b)", p.Flags)
        }
    }

    validate(dec(0x02))
    validate(dex(0x02))
    validate(dey(0x02))
}

func TestDecrementZero(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Zero() {
            t.Errorf("Zero flag should be set when decrement is zero (flags: %08b)", p.Flags)
        }
    }

    validate(dec(0x01))
    validate(dex(0x01))
    validate(dey(0x01))
}

func TestDecrementNegative(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Negative() {
            t.Errorf("Negative flag should be set when decrement is negative (flags: %08b)", p.Flags)
        }
    }

    validate(dec(0x00))
    validate(dex(0x00))
    validate(dey(0x00))
}
