package cpu

import "testing"

type TimingTest struct {
    run (func(*CPU) int)
    expected int
}

func testTiming(t *testing.T, name string, run (func(*CPU) int), expected int) {
    var p = NewCPU()
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)
    p.Reset()

    actual := run(p)

    if actual != expected {
        t.Errorf("%s failed", name)
        t.Errorf("Expected %d, got %d", expected, actual)
    }
}

func timeimpl(op func(*CPU)) (func(*CPU) int) {
    return func(p *CPU) int {
        op(p)
        return p.cycles
    }
}

func time(op func(*CPU, Address)) (func(*CPU) int) {
    return func(p *CPU) int {
        op(p, 0x0000)
        return p.cycles
    }
}

func timeMode(mode int) (func(*CPU) int){
    return func(p *CPU) int {
        addressing[mode](p)

        return p.cycles
    }
}

func TestAddressModeXPage(t *testing.T) {
    tests := map[string]TimingTest {
        "Absolute X crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0xff

                    // 0xff + 0x0001 == 0x100 which crosses a page boundary
                    p.Memory.Copy([]byte{0x01, 0x00}, 0x0000)
                    p.AbsoluteX()

                    return p.cycles
                },
                3,
            },
        "Absolute Y crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0xff

                    // 0xff + 0x0001 == 0x100 which crosses a page boundary
                    p.Memory.Copy([]byte{0x01, 0x00}, 0x0000)
                    p.AbsoluteY()

                    return p.cycles
                },
                3,
            },
        "Indirect indexed crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0xff

                    // 0xff + 0x0001 == 0x100 which crosses a page boundary
                    p.Memory.Copy([]byte{0x01, 0x02, 0x00}, 0x0000)
                    p.IndirectIndexed()

                    return p.cycles
                },
                4,
            },
    }

    for name, test := range tests {
        testTiming(t, name, test.run, test.expected)
    }
}

func TestAddressModeTiming(t *testing.T) {
    tests := map[string]TimingTest {
        "Immediate": TimingTest { timeMode(Immediate), 0 },
        "Zero page": TimingTest { timeMode(ZeroPage), 1 },
        "Zero page x": TimingTest { timeMode(ZeroPageX), 2 },
        "Zero page y": TimingTest { timeMode(ZeroPageY), 2 },
        "Absolute": TimingTest { timeMode(Absolute), 2 },
        "Absolute x": TimingTest { timeMode(AbsoluteX), 2 },
        "Absolute y": TimingTest { timeMode(AbsoluteY), 2 },
        "Indirect": TimingTest { timeMode(Indirect), 4 },
        "Indexed indirect": TimingTest { timeMode(IndexedIndirect), 4 },
        "Indirect indexed": TimingTest { timeMode(IndirectIndexed), 3},
        "Relative": TimingTest { timeMode(Relative), 1},
    }

    for name, test := range tests {
        testTiming(t, name, test.run, test.expected)
    }
}

func branchTimingTest(t *testing.T, name string, branch (func(*CPU, Address)), yes (func(*CPU)), no (func(*CPU))) {
    tests := map[string]TimingTest {
        (name + " no branch"): TimingTest {
            func(p *CPU) int {
                no(p)
                branch(p, 0x0000)
                return p.cycles
            },
            0,
        },
        (name + " branch"): TimingTest {
            func(p *CPU) int {
                yes(p)
                branch(p, 0x0000)
                return p.cycles
            },
            1,
        },
        (name + " branch, new page"): TimingTest {
            func(p *CPU) int {
                yes(p)
                branch(p, 0x0100)
                return p.cycles
            },
            2,
        },
    }

    for name, test := range tests {
        testTiming(t, name, test.run, test.expected)
    }
}

func TestBranchTiming(t *testing.T) {
    branchTimingTest(t, "Bcc", (*CPU).Bcc,
        func(p *CPU) { p.setCarryFlag(false) },
        func(p *CPU) { p.setCarryFlag(true) })

    branchTimingTest(t, "Bcs", (*CPU).Bcs,
        func(p *CPU) { p.setCarryFlag(true) },
        func(p *CPU) { p.setCarryFlag(false) })

    branchTimingTest(t, "Beq", (*CPU).Beq,
        func(p *CPU) { p.setZeroFlag(true) },
        func(p *CPU) { p.setZeroFlag(false) })

    branchTimingTest(t, "Bmi", (*CPU).Bmi,
        func(p *CPU) { p.setNegativeFlag(true) },
        func(p *CPU) { p.setNegativeFlag(false) })

    branchTimingTest(t, "Bpl", (*CPU).Bpl,
        func(p *CPU) { p.setNegativeFlag(false) },
        func(p *CPU) { p.setNegativeFlag(true) })

    branchTimingTest(t, "Bne", (*CPU).Bne,
        func(p *CPU) { p.setZeroFlag(false) },
        func(p *CPU) { p.setZeroFlag(true) })

    branchTimingTest(t, "Bvc", (*CPU).Bvc,
        func(p *CPU) { p.setOverflowFlag(false) },
        func(p *CPU) { p.setOverflowFlag(true) })

    branchTimingTest(t, "Bvs", (*CPU).Bvs,
        func(p *CPU) { p.setOverflowFlag(true) },
        func(p *CPU) { p.setOverflowFlag(false) })
}

