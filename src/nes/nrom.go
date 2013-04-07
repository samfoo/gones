package nes

import "cpu"

type NROM struct {
    Rom *ROM
}

func (n *NROM) Read(location cpu.Address) byte {
    var bank = 0
    if len(n.Rom.Banks) > 1 && location >= 0x4000 {
        bank = 1
    }

    // If NROM-128, mirror the first bank on the second
    normalized := location & 0x3fff
    return n.Rom.Banks[bank][normalized]
}

func (n *NROM) Write(val byte, location cpu.Address) {
}
