package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Position struct {
	X int
	Y int
}

type Life struct {
	Pos       Position
	Neighbors uint8
	Alive     bool
	Locked    uint8
	Color     sdl.Color
	Rules     [][]int
}

func convertStringToLifeType(str string) [][]int {
	if str == "" {
		str = "23/3/2"
	}
	nums := strings.Split(str, "/")
	lifeType := make([][]int, 3)
	for i := 0; i < 3; i++ {
		numbers := strings.Split(nums[i], "")
		intArr := make([]int, len(numbers))
		for j := 0; j < len(numbers); j++ {
			intArr[j], _ = strconv.Atoi(numbers[j])
		}
		lifeType[i] = intArr
	}
	return lifeType
}

func main() {
	life := "gol"
	var lifeType [][]int
	if len(os.Args) > 3 {
		life = os.Args[3]
		switch life {
		default:
			lifeType = convertStringToLifeType(life)
		case "gol":
			// 23/3/2/
			lifeType = convertStringToLifeType("23/3/2")
		case "maze":
			lifeType = convertStringToLifeType("12345/3/2")
		case "repl":
			lifeType = convertStringToLifeType("1357/1357/2")
		case "wall":
			lifeType = convertStringToLifeType("2345/45678/2")
		case "34":
			lifeType = convertStringToLifeType("34/34/2")
		case "star":
			lifeType = convertStringToLifeType("3456/278/8")
		}
	}
	lifeMap := makeLifeMap()
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("gol", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1000, 1000, sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	w, err := strconv.Atoi(os.Args[1])
	h, err := strconv.Atoi(os.Args[2])
	width, height := window.GetSize()
	height /= int32(h)
	if width < height {
		height = width
	}
	width = height

	// fmt.Println("width ", width, "height", height)
	window.UpdateSurface()
	// randomizeMap(lifeMap)
	// lifeMap[1][1] = true
	// lifeMap[0][1], lifeMap[1][2], lifeMap[2][0], lifeMap[2][1], lifeMap[2][2] = true, true, true, true, true
	lifeMap = updateNeighbors(lifeMap)
	for i := range lifeMap {
		for j := range lifeMap[i] {
			renderCell(width, height, j, i, lifeMap[i][j], surface)
		}
	}
	// printNeighborMap(neighborMap)
	running := true
	avg := time.Duration(0)
	paused := false
	dragging := false
	draggingState := true
	// fmt.Print("\033[H\033[2J")
	for running {
		// fmt.Println(height)
		start := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("exit")
				running = false
			case *sdl.MouseMotionEvent:
				if dragging {
					x := t.X / width
					y := t.Y / height
					if int(y) >= len(lifeMap) || int(x) >= len(lifeMap[y]) {
						break
					}
					lifeMap[y][x].Alive = draggingState
					life := Life{}
					life.Pos.X, life.Pos.Y, life.Alive = int(x), int(y), lifeMap[y][x].Alive
					changeNeighborOfCells(life, lifeMap)
					lifeMap = updateNeighbors(lifeMap)
					renderCell(width, height, int(x), int(y), lifeMap[y][x], surface)
					go func() {
						window.UpdateSurface()
					}()
				}
			case *sdl.MouseButtonEvent:
				if t.State != sdl.PRESSED {
					dragging = false
					break
				}
				dragging = true
				x, y := int(t.X/width), int(t.Y/height)
				if y >= len(lifeMap) || x >= len(lifeMap[y]) {
					break
				}
				lifeMap[y][x].Alive = !lifeMap[y][x].Alive
				draggingState = lifeMap[y][x].Alive
				p := Position{x, y}
				life := Life{}
				life.Pos, life.Alive = p, draggingState
				changeNeighborOfCells(life, lifeMap)
				renderLifeMap(lifeMap, width, height, surface)
				window.UpdateSurface()
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYUP {
					break
				}
				switch t.Keysym.Sym {
				case sdl.K_r:
					clearWindow(surface, window)
					randomizeMap(lifeMap)
					lifeMap = updateNeighbors(lifeMap)
					renderLifeMap(lifeMap, width, height, surface)
					go func() {
						window.UpdateSurface()
					}()
				case sdl.K_SPACE:
					paused = !paused
				case sdl.K_ESCAPE, sdl.K_q:
					fmt.Println("Exit")
					running = false
				case sdl.K_a:
					clearWindow(surface, window)
					switch t.Keysym.Mod {
					default:
						w, width, height = expandWindowRight(window, w, height)
						lifeMap = expandMapsRight(lifeMap)
						for i := 0; i < len(lifeMap); i++ {
							rotateRight(lifeMap[i], 1)
						}
					case sdl.KMOD_LSHIFT:
						w--
						width, height = updateWidthAndHeight(w, h, window)
						for i := 0; i < len(lifeMap); i++ {
							lifeMap[i][0].Alive, lifeMap[i][0].Neighbors = false, 0
							lifeMap[i] = lifeMap[i][1:]
						}
					}
					renderLifeMap(lifeMap, width, height, surface)
					go func() {
						window.UpdateSurface()
					}()
				case sdl.K_w:
					clearWindow(surface, window)
					switch t.Keysym.Mod {
					default:
						h, width, height = expandWindowDown(window, h, width)
						lifeMap = expandMapsDown(lifeMap)
						rotateRight(lifeMap, 1)
					case sdl.KMOD_LSHIFT:
						h--
						width, height = updateWidthAndHeight(w, h, window)
						lifeMap = lifeMap[1:]
					}
					renderLifeMap(lifeMap, width, height, surface)
					go func() {
						window.UpdateSurface()
					}()
				case sdl.K_s:
					clearWindow(surface, window)
					switch t.Keysym.Mod {
					default:
						h, width, height = expandWindowDown(window, h, width)
						lifeMap = expandMapsDown(lifeMap)
					case sdl.KMOD_LSHIFT:
						h--
						width, height = updateWidthAndHeight(w, h, window)
						lifeMap = lifeMap[:len(lifeMap)-1]
					}
					renderLifeMap(lifeMap, width, height, surface)
					go func() {
						window.UpdateSurface()
					}()
				case sdl.K_f:
					surface, err := window.GetSurface()
					if err != nil {
						panic(err)
					}
					passFrame(lifeMap, width, height, surface, lifeType)
					err = window.UpdateSurface()
					if err != nil {
						panic(err)
					}
				case sdl.K_d:
					clearWindow(surface, window)
					switch t.Keysym.Mod {
					default:
						w, width, height = expandWindowRight(window, w, height)
						lifeMap = expandMapsRight(lifeMap)
					case sdl.KMOD_LSHIFT:
						w--
						width, height = updateWidthAndHeight(w, h, window)
						for i := 0; i < len(lifeMap); i++ {
							lifeMap[i][len(lifeMap[i])-1].Alive, lifeMap[i][len(lifeMap[i])-1].Neighbors, lifeMap[i][len(lifeMap[i])-1].Locked = false, 0, 0
							lifeMap[i] = lifeMap[i][:len(lifeMap[i])-1]
						}
					}
					renderLifeMap(lifeMap, width, height, surface)
					go func() {
						window.UpdateSurface()
					}()
				case sdl.K_c:
					h, w := len(lifeMap), len(lifeMap[0])
					lifeMap = make([][]Life, h)
					for i := 0; i < len(lifeMap); i++ {
						lifeMap[i] = make([]Life, w)
					}
					lifeMap = updateNeighbors(lifeMap)
					renderLifeMap(lifeMap, width, height, surface)
				}
			}
		}
		if !paused {
			surface, err := window.GetSurface()
			if err != nil {
				panic(err)
			}
			passFrame(lifeMap, width, height, surface, lifeType)
			err = window.UpdateSurface()
			avg = handleFrameTime(start, avg)
			if err != nil {
				panic(err)
			}
		}
		err := window.UpdateSurface()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Frame avg time ", avg)
}

