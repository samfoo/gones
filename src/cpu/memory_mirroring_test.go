package cpu

import "testing"

func TestSystemMemoryIsMirrored(t *testing.T) {
    var r = new(RAM)

    mirroredPages := []Address { 0x0000, 0x0800, 0x1000, 0x1800 }

    r.Write(0xf4, 0x0000)

    for _, p := range mirroredPages {
        if r.Read(p) != 0xf4 {
            t.Errorf("System memory not properly mirrored at %#04x", p)
        }
    }
}
