package cpu

import "testing"

func TestNewCPUShouldHaveRegistersZerod(t *testing.T) {
    var p *CPU = new(CPU)

    if p.A != 0x00 {
        t.Errorf("Accumulator is not zero'd")
        t.FailNow()
    }
}

