package chip8

import (
	"fmt"
	"os"
)

const (
	DisplayWidth  = 64
	DisplayHeight = 32
)

type Opcode uint16

type Cpu struct {
	Memory         []byte
	ProgramCounter uint16
	ValueRegister  [0x10]byte
	DelayTimer     byte
	SoundTimer     byte
	IndexRegister  uint16
	Display        []bool
}

func (cpu *Cpu) Reset() {
	cpu.Memory = make([]byte, 0x1000)
	cpu.Display = make([]bool, DisplayWidth*DisplayHeight)
	cpu.ProgramCounter = 0x200
	cpu.DelayTimer = 0
	cpu.SoundTimer = 0
	cpu.IndexRegister = 0
}

func (cpu *Cpu) Execute() {
	cpu.Decode(cpu.FetchInstruction())
}

func (cpu *Cpu) FetchInstruction() Opcode {
	high := cpu.Memory[cpu.ProgramCounter]
	low := cpu.Memory[cpu.ProgramCounter+1]

	return Opcode((uint16(high) << 8) + uint16(low))
}

func (cpu *Cpu) Decode(opcode Opcode) {
	f := byte((opcode & 0xF000) >> 12)
	x := byte((opcode & 0x0F00) >> 8)
	y := byte((opcode & 0x00F0) >> 4)
	n := byte(opcode & 0x000F)
	nn := byte(opcode & 0x00FF)
	nnn := uint16(opcode & 0x0FFF)
	switch f {
	case 0:
		if nn == 0xE0 {
			for i := 0; i < DisplayWidth*DisplayHeight; i++ {
				cpu.Display[i] = false
			}
			cpu.ProgramCounter += 2
		} else if nn == 0xEE {
			// 00EE Returns from a subroutine.
			// cpu.ProgramCounter +=
			fmt.Println("Place holder")
			os.Exit(-1)
		} else {
			os.Exit(-1)
		}
	case 1:
		// 1NNN	Jumps to address NNN.
		cpu.ProgramCounter = nnn
	case 2:
		// 2NNN Calls subroutine at NNN.
		// cpu.Stack ( pc + 2)
		cpu.ProgramCounter = nnn
	case 3:
		// 3XNN	Skips the next instruction if VX equals NN.
		if cpu.ValueRegister[x] == nn {
			cpu.ProgramCounter += 2
		}
		cpu.ProgramCounter += 2
	case 4:
		// 4XNN Skips the next instruction if VX doesn't equal NN.
		if cpu.ValueRegister[x] != nn {
			cpu.ProgramCounter += 2
		}
		cpu.ProgramCounter += 2
	case 5:
		// 5XY0 Skips the next instruction if VX equals VY.
		if cpu.ValueRegister[x] == cpu.ValueRegister[y] {
			cpu.ProgramCounter += 2
		}
		cpu.ProgramCounter += 2
	case 6:
		// 6XNN Sets VX to NN.
		cpu.ValueRegister[x] = nn
		cpu.ProgramCounter += 2
	case 7:
		// 7XNN Adds NN to VX.
		cpu.ValueRegister[x] += nn
		cpu.ProgramCounter += 2
	case 8:
		cpu.Decode8(opcode)
		cpu.ProgramCounter += 2
	case 9:
		// 9XY0 Skips the next instruction if VX doesn't equal VY.
		if cpu.ValueRegister[x] != cpu.ValueRegister[y] {
			cpu.ProgramCounter += 2
		}
		cpu.ProgramCounter += 2
	case 0xA:
		// ANNN Sets I to the address NNN.
		cpu.IndexRegister = nnn
		cpu.ProgramCounter += 2
	case 0xB:
		// BNNN Jumps to the address NNN plus V0.
		cpu.ProgramCounter = cpu.IndexRegister + uint16(cpu.ValueRegister[0])
	case 0xC:
		// CXNN Sets VX to the result of a bitwise and operation on a
		fmt.Println("Place holder")
		os.Exit(-1)
		cpu.ProgramCounter += 2
	case 0xD:
		cpu.Draw(x, y, n)
		cpu.ProgramCounter += 2
	case 0xE:
		// EX9E 	Skips the next instruction if the key stored in VX is pressed.
		// EXA1 	Skips the next instruction if the key stored in VX isn't pressed.

		if nn == 0x9E {
			fmt.Println("Place holder")
			os.Exit(-1)
		} else if nn == 0xA1 {
			fmt.Println("Place holder")
			os.Exit(-1)
			cpu.ProgramCounter += 2
		} else {
			fmt.Printf("Invalid instruction: %x\n", opcode)
			os.Exit(-1)
		}
		cpu.ProgramCounter += 2
	case 0xF:
		cpu.DecodeF(opcode)
		cpu.ProgramCounter += 2
	}
}

