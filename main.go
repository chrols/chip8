package main

import (
	"fmt"
	"os"

	"github.com/chrols/chip8/cpu"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No ROM file provided")
		os.Exit(-1)
	}

	cpu := chip8.Cpu{}

	cpu.Reset()
	cpu.LoadFile(os.Args[1])

	for i := 0; i < 100; i++ {
		fmt.Printf("%X\n", cpu.FetchInstruction())
		cpu.Execute()
	}
}