func updateWidthAndHeight(w, h int, window *sdl.Window) (int32, int32) {
	width, height := window.GetSize()
	width /= int32(w)
	height /= int32(h)
	if width > height {
		width = height
	}
	height = width
	return width, height
}

func rotateRight[T any](nums []T, k int) {
	k %= len(nums)
	new_array := make([]T, len(nums))
	copy(new_array[:k], nums[len(nums)-k:])
	copy(new_array[k:], nums[:len(nums)-k])
	copy(nums, new_array)
}

func expandMapsRight(lifeMap [][]Life) [][]Life {
	for i := range lifeMap {
		lifeMap[i] = append(lifeMap[i], Life{})
	}
	for i := range lifeMap {
		for j := range lifeMap[i] {
			lifeMap[i][j].Neighbors = uint8(getNeighbors(j, i, lifeMap))
		}
	}
	return lifeMap
}

func expandWindowRight(window *sdl.Window, w int, height int32) (int, int32, int32) {
	width, _ := window.GetSize()
	w++
	width /= int32(w)
	if width < height {
		height = width
	}
	width = height
	return w, width, height
}

func clearWindow(surface *sdl.Surface, window *sdl.Window) {
	width, height := window.GetSize()
	renderCell(width, height, 0, 0, Life{}, surface)
}

