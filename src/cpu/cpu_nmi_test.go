package cpu

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

func TestHandleNMISetPCToNMIVector(t *testing.T) {
    p := NewCPU()
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)
    p.Reset()

    p.Memory.Write(0xef, 0xfffa)
    p.Memory.Write(0xbe, 0xfffb)

    p.HandleNMI()

    assert.Equal(t, p.PC, Address(0xbeef))
}

func TestHandleNMIPushesCurrentPCToStack(t *testing.T) {
    p := NewCPU()
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)
    p.Reset()
    p.PC = 0xbeef

    p.HandleNMI()

    assert.Equal(t, p.Memory.Read(0x01fd), byte(0xbe))
    assert.Equal(t, p.Memory.Read(0x01fc), byte(0xef))
}

func TestHandleNMIPushesFlagsToStag(t *testing.T) {
    p := NewCPU()
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)
    p.Reset()
    p.Flags = 0xbe

    p.HandleNMI()

    assert.Equal(t, p.Memory.Read(0x01fb), byte(0xbe))
}

func TestSteppingWithNMISetShouldExecuteNMIAfterInstructions(t *testing.T) {
    p := NewCPU()
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)
    p.Reset()

    p.nmi = true
    p.Memory.Write(0xef, 0xfffa)
    p.Memory.Write(0xbe, 0xfffb)

    p.Memory.Write(0x69, 0x0000)
    p.PC = 0x0000

    p.Step()

    assert.Equal(t, p.PC, Address(0xbeef))
    assert.False(t, p.nmi)
}

func TestInterruptingNMISetsNMI(t *testing.T) {
    p := NewCPU()

    assert.False(t, p.nmi)
    p.Interrupt(NMI)
    assert.True(t, p.nmi)
}
