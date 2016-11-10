package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/chrols/chip8/cpu"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	Factor       = 20
	WindowTitle  = "CHIP-8 Emulator"
	WindowWidth  = 64 * Factor
	WindowHeight = 32 * Factor
	FrameRate    = 60

	RectWidth  = WindowWidth / 64
	RectHeight = WindowHeight / 32
	NumRects   = 64 * 32

	SystemClock = 1000
)

var rects [NumRects]sdl.Rect

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("No ROM file provided")
		os.Exit(-1)
	}

	cpu := chip8.Cpu{}

	cpu.Reset()
	cpu.LoadFile(os.Args[1])

	delay_ticker := time.NewTicker(1000 / 60 * time.Millisecond)
	system_ticker := time.NewTicker(time.Second / SystemClock)

	go cpu.DelayTick(delay_ticker.C)
	go cpu.CycleTick(system_ticker.C)

	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	sdl.Do(func() {
		window, err = sdl.CreateWindow(WindowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_OPENGL)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		sdl.Do(func() {
			window.Destroy()
		})
	}()

	sdl.Do(func() {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	})
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(2)
	}
	defer func() {
		sdl.Do(func() {
			renderer.Destroy()
		})
	}()

	sdl.Do(func() {
		renderer.Clear()
	})

	for i := range rects {
		rects[i] = sdl.Rect{
			X: int32(i % 64 * RectWidth),
			Y: int32(i / 64 * RectHeight),
			W: RectWidth,
			H: RectHeight,
		}
	}
	sdl.Do(func() {
		renderer.Clear()
		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.FillRect(&sdl.Rect{0, 0, WindowWidth, WindowHeight})
	})

	running := true
	for running {

		sdl.Do(func() {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					running = false
				case *sdl.KeyDownEvent:
					switch t.Keysym.Sym {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_F5:
						cpu.Reset()
						cpu.LoadFile(os.Args[1])
					case sdl.K_p:
						cpu.PrintDisplay()
					default:
						cpu.KeyPressed[Convert(t.Keysym.Scancode)] = true
					}
				case *sdl.KeyUpEvent:
					cpu.KeyPressed[Convert(t.Keysym.Scancode)] = false
				}
			}
		})

		sdl.Do(func() {
			for i := 0; i < len(rects); i++ {
				if cpu.Display[i] {
					renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
				} else {
					renderer.SetDrawColor(0, 0, 0, 0)
				}
				renderer.FillRect(&rects[i])
			}
		})

		sdl.Do(func() {
			renderer.Present()
			sdl.Delay(1000 / FrameRate)
		})
	}

	os.Exit(0)
}
func Convert(scancode sdl.Scancode) byte {
	switch scancode {
	case sdl.SCANCODE_X:
		return (0x00)
	case sdl.SCANCODE_1:
		return (0x01)
	case sdl.SCANCODE_2:
		return (0x02)
	case sdl.SCANCODE_3:
		return (0x03)
	case sdl.SCANCODE_Q:
		return (0x04)
	case sdl.SCANCODE_W:
		return (0x05)
	case sdl.SCANCODE_E:
		return (0x06)
	case sdl.SCANCODE_A:
		return (0x07)
	case sdl.SCANCODE_S:
		return (0x08)
	case sdl.SCANCODE_D:
		return (0x09)
	case sdl.SCANCODE_Z:
		return (0x0A)
	case sdl.SCANCODE_C:
		return (0x0B)
	case sdl.SCANCODE_4:
		return (0x0C)
	case sdl.SCANCODE_R:
		return (0x0D)
	case sdl.SCANCODE_F:
		return (0x0E)
	case sdl.SCANCODE_V:
		return (0x0F)
	default:
		// FIXME Handle
		return 0x00
	}
}