func TestNoArgsNonBranchTiming(t *testing.T) {
    tests := map[string]TimingTest {
        "Brk": TimingTest { timeimpl((*CPU).Brk), 6 },
        "Clc": TimingTest { timeimpl((*CPU).Clc), 1 },
        "Cld": TimingTest { timeimpl((*CPU).Cld), 1 },
        "Cli": TimingTest { timeimpl((*CPU).Cli), 1 },
        "Clv": TimingTest { timeimpl((*CPU).Clv), 1 },
        "Dex": TimingTest { timeimpl((*CPU).Dex), 1 },
        "Dey": TimingTest { timeimpl((*CPU).Dey), 1 },
        "Inx": TimingTest { timeimpl((*CPU).Inx), 1 },
        "Iny": TimingTest { timeimpl((*CPU).Iny), 1 },
        "Nop": TimingTest { timeimpl((*CPU).Nop), 1 },
        "Pha": TimingTest { timeimpl((*CPU).Pha), 2 },
        "Php": TimingTest { timeimpl((*CPU).Php), 2 },
        "Pla": TimingTest { timeimpl((*CPU).Pla), 3 },
        "Plp": TimingTest { timeimpl((*CPU).Plp), 3 },
        "Rti": TimingTest { timeimpl((*CPU).Rti), 5 },
        "Rts": TimingTest { timeimpl((*CPU).Rts), 5 },
        "Sec": TimingTest { timeimpl((*CPU).Sec), 1 },
        "Sed": TimingTest { timeimpl((*CPU).Sed), 1 },
        "Sei": TimingTest { timeimpl((*CPU).Sei), 1 },
        "Tax": TimingTest { timeimpl((*CPU).Tax), 1 },
        "Tay": TimingTest { timeimpl((*CPU).Tay), 1 },
        "Tsx": TimingTest { timeimpl((*CPU).Tsx), 1 },
        "Txa": TimingTest { timeimpl((*CPU).Txa), 1 },
        "Txs": TimingTest { timeimpl((*CPU).Txs), 1 },
        "Tya": TimingTest { timeimpl((*CPU).Tya), 1 },
        "Lsr acc": TimingTest { timeimpl((*CPU).LsrAcc), 1 },
        "Asl acc": TimingTest { timeimpl((*CPU).AslAcc), 1 },
        "Rol acc": TimingTest { timeimpl((*CPU).RolAcc), 1 },
        "Ror acc": TimingTest { timeimpl((*CPU).RorAcc), 1 },
    }

    for name, test := range tests {
        testTiming(t, name, test.run, test.expected)
    }
}

func TestArgsNonBranchTiming(t *testing.T) {
    tests := map[string]TimingTest {
        "Adc": TimingTest { time((*CPU).Adc), 1 },
        "Sbc": TimingTest { time((*CPU).Sbc), 1 },
        "And": TimingTest { time((*CPU).And), 1 },
        "Asl": TimingTest { time((*CPU).Asl), 3 },
        "Bit": TimingTest { time((*CPU).Bit), 1 },
        "Cmp": TimingTest { time((*CPU).Cmp), 1 },
        "Cpx": TimingTest { time((*CPU).Cpx), 1 },
        "Cpy": TimingTest { time((*CPU).Cpy), 1 },
        "Dec": TimingTest { time((*CPU).Dec), 3 },
        "Inc": TimingTest { time((*CPU).Inc), 3 },
        "Eor": TimingTest { time((*CPU).Eor), 1 },
        "Jmp": TimingTest { time((*CPU).Jmp), 0 },
        "Jsr": TimingTest { time((*CPU).Jsr), 3 },
        "Lda": TimingTest { time((*CPU).Lda), 1 },
        "Sta": TimingTest { time((*CPU).Sta), 1 },
        "Stx": TimingTest { time((*CPU).Stx), 1 },
        "Sty": TimingTest { time((*CPU).Sty), 1 },
        "Ldx": TimingTest { time((*CPU).Ldx), 1 },
        "Ldy": TimingTest { time((*CPU).Ldy), 1 },
        "Lsr": TimingTest { time((*CPU).Lsr), 3 },
        "Ora": TimingTest { time((*CPU).Ora), 1 },
        "Rol": TimingTest { time((*CPU).Rol), 3 },
        "Ror": TimingTest { time((*CPU).Ror), 3 },
    }

    for name, test := range tests {
        testTiming(t, name, test.run, test.expected)
    }
}
