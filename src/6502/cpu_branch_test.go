package cpu

import "testing"

func TestBccBranchesIfNoCarry(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Bcc(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("BCC didn't branch when carry flag was not set")
        t.Errorf("Expected %#04x, got %#04x", 0x0025, p.PC)
    }
}

func TestBccDoesntBranchIfCarry(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x01
    p.Bcc(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("BCC branched when carry flag was set")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBcsNoBranchIfNoCarry(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Bcs(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("BCS branched when carry flag was not set")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBcsBranchesIfCarry(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x01
    p.Bcs(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("BCS didn't branch when carry flag was set")
        t.Errorf("Expected %#04x, got %#04x", 0x0025, p.PC)
    }
}

func TestBeqDoesntBranchIfNoZero(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Beq(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("BEQ branched when zero flag was not set")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBeqBranchesIfZero(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x02
    p.Beq(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("BEQ didn't branch when zero flag was set")
        t.Errorf("Expected %#04x, got %#04x", 0x0025, p.PC)
    }
}

func TestBmiBranchesIfNegative(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x80
    p.Bmi(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("Bmi didn't branch when negative flag was set")
        t.Errorf("Expected %#04x, got %#04x", 0x0025, p.PC)
    }
}

func TestBmiDoesntBranchIfNoNegative(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Bmi(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("Bmi branched when negative flag was not set")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBneBranchesIfNoZero(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Bne(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("BNE didn't branch when zero flag was not set")
        t.Errorf("Expected %#04x, got %#04x", 0x0025, p.PC)
    }
}

func TestBneDoesntBranchIfZero(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x02
    p.Bne(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("BNE branched when zero flag was set")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBplBranchesIfNotNegative(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Bpl(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("Bpl didn't branch when negative flag was not set")
        t.Errorf("Expected %#04x, got %#04x", 0x0025, p.PC)
    }
}

func TestBplDoesntBranchIfNegative(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x80
    p.Bpl(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("Bpl branched when negative flag was set")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBvcBranchIfOverflowClear(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Bvc(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("Bvc didn't branch when overflow clear")
        t.Errorf("Expected %#04x, got %#04x", 0x0025, p.PC)
    }
}

func TestBvcNoBranchIfOverflow(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x40
    p.Bvc(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("Bvc branched when overflow present")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBvsNoBranchIfOverflowClear(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Bvs(0x0025)

    if p.PC != 0x0000 {
        t.Errorf("Bvs branched when overflow clear")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.PC)
    }
}

func TestBvsBranchIfOverflow(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x40
    p.Bvs(0x0025)

    if p.PC != 0x0025 {
        t.Errorf("Bvs didn't branch when overflow present")
        t.Errorf("Expected %#04x, got %#04x", 0x0001, p.PC)
    }
}