func (cpu *Cpu) Decode8(opcode Opcode) {
	x := byte((opcode & 0x0F00) >> 8)
	y := byte((opcode & 0x00F0) >> 4)
	n := byte(opcode & 0x000F)
	switch n {
	case 0x0:
		// 8XY0 Sets VX to the value of VY.
		cpu.ValueRegister[x] = cpu.ValueRegister[y]
	case 0x1:
		// 8XY1 Sets VX to VX bitwise or VY.
		cpu.ValueRegister[x] |= cpu.ValueRegister[y]
	case 0x2:
		// 8XY2 Sets VX to VX bitwise and VY.
		cpu.ValueRegister[x] &= cpu.ValueRegister[y]
	case 0x3:
		// 8XY3 Sets VX to VX bitwise xor VY.
		cpu.ValueRegister[x] ^= cpu.ValueRegister[y]
	case 0x4:
		// 8XY4 Adds VY to VX. VF is set to 1 when there's a carry,
		// and to 0 when there isn't.
		res := uint16(cpu.ValueRegister[x]) + uint16(cpu.ValueRegister[y])
		if res&0xFF00 > 0 {
			cpu.ValueRegister[0xF] = 1
		} else {
			cpu.ValueRegister[0xF] = 0
		}
		cpu.ValueRegister[x] = byte(0x00FF & res)
	case 0x5:
		// 8XY5 VY is subtracted from VX. VF is set to 0 when there's
		// a borrow, and 1 when there isn't.
		if cpu.ValueRegister[x] > cpu.ValueRegister[y] {
			cpu.ValueRegister[0xF] = 1
		} else {
			cpu.ValueRegister[0xF] = 0
		}
		cpu.ValueRegister[x] -= cpu.ValueRegister[y]
	case 0x6:
		// 8XY6 Shifts VX right by one. VF is set to the value of the
		// least significant bit of VX before the shift.
		cpu.ValueRegister[0xF] = 0x01 & cpu.ValueRegister[x]
		cpu.ValueRegister[x] >>= 1
	case 0x7:
		// 8XY7 Sets VX to VY minus VX. VF is set to 0 when there's a
		// borrow, and 1 when there isn't.
		if cpu.ValueRegister[y] > cpu.ValueRegister[x] {
			cpu.ValueRegister[0xF] = 1
		} else {
			cpu.ValueRegister[0xF] = 0
		}
		cpu.ValueRegister[x] = cpu.ValueRegister[y] - cpu.ValueRegister[x]
	case 0xE:
		// 8XYE Shifts VX left by one. VF is set to the value of the
		// most significant bit of VX before the shift.
		cpu.ValueRegister[0xF] = (0x80 & cpu.ValueRegister[x] >> 7)
		cpu.ValueRegister[x] <<= 1
	default:
		fmt.Printf("Invalid instruction: %x\n", opcode)
		os.Exit(-1)
	}
}

