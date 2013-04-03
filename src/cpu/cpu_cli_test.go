package cpu

import "testing"

func TestCliClearsInterruptDisable(t *testing.T) {
    var p *CPU = NewCPU()

    p.Flags = 0xff
    p.Cli()

    if p.InterruptDisable() {
        t.Errorf("Interrupt disable set when it should have been cleared (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