func renderLifeMap(lifeMap [][]Life, width, height int32, surface *sdl.Surface) {
	for i := range lifeMap {
		for j := range lifeMap[i] {
			renderCell(width, height, j, i, lifeMap[i][j], surface)
		}
	}
}

func expandMapsDown(lifeMap [][]Life) [][]Life {
	lifeMap = append(lifeMap, make([]Life, len(lifeMap[0])))
	for i := range lifeMap {
		for j := range lifeMap[i] {
			lifeMap[i][j].Neighbors = uint8(getNeighbors(j, i, lifeMap))
		}
	}
	return lifeMap
}

func expandWindowDown(window *sdl.Window, h int, width int32) (int, int32, int32) {
	_, height := window.GetSize()
	h++
	height /= int32(h)
	if width < height {
		height = width
	}
	width = height
	return h, width, height
}

func handleFrameTime(start time.Time, avg time.Duration) time.Duration {
	frameTime := time.Since(start)
	avg = (avg + frameTime) / 2
	if time.Since(start) < 32*time.Millisecond {
		time.Sleep(32*time.Millisecond - time.Since(start))
	} else if frameTime > 100*time.Millisecond {
		go fmt.Println("Long frame comp", frameTime)
	}
	return avg
}

func updateNeighbors(lifeMap [][]Life) [][]Life {
	for i := 0; i < len(lifeMap); i++ {
		for j := 0; j < len(lifeMap[i]); j++ {
			lifeMap[i][j].Neighbors = getNeighbors(j, i, lifeMap)
		}
	}
	return lifeMap
}

func renderCell(width, height int32, x, y int, life Life, surface *sdl.Surface) {
	rect := sdl.Rect{int32(x) * width, int32(y) * height, width, height}
	colour := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	if width == height {
		colour = sdl.Color{R: 25, G: 25, B: 25, A: 255}
	}
	if life.Alive {
		colour = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	}
	if life.Locked > 0 {
		colour = sdl.Color{R: (uint8(255.0 / ((1.0/float64(life.Locked))+1))), G: (uint8(255.0 / ((1.0/float64(life.Locked))+1))), B: 255, A: 255}
	}
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	// err := surface.Lock()
	// if err != nil {
	// panic(err)
	// }
	// defer surface.Unlock()
	err := surface.FillRect(&rect, pixel)
	if err != nil {
		panic(err)
	}
}

