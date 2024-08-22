package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type Position struct {
	X   int
	Y   int
	Val bool
}

func main() {
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
	width /= int32(w)
	height /= int32(h)
	if width < height {
		height = width
	}
	width = height

	// fmt.Println("width ", width, "height", height)
	window.UpdateSurface()
	randomizeMap(lifeMap)
	// lifeMap[1][1] = true
	// lifeMap[0][1], lifeMap[1][2], lifeMap[2][0], lifeMap[2][1], lifeMap[2][2] = true, true, true, true, true
	neighborMap := generateNeighborMap(lifeMap)
	for i := range lifeMap {
		for j := range lifeMap[i] {
			renderCell(width, height, j, i, lifeMap[i][j], surface)
		}
	}
	// printNeighborMap(neighborMap)
	running := true
	avg := time.Duration(0)
	// fmt.Print("\033[H\033[2J")
	ch := make(chan [][]bool, 60)
	for running {
		start := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("exit", t)
				running = false
			
      case *sdl.KeyboardEvent:
        if t.Type == sdl.KEYUP {
          break
        }
        switch t.Keysym.Sym{
        case sdl.K_w:
          // implement logic for adding to the "top" of the array
        case sdl.K_s:
          clearWindow(surface, window)
          h, width, height = expandWindowDown(window, h)
          lifeMap, neighborMap = expandMapsDown(lifeMap, neighborMap)
          renderLifeMap(lifeMap, width, height, surface)
        
        case sdl.K_d:
          clearWindow(surface, window)
          w, width, height = expandWindowRight(window, w)
          lifeMap, neighborMap = expandMapsRight(lifeMap, neighborMap)
        }
      }
		}
  
		surface, err := window.GetSurface()
		if err != nil {
			panic(err)
		}
		passFrame(lifeMap, neighborMap, width, height, surface, ch)
    // fmt.Println(height)
		err = window.UpdateSurface()
    avg = handleFrameTime(start, avg)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Frame avg time ", avg)
}

func expandMapsRight(lifeMap [][]bool, neighborMap[][]int16) ([][]bool, [][]int16) {
  for i := range neighborMap {
    neighborMap[i] = append(neighborMap[i], 0)
    lifeMap[i] = append(lifeMap[i], false)
  }
	for i := range neighborMap {
	  for j := range neighborMap[i] {
			neighborMap[i][j] = int16(getNeighbors(j, i, lifeMap))
		}
	}
  return lifeMap, neighborMap
}

func expandWindowRight(window *sdl.Window, w int) (int, int32, int32) {
  width, height := window.GetSize()
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
  renderCell(width, height, 0, 0, false, surface)
}

func renderLifeMap(lifeMap [][]bool, width, height int32, surface *sdl.Surface) {
  for i := range lifeMap {
		for j := range lifeMap[i] {
			renderCell(width, height, j, i, lifeMap[i][j], surface)
		}
	}
}

func expandMapsDown(lifeMap [][]bool, neighborMap[][]int16) ([][]bool, [][]int16) {
  lifeMap = append(lifeMap, make([]bool, len(lifeMap[0])))
  neighborMap = append(neighborMap, make([]int16, len(neighborMap[0])))
	for i := range neighborMap {
	  for j := range neighborMap[i] {
			neighborMap[i][j] = int16(getNeighbors(j, i, lifeMap))
		}
	}
  return lifeMap, neighborMap
}

func expandWindowDown(window *sdl.Window, h int) (int, int32, int32) {
  width, height := window.GetSize()
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
	} else if frameTime > 32*time.Millisecond {
		go fmt.Println("Long frame comp", frameTime)
	}
  return avg
}

func generateNeighborMap(lifeMap [][]bool) [][]int16 {
	neighborMap := make([][]int16, len(lifeMap))
	for i := 0; i < len(lifeMap); i++ {
		neighborMap[i] = make([]int16, len(lifeMap[i]))
	}
	for i := 0; i < len(lifeMap); i++ {
		for j := 0; j < len(lifeMap[i]); j++ {
			neighborMap[i][j] = int16(getNeighbors(j, i, lifeMap))
		}
	}
	return neighborMap
}

func renderCell(width, height int32, x, y int, alive bool, surface *sdl.Surface) {
	rect := sdl.Rect{int32(x) * width, int32(y) * height, width, height}
	colour := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	if alive {
		colour = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	}
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	err := surface.FillRect(&rect, pixel)
	if err != nil {
		panic(err)
	}
}

func randomizeMap(arr [][]bool) {
	for i := range arr {
		for j := range arr[i] {
			val := rand.Intn(2)
			if val == 1 {
				arr[i][j] = true
			} else {
				arr[i][j] = false
			}
		}
	}
}

