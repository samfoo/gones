package main

import (
    "os"
    "cpu"
    "log"
)

func main() {
    var file *os.File
    var err error
    if file, err = os.Open("assets/nestest"); err != nil {
        log.Fatal(err)
        return
    }

    proc := new(cpu.CPU)
    proc.Reset()

    programSpace := proc.Memory[0xc000:0xffff]
    if _, err := file.Read(programSpace); err != nil {
        log.Fatal(err)
        return
    }

    proc.PC = 0xc000

    for i:=0; i < 40; i++ {
        proc.Step()
    }
}
