package cpu

import "testing"

func TestSec(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Sec()

    if !p.Carry() {
        t.Errorf("Carry flag not set")
    }
}

func TestSed(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Sed()

    if !p.Decimal() {
        t.Errorf("Decimal flag not set")
    }
}

func TestSei(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0x00
    p.Sei()

    if !p.InterruptDisable() {
        t.Errorf("Interrupt disable flag not set")
    }
}