func randomizeMap(arr [][]Life) {
	for i := range arr {
		for j := range arr[i] {
			val := rand.Intn(2)
			if val == 1 {
				arr[i][j].Alive = true
			} else {
				arr[i][j].Alive = false
			}
		}
	}
}

func printMap(arr [][]Life) {
	// fmt.Print("\033[s")
	str := ""
	for i := range arr {
		for j := range arr[i] {
			if !arr[i][j].Alive {
				str += ""
			} else {
				str += "󰝤"
			}
		}
		str += "\n"
	}
	// fmt.Print("\033[u\033[K")
	fmt.Printf("%s", str)
}

func makeLifeMap() [][]Life {
	if len(os.Args) < 3 {
		fmt.Println("Not enough args, usage: ./main height width")
		os.Exit(1)
	}
	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid argument 1(height)")
		os.Exit(1)
	}
	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid argument 2(width)")
		os.Exit(1)
	}
	var lifeMap = make([][]Life, height)

	for i := 0; i < len(lifeMap); i++ {
		lifeMap[i] = make([]Life, width)
	}
	for y := 0; y < len(lifeMap); y++ {
		for x := 0; x < len(lifeMap[y]); x++ {
			lifeMap[y][x].Pos = Position{x, y}
		}
	}
	return lifeMap
}

func passFrame(lifeMap [][]Life, width, height int32, surface *sdl.Surface, lifeType [][]int) {
	var wg sync.WaitGroup
	ch := make(chan Life, (len(lifeMap) * len(lifeMap[0])))
	for y := 0; y < len(lifeMap); y++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processRow(y, lifeMap, ch, lifeType)
		}()
	}
	wg.Wait()
	close(ch)
	for i := range ch {
    if i.Locked == 100 {
      lifeMap[i.Pos.Y][i.Pos.X].Locked = 0
      lifeMap[i.Pos.Y][i.Pos.X].Alive = false
      renderCell(width, height, i.Pos.X, i.Pos.Y, lifeMap[i.Pos.Y][i.Pos.X], surface)
      continue
    }
    if i.Locked > 0 && i.Alive {
      lifeMap[i.Pos.Y][i.Pos.X].Locked--
      renderCell(width, height, i.Pos.X, i.Pos.Y, i, surface)
      continue
    }
		lifeMap[i.Pos.Y][i.Pos.X].Alive = i.Alive
		lifeMap[i.Pos.Y][i.Pos.X].Locked = i.Locked
		renderCell(width, height, i.Pos.X, i.Pos.Y, i, surface)
		changeNeighborOfCells(i, lifeMap)
	}
  // printMap(lifeMap)
  // printNeighborMap(lifeMap)
}

