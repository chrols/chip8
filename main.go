package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
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
)

var rects [NumRects]sdl.Rect
var runningMutex sync.Mutex

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("No ROM file provided")
		os.Exit(-1)
	}

	cpu := chip8.Cpu{}

	cpu.Reset()
	cpu.LoadFile(os.Args[1])

	timer := time.NewTicker(1000 / 60 * time.Millisecond)

	go cpu.Update(timer.C)

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

	running := true
	for running {
		cpu.Execute()

		sdl.Do(func() {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					runningMutex.Lock()
					running = false
					runningMutex.Unlock()
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

			renderer.Clear()
			renderer.SetDrawColor(0, 0, 0, 0)
			renderer.FillRect(&sdl.Rect{0, 0, WindowWidth, WindowHeight})
		})

		// Do expensive stuff using goroutines
		wg := sync.WaitGroup{}
		//for i := range rects {
		for i := 0; i < len(rects); i++ {
			wg.Add(1)
			go func(i int) {
				//rects[i].X = (rects[i].X + 10) % WindowWidth
				sdl.Do(func() {
					if cpu.Display[i] {
						renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
					} else {
						renderer.SetDrawColor(0, 0, 0, 0)
					}
					//renderer.DrawRect(&rects[i])
					renderer.FillRect(&rects[i])
				})
				wg.Done()
			}(i)
		}
		wg.Wait()

		sdl.Do(func() {
			renderer.Present()
			//sdl.Delay(1000 / FrameRate)
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
