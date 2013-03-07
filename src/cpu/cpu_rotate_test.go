package cpu

import "testing"

func ror(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = value
    p.Ror(&p.A)

    return p
}

func rol(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = value
    p.Rol(&p.A)

    return p
}

func TestRotatingWithoutCarry(t *testing.T) {
    var p = ror(0x40)

    if p.A != 0x20 {
        t.Errorf("Rotating right without a carry failed")
        t.Errorf("Expected %#02x, got %#02x", 0x20, p.A)
    }

    p = rol(0x40)

    if p.A != 0x80 {
        t.Errorf("Rotating left without a carry failed")
        t.Errorf("Expected %#02x, got %#02x", 0x80, p.A)
    }
}

func TestRotatingWithCarry(t *testing.T) {
    var p = ror(0x01)

    if p.A != 0x80 {
        t.Errorf("Rotating right with a carry failed")
        t.Errorf("Expected %#02x, got %#02x", 0x80, p.A)
    }

    p = rol(0x80)

    if p.A != 0x01 {
        t.Errorf("Rotating leftwith a carry failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }
}

func TestRotatingSetsZero(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Zero() {
            t.Errorf("Rotating didn't set the zero flag")
        }
    }

    validate(rol(0x00))
    validate(ror(0x00))
}

func TestRotatingSetsNegative(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Negative() {
            t.Errorf("Rotating didn't set the negative flag")
        }
    }

    validate(rol(0x40))
    validate(ror(0x01))
}