func printNeighborMap(neighborMap [][]Life) {
	for i := 0; i < len(neighborMap); i++ {
		for j := 0; j < len(neighborMap[i]); j++ {
			fmt.Print(neighborMap[i][j].Neighbors, " ")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func processRow(y int, lifeMap [][]Life, ch chan Life, lifeType [][]int) {
	for x := 0; x < len(lifeMap[y]); x++ {
		life := processCell(x, y, lifeMap, lifeType)
		if life.Pos.X != -1 {
			ch <- life
		}
	}
}

func changeNeighborOfCells(l Life, lifeMap [][]Life) {
	x, y := l.Pos.X, l.Pos.Y
	rows := len(lifeMap)
	cols := len(lifeMap[0])
	if l.Alive {
		lifeMap[(y-1+rows)%rows][(x-1+cols)%cols].Neighbors++
		lifeMap[(y-1+rows)%rows][x].Neighbors++
		lifeMap[(y-1+rows)%rows][(x+1)%cols].Neighbors++
		lifeMap[y][(x-1+cols)%cols].Neighbors++
		lifeMap[y][(x+1)%cols].Neighbors++
		lifeMap[(y+1)%rows][(x-1+cols)%cols].Neighbors++
		lifeMap[(y+1)%rows][x].Neighbors++
		lifeMap[(y+1)%rows][(x+1)%cols].Neighbors++
	} else {
		lifeMap[(y-1+rows)%rows][(x-1+cols)%cols].Neighbors--
		lifeMap[(y-1+rows)%rows][x].Neighbors--
		lifeMap[(y-1+rows)%rows][(x+1)%cols].Neighbors--
		lifeMap[y][(x-1+cols)%cols].Neighbors--
		lifeMap[y][(x+1)%cols].Neighbors--
		lifeMap[(y+1)%rows][(x-1+cols)%cols].Neighbors--
		lifeMap[(y+1)%rows][x].Neighbors--
		lifeMap[(y+1)%rows][(x+1)%cols].Neighbors--
	}
}

func processCell(x, y int, lifeMap [][]Life, lifeType [][]int) Life {
	alive := lifeMap[y][x].Alive
	stay := Life{}
	stay.Pos.X = -1
	change := lifeMap[y][x]
	change.Pos.X = x
	change.Pos.Y = y
	change.Alive = !change.Alive
	neighbors := lifeMap[y][x].Neighbors
	if lifeMap[y][x].Locked > 0 {
    change.Locked--
    if change.Locked == 0 {
      change.Locked = 100
    }
    // fmt.Println("hi")
    return change
  }
  if !alive {
		for i := 0; i < len(lifeType[1]); i++ {
			if neighbors == uint8(lifeType[1][i]) {
				return change
			}
		}
	}
	if alive {
		for i := 0; i < len(lifeType[0]); i++ {
			if neighbors == uint8(lifeType[0][i]) {
				return stay
			}
		}
		change.Locked = uint8(lifeType[2][0] - 2)
		return change
	}
	return stay
}

func getNeighbors(x, y int, arr [][]Life) uint8 {
	return uint8(checkLeft(x, y, arr) + checkRight(x, y, arr) + checkDown(x, y, arr) + checkUp(x, y, arr) + checkBottomCorners(x, y, arr) + checkTopCorners(x, y, arr))
}

func checkRight(x, y int, arr [][]Life) int {
	if x == len(arr[y])-1 {
		if !arr[y][0].Alive {
			return 0
		}
		return 1
	}
	if !arr[y][x+1].Alive {
		return 0
	}
	return 1
}

func checkLeft(x, y int, arr [][]Life) int {
	if x == 0 {
		if !arr[y][len(arr[y])-1].Alive {
			return 0
		}
		return 1
	}
	if !arr[y][x-1].Alive {
		return 0
	}
	return 1
}

func checkUp(x, y int, arr [][]Life) int {
	if y == 0 {
		if !arr[len(arr)-1][x].Alive {
			return 0
		}
		return 1
	}
	if !arr[y-1][x].Alive {
		return 0
	}
	return 1
}

func checkDown(x, y int, arr [][]Life) int {
	if y == len(arr)-1 {
		if !arr[0][x].Alive {
			return 0
		}
		return 1
	}
	if !arr[y+1][x].Alive {
		return 0
	}
	return 1
}

func checkBottomCorners(x, y int, arr [][]Life) int {
	counter := 0
	if y == len(arr)-1 {
		counter += checkRight(x, 0, arr)
		counter += checkLeft(x, 0, arr)
		return counter
	}
	counter += checkRight(x, y+1, arr)
	counter += checkLeft(x, y+1, arr)
	return counter
}

func checkTopCorners(x, y int, arr [][]Life) int {
	counter := 0
	if y == 0 {
		counter += checkRight(x, len(arr)-1, arr)
		counter += checkLeft(x, len(arr)-1, arr)
		return counter
	}
	counter += checkRight(x, y-1, arr)
	counter += checkLeft(x, y-1, arr)
	return counter
}