func printMap(arr [][]bool) {
	// fmt.Print("\033[s")
	str := ""
	for i := range arr {
		for j := range arr[i] {
			if !arr[i][j] {
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

func makeLifeMap() [][]bool {
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
	var lifeMap = make([][]bool, height)

	for i := 0; i < len(lifeMap); i++ {
		lifeMap[i] = make([]bool, width)
	}
	return lifeMap
}

func passFrame(lifeMap [][]bool, neighborMap [][]int16, width, height int32, surface *sdl.Surface, frameBuffer chan [][]bool) {
	// printNeighborMap(neighborMap)
	// printMap(lifeMap)
	// fmt.Println(lifeMap)
	var wg sync.WaitGroup
	ch := make(chan Position, (len(lifeMap) * len(lifeMap[1])))
	for y := 0; y < len(lifeMap); y++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processRow(y, lifeMap, neighborMap, ch)
		}()
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		lifeMap[i.Y][i.X] = i.Val
		renderCell(width, height, i.X, i.Y, i.Val, surface)
		changeNeighborOfCells(i, neighborMap)
	}
}

func printNeighborMap(neighborMap [][]int16) {
	for i := 0; i < len(neighborMap); i++ {
		for j := 0; j < len(neighborMap[i]); j++ {
			fmt.Print(neighborMap[i][j], " ")
		}
		fmt.Print("\n")
	}
}

func processRow(y int, lifeMap [][]bool, neighborMap [][]int16, ch chan Position) {
	for x := 0; x < len(lifeMap[y]); x++ {
		pos := processCell(x, y, lifeMap, neighborMap)
		if pos.X != -1 {
			ch <- pos
		}
	}
}

func changeNeighborOfCells(p Position, neighborMap [][]int16) {
	x, y := p.X, p.Y
	rows := len(neighborMap)
	cols := len(neighborMap[0])
	if p.Val {
		neighborMap[(y-1+rows)%rows][(x-1+cols)%cols]++
		neighborMap[(y-1+rows)%rows][x]++
		neighborMap[(y-1+rows)%rows][(x+1)%cols]++
		neighborMap[y][(x-1+cols)%cols]++
		neighborMap[y][(x+1)%cols]++
		neighborMap[(y+1)%rows][(x-1+cols)%cols]++
		neighborMap[(y+1)%rows][x]++
		neighborMap[(y+1)%rows][(x+1)%cols]++
	} else {
		neighborMap[(y-1+rows)%rows][(x-1+cols)%cols]--
		neighborMap[(y-1+rows)%rows][x]--
		neighborMap[(y-1+rows)%rows][(x+1)%cols]--
		neighborMap[y][(x-1+cols)%cols]--
		neighborMap[y][(x+1)%cols]--
		neighborMap[(y+1)%rows][(x-1+cols)%cols]--
		neighborMap[(y+1)%rows][x]--
		neighborMap[(y+1)%rows][(x+1)%cols]--
	}
}

func processCell(x, y int, lifeMap [][]bool, neighborMap [][]int16) Position {
	neighbors := neighborMap[y][x]
	if neighbors == 3 && !lifeMap[y][x] {
		// fmt.Println("Turning ",x,y,"alive")

		return Position{x, y, true}
	}
	if (neighbors < 2 || neighbors > 3) && lifeMap[y][x] {
		// fmt.Println("Killing ",x,y)
		return Position{x, y, false}
	}
	return Position{-1, 0, false}
}

func getNeighbors(x, y int, arr [][]bool) int {
	return checkLeft(x, y, arr) + checkRight(x, y, arr) + checkDown(x, y, arr) + checkUp(x, y, arr) + checkBottomCorners(x, y, arr) + checkTopCorners(x, y, arr)
}

func checkRight(x, y int, arr [][]bool) int {
	if x == len(arr[y])-1 {
		if !arr[y][0] {
			return 0
		}
		return 1
	}
	if !arr[y][x+1] {
		return 0
	}
	return 1
}

func checkLeft(x, y int, arr [][]bool) int {
	if x == 0 {
		if !arr[y][len(arr[y])-1] {
			return 0
		}
		return 1
	}
	if !arr[y][x-1] {
		return 0
	}
	return 1
}

func checkUp(x, y int, arr [][]bool) int {
	if y == 0 {
		if !arr[len(arr)-1][x] {
			return 0
		}
		return 1
	}
	if !arr[y-1][x] {
		return 0
	}
	return 1
}

func checkDown(x, y int, arr [][]bool) int {
	if y == len(arr)-1 {
		if !arr[0][x] {
			return 0
		}
		return 1
	}
	if !arr[y+1][x] {
		return 0
	}
	return 1
}

func checkBottomCorners(x, y int, arr [][]bool) int {
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

func checkTopCorners(x, y int, arr [][]bool) int {
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
