package nes

import "cpu"

type MMC1 struct {
    Rom *ROM
}

func (m *MMC1) Read(location cpu.Address) byte {
    if location < 0x4000 {
        return m.Rom.Banks[0][location]
    }

    return m.Rom.Banks[len(m.Rom.Banks)-1][location & 0x3fff]
}

func (m *MMC1) Write(val byte, location cpu.Address) {
}
