package cpu

import "testing"

func ror(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = value
    p.RorAcc()

    return p
}

func rol(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = value
    p.RolAcc()

    return p
}

func TestRorWithCarrySetBefore(t *testing.T) {
    var p *CPU = new(CPU)

    p.setCarryFlag(true)
    p.A = 0x00
    p.RorAcc()

    if p.A != 0x80 {
        t.Errorf("Rotating right with carry set should have set bit 7 to 1")
    }
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

    if p.A != 0x00 {
        t.Errorf("Rotating right with a carry failed")
        t.Errorf("Expected %#02x, got %#02x", 0x00, p.A)
    }

    if !p.Carry() {
        t.Errorf("Carry flag not set when rotating right")
    }

    p = rol(0x80)

    if p.A != 0x00 {
        t.Errorf("Rotating leftwith a carry failed")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }

    if !p.Carry() {
        t.Errorf("Carry flag not set when rotating left")
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
}
