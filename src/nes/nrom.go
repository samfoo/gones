package nes

import "cpu"

type NROM struct {
    Rom *ROM
}

func (n *NROM) Graphics() cpu.Mountable {
    return &MountableStruct {
        func(location cpu.Address) byte {
            if location & 0x1000 == 0x1000 {
                return n.Rom.ChrBanks[1][location & 0xfff]
            } else {
                return n.Rom.ChrBanks[0][location]
            }
        },

        func(val byte, location cpu.Address) {
        },
    }
}

func (n *NROM) Program() cpu.Mountable {
    return &MountableStruct {
        func(location cpu.Address) byte {
            var bank = 0
            if len(n.Rom.PrgBanks) > 1 && location >= 0x4000 {
                bank = 1
            }

            // If NROM-128, mirror the first bank on the second
            normalized := location & 0x3fff
            return n.Rom.PrgBanks[bank][normalized]
        },

        func(val byte, location cpu.Address) {
        },
    }
}



