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
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("gol", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 500, 500, 0x00000020)
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
	width, height := int32(500), int32(500)
	width /= int32(w)
	height /= int32(h)

	// fmt.Println("width ", width, "height", height)
	window.UpdateSurface()
	lifeMap := makeLifeMap()
	randomizeMap(lifeMap)
	// lifeMap[0][1], lifeMap[1][2], lifeMap[2][0], lifeMap[2][1], lifeMap[2][2] = true, true, true, true, true
	start := time.Now()
	for i := range lifeMap {
		for j := range lifeMap[i] {
			renderCell(width, height, j, i, lifeMap[i][j], surface, window)
		}
	}
	// fmt.Print("\033[H\033[2J")
	for {
		// printMap(lifeMap)
		// fmt.Println()
		passFrame(lifeMap, width, height, surface, window)
		time.Sleep(16 * time.Millisecond)
	}
	_ = time.Since(start)
}

func renderCell(width, height int32, x, y int, alive bool, surface *sdl.Surface, window *sdl.Window) {
	rect := sdl.Rect{int32(x) * width, int32(y) * height, width, height}
	colour := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	if alive {
		colour = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	}
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)
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
	fmt.Print("\033[s")
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
	fmt.Print("\033[u\033[K")
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

func passFrame(lifeMap [][]bool, width, height int32, surface *sdl.Surface, window *sdl.Window) {
	var wg sync.WaitGroup
	ch := make(chan Position, (len(lifeMap) * len(lifeMap[1])))
	for y := 0; y < len(lifeMap); y++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processRow(y, lifeMap, ch)
		}()
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		lifeMap[i.Y][i.X] = i.Val
		renderCell(width, height, i.X, i.Y, i.Val, surface, window)
	}
	window.UpdateSurface()
}

func processRow(y int, lifeMap [][]bool, ch chan Position) {
	for x := 0; x < len(lifeMap[y]); x++ {
		pos := processCell(x, y, lifeMap)
		if pos.X != -1 {
			ch <- pos
		}
	}
}

func processCell(x, y int, lifeMap [][]bool) Position {
	neighbors := getNeighbors(x, y, lifeMap)
	if neighbors == 3 && !lifeMap[y][x] {
		// fmt.Println("Turning ",x,y,"alive")
		return Position{x, y, true}
	}
	if neighbors < 2 || neighbors > 3 && lifeMap[y][x] {
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
