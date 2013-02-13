package cpu

import "testing"

func validateGreaterThan(p *CPU, t *testing.T) {
    if !p.Carry() {
        t.Errorf("Carry flag should be set when {register} > {value} but isn't (flags: %08b)", p.Flags)
    }
}

func validateEqual(p *CPU, t *testing.T) {
    if !p.Carry() {
        t.Errorf("Carry flag should be set when {register} == {value} but isn't (flags: %08b)", p.Flags)
    }

    if !p.Zero() {
        t.Errorf("Zero flag should be set when {register} == {value} but isn't (flags: %08b)", p.Flags)
    }
}

func validateNegative(p *CPU, t *testing.T) {
    if !p.Negative() {
        t.Errorf("Negative flag should be set when {register} < {value} but isn't (flags: %08b)", p.Flags)
    }
}

func cmp(first byte, second byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = first
    p.Memory[0] = second
    p.Cmp(0)

    return p
}

func cpx(first byte, second byte) (*CPU) {
    var p *CPU = new(CPU)

    p.X = first
    p.Memory[0] = second
    p.Cpx(0)

    return p
}

func cpy(first byte, second byte) (*CPU) {
    var p *CPU = new(CPU)

    p.Y = first
    p.Memory[0] = second
    p.Cpy(0)

    return p
}

func TestCompareGreaterThan(t *testing.T) {
    validateGreaterThan(cmp(0x01, 0x00), t)
    validateGreaterThan(cpy(0x01, 0x00), t)
    validateGreaterThan(cpx(0x01, 0x00), t)
}

func TestCompareEqual(t *testing.T) {
    validateEqual(cmp(0x01, 0x01), t)
    validateEqual(cpy(0x01, 0x01), t)
    validateEqual(cpx(0x01, 0x01), t)
}

func TestCompareNegative(t *testing.T) {
    validateNegative(cmp(0x00, 0x01), t)
    validateNegative(cpy(0x00, 0x01), t)
    validateNegative(cpx(0x00, 0x01), t)
}
