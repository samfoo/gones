package nes

import (
    "cpu"
    "ppu"
)

type MMC1 struct {
    Rom *ROM
}

func (n *MMC1) Patterntable(i int) *ppu.Patterntable {
    return ppu.NewPatterntable(make([]byte, 0x1000))
}

func (m *MMC1) Graphics() cpu.Mountable {
    return &MountableStruct {
        func(location cpu.Address) byte { return 0x00 },
        func(val byte, location cpu.Address) {},
    }
}

func (m *MMC1) Program() cpu.Mountable {
    return &MountableStruct {
        func(location cpu.Address) byte {
            if location < 0x4000 {
                return m.Rom.PrgBanks[0][location]
            }

            return m.Rom.PrgBanks[len(m.Rom.PrgBanks)-1][location & 0x3fff]
        },

        func(val byte, location cpu.Address) {
        },
    }
}

