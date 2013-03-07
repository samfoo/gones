package cpu

import "testing"

func tax(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = value
    p.Tax()

    return p
}

func tay(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = value
    p.Tay()

    return p
}

func tsx(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.SP = value
    p.Tsx()

    return p
}

func txa(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.X = value
    p.Txa()

    return p
}

func txs(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.X = value
    p.Txs()

    return p
}

func tya(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.Y = value
    p.Tya()

    return p
}

func TestTransferNegativeFlag(t *testing.T) {
    if !tax(0x80).Negative() {
        t.Errorf("Negative flag not set for TAX")
    }

    if !tay(0x80).Negative() {
        t.Errorf("Negative flag not set for TAY")
    }

    if !tsx(0x80).Negative() {
        t.Errorf("Negative flag not set for TSX")
    }

    if !txa(0x80).Negative() {
        t.Errorf("Negative flag not set for TXA")
    }

    if !tya(0x80).Negative() {
        t.Errorf("Negative flag not set for TYA")
    }
}

func TestTransferZeroFlag(t *testing.T) {
    if !tax(0x00).Zero() {
        t.Errorf("Zero flag not set for TAX")
    }

    if !tay(0x00).Zero() {
        t.Errorf("Zero flag not set for TAY")
    }

    if !tsx(0x00).Zero() {
        t.Errorf("Zero flag not set for TSX")
    }

    if !txa(0x00).Zero() {
        t.Errorf("Zero flag not set for TXA")
    }

    if !tya(0x00).Zero() {
        t.Errorf("Zero flag not set for TYA")
    }
}

func TestTxsDoesntAffectFlags(t *testing.T) {
    if txs(0x00).Zero() {
        t.Errorf("TXS shouldn't affect zero flag")
    }

    if txs(0x80).Negative() {
        t.Errorf("TXS shouldn't affect negative flag")
    }
}

func TestTransfers(t *testing.T) {
    if tax(0xfa).X != 0xfa {
        t.Errorf("Transfering accumulator to X failed")
    }

    if tay(0xfa).Y != 0xfa {
        t.Errorf("Transfering accumulator to Y failed")
    }

    if tsx(0xfa).X != 0xfa {
        t.Errorf("Transfering SP to X failed")
    }

    if txa(0xfa).A != 0xfa {
        t.Errorf("Transfering X to accumulator failed")
    }

    if txs(0xfa).SP != 0xfa {
        t.Errorf("Transfering X to SP failed")
    }

    if tya(0xfa).A != 0xfa {
        t.Errorf("Transfering Y to accumulator failed")
    }
}

