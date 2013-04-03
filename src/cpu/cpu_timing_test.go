package cpu

import "testing"

type TimingTest struct {
    run (func(*CPU) int)
    expected int
}

func testTiming(t *testing.T, name string, run (func(*CPU) int), expected int) {
    var p = NewCPU()
    p.Reset()

    actual := run(p)

    if actual != expected {
        t.Errorf("%s failed", name)
        t.Errorf("Expected %d, got %d", expected, actual)
    }
}

func time(p *CPU, arguments []byte) int {
    for i := range arguments {
        p.Memory.buffer[i] = arguments[i]
    }

    p.PC = 0
    return p.Step()
}

func TestAddAndSubTiming(t *testing.T) {
    tests := map[string]TimingTest {
        "Adc immediate":
            TimingTest {
                func(p *CPU) int { return time(p, []byte{0x69}) },
                2,
            },
        "Adc zero page":
            TimingTest {
                func(p *CPU) int { return time(p, []byte{0x65, 0x01}) },
                3,
            },
        "Adc zero page X":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0x01
                    return time(p, []byte{0x75, 0x01, 0x00})
                },
                4,
            },
        "Adc absolute":
            TimingTest {
                func(p *CPU) int { return time(p, []byte{0x6d, 0x02, 0x00}) },
                4,
            },
        "Adc absolute X":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0x02
                    return time(p, []byte{0x7d, 0x00, 0x00})
                },
                4,
            },
        "Adc absolute X crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0xff
                    // 0xff + 0x0001 == 0x100 which crosses a page boundary
                    return time(p, []byte{0x7d, 0x01, 0x00})
                },
                5,
            },
        "Adc absolute Y":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0x02
                    return time(p, []byte{0x79, 0x00, 0x00})
                },
                4,
            },
        "Adc absolute Y crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0xff
                    // 0xff + 0x0001 == 0x100 which crosses a page boundary
                    return time(p, []byte{0x79, 0x01, 0x00})
                },
                5,
            },
        "Adc indexed indirect":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0x01
                    return time(p, []byte{0x61, 0x00, 0x03, 0x00})
                },
                6,
            },
        "Adc indirect indexed":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0x01
                    return time(p, []byte{0x71, 0x01, 0x02, 0x00})
                },
                5,
            },
        "Adc indirect indexed crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0xff
                    return time(p, []byte{0x71, 0x01, 0x02, 0x00})
                },
                6,
            },
        "Sbc immediate":
            TimingTest {
                func(p *CPU) int { return time(p, []byte{0xe9}) },
                2,
            },
        "Sbc zero page":
            TimingTest {
                func(p *CPU) int { return time(p, []byte{0xe5, 0x01}) },
                3,
            },
        "Sbc zero page X":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0x01
                    return time(p, []byte{0xf5, 0x01, 0x00})
                },
                4,
            },
        "Sbc absolute":
            TimingTest {
                func(p *CPU) int { return time(p, []byte{0xed, 0x02, 0x00}) },
                4,
            },
        "Sbc absolute X":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0x02
                    return time(p, []byte{0xfd, 0x00, 0x00})
                },
                4,
            },
        "Sbc absolute X crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0xff
                    // 0xff + 0x0001 == 0x100 which crosses a page boundary
                    return time(p, []byte{0xfd, 0x01, 0x00})
                },
                5,
            },
        "Sbc absolute Y":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0x02
                    return time(p, []byte{0xf9, 0x00, 0x00})
                },
                4,
            },
        "Sbc absolute Y crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0xff
                    // 0xff + 0x0001 == 0x100 which crosses a page boundary
                    return time(p, []byte{0xf9, 0x01, 0x00})
                },
                5,
            },
        "Sbc indexed indirect":
            TimingTest {
                func(p *CPU) int {
                    p.X = 0x01
                    return time(p, []byte{0xe1, 0x00, 0x03, 0x00})
                },
                6,
            },
        "Sbc indirect indexed":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0x01
                    return time(p, []byte{0xf1, 0x01, 0x02, 0x00})
                },
                5,
            },
        "Sbc indirect indexed crosses page":
            TimingTest {
                func(p *CPU) int {
                    p.Y = 0xff
                    return time(p, []byte{0xf1, 0x01, 0x02, 0x00})
                },
                6,
            },
    }

    for name, test := range tests {
        testTiming(t, name, test.run, test.expected)
    }
}