func (cpu *Cpu) DecodeF(opcode Opcode) {
	x := byte((opcode & 0x0F00) >> 8)
	nn := byte(opcode & 0x00FF)

	switch nn {
	case 0x07:
		// FX07	Sets VX to the value of the delay timer.
		cpu.ValueRegister[x] = cpu.DelayTimer
	case 0x0A:
		// FX0A	A key press is awaited, and then stored in VX.
		fmt.Println("PlaceHolder")
		os.Exit(-1)
		cpu.ProgramCounter -= 2
	case 0x15:
		// FX15	Sets the delay timer to VX.
		cpu.DelayTimer = cpu.ValueRegister[x]
	case 0x18:
		// FX18	Sets the sound timer to VX.
		cpu.SoundTimer = cpu.ValueRegister[x]
	case 0x1E:
		// FX1E Adds VX to I.
		cpu.IndexRegister += uint16(cpu.ValueRegister[x])
	case 0x29:
		// FX29 Sets I to the location of the sprite for the
		// character in VX. Characters 0-F (in hexadecimal) are
		// represented by a 4x5 font.
		fmt.Println("PlaceHolder")
	case 0x33:
		// FX33 Stores the binary-coded decimal representation of
		// VX, with the most significant of three digits at the
		// address in I, the middle digit at I plus 1, and the
		// least significant digit at I plus 2. (In other words,
		// take the decimal representation of VX, place the
		// hundreds digit in memory at location in I, the tens
		// digit at location I+1, and the ones digit at location
		// I+2.)
		cpu.Memory[cpu.IndexRegister] = cpu.ValueRegister[x] / 100
		cpu.Memory[cpu.IndexRegister+1] = cpu.ValueRegister[x] / 10
		cpu.Memory[cpu.IndexRegister+2] = cpu.ValueRegister[x] % 10
	case 0x55:
		// FX55 Stores V0 to VX (including VX) in memory starting at address I.
		for i := 0; i <= int(x); i++ {
			cpu.Memory[int(cpu.IndexRegister)+i] = cpu.ValueRegister[i]
		}
	case 0x65:
		// FX65	Fills V0 to VX (including VX) with values from memory starting at address I.
		for i := 0; i <= int(x); i++ {
			cpu.ValueRegister[i] = cpu.Memory[int(cpu.IndexRegister)+i]
		}
	default:
		fmt.Printf("Invalid instruction: %x\n", opcode)
		os.Exit(-1)
	}
}

func (cpu *Cpu) Draw(x byte, y byte, n byte) {
	// DXYN Draws a sprite at coordinate (VX, VY) that has a width of
	// 8 pixels and a height of N pixels. Each row of 8 pixels is read
	// as bit-coded starting from memory location I; I value doesn’t
	// change after the execution of this instruction. As described
	// above, VF is set to 1 if any screen pixels are flipped from set
	// to unset when the sprite is drawn, and to 0 if that doesn’t
	// happen

	cx := cpu.ValueRegister[x]
	cy := cpu.ValueRegister[y]

	pixel_deleted := false

	for yy := byte(0); yy < n; yy++ {
		b := cpu.Memory[cpu.IndexRegister+uint16(yy)]
		for xx := byte(0); xx < 8; xx++ {
			pixel_x := uint16(xx) + uint16(cx)
			pixel_y := uint16(yy) + uint16(cy)

			pos := DisplayWidth*pixel_y + pixel_x

			active := (b & (0x80 >> xx)) != 0

			if cpu.Display[pos] && active {
				pixel_deleted = true
			}

			cpu.Display[pos] = (active != cpu.Display[pos])
		}
	}

	if pixel_deleted {
		cpu.ValueRegister[0x0F] = 1
	} else {
		cpu.ValueRegister[0x0F] = 0
	}

	cpu.PrintDisplay()
}

func (cpu *Cpu) PrintDisplay() {
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			if cpu.Display[x+y*DisplayWidth] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}

func (cpu *Cpu) LoadFile(name string) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("Could not open file: ", name, err)
		os.Exit(-1)
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		fmt.Println("Could not open file: ", name, err)
		os.Exit(-1)
	}

	if stat.Size() > (0x800) {
		fmt.Println("File will not fit into memory")
		os.Exit(-1)
	}

	_, err = file.Read(cpu.Memory[0x200:])

	if err != nil {
		fmt.Println("Error reading file: ", err)
		os.Exit(-1)
	}
}
