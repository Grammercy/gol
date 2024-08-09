package main

import (
	"fmt"
	"time"
	// "sync"
	"os"
	"strconv"
	// "time"
)

type Position struct {
  X int
  Y int
  Val bool
}

func main(){
  lifeMap := makeLifeMap()
  lifeMap[0][1], lifeMap[1][2], lifeMap[2][0], lifeMap[2][1], lifeMap[2][2] = true, true, true, true, true
  
  for {
    printMap(lifeMap)
    fmt.Println()
    passFrame(lifeMap)
  }
}

func printMap(arr [][]bool) {
  for i := range(arr) {
    for j := range(arr[i]) {
      if !arr[i][j] {
        fmt.Print("")
      } else {
        fmt.Print("󰝤")
      }
    }
    fmt.Println()
  }
} 

func makeLifeMap() [][]bool {
  if len(os.Args) < 3 {
    fmt.Println("Not enough args")
    os.Exit(1)
  }
  width, err := strconv.Atoi(os.Args[1])
  if err != nil {
    fmt.Println("Invalid argument 1")
    os.Exit(1)
  }
  height, err := strconv.Atoi(os.Args[2])
  if err != nil {
    fmt.Println("Invalid argument 1")
    os.Exit(1)
  }
  var lifeMap = make([][]bool, height)

  for i := 0; i < len(lifeMap); i++ {
    lifeMap[i] = make([]bool, width)
  }
  return lifeMap
}

func passFrame(lifeMap[][]bool) {
  var listOfChanges []Position
  for y := 0; y < len(lifeMap); y++ {
    for x := 0; x < len(lifeMap[y]); x++ {
      pos := processCell(x, y, lifeMap)
      if pos.X != -1 {
        listOfChanges = append(listOfChanges, pos)
      }
    }
  }
  time.Sleep(100 * time.Millisecond)
  for i := 0; i < len(listOfChanges); i++ {
    p := listOfChanges[i]
    lifeMap[p.Y][p.X] = p.Val
  }
}

func processCell(x, y int, lifeMap [][]bool) Position{
  neighbors := getNeighbors(x, y, lifeMap)
  if neighbors == 3 && !lifeMap[y][x]{
    // fmt.Println("Turning ",x,y,"alive")
    return Position{x, y, true}
  }
  if neighbors < 2 || neighbors > 3  && lifeMap[y][x]{
    // fmt.Println("Killing ",x,y)
    return Position{x, y, false}
  }
  return Position{-1, 0, false}
} 

func getNeighbors(x,y int, arr[][]bool) int {
  return checkLeft(x, y, arr) + checkRight(x, y, arr) + checkDown(x, y, arr) + checkUp(x, y, arr) + checkBottomCorners(x, y, arr) + checkTopCorners(x, y, arr)
}

func checkRight(x, y int, arr [][]bool) int{
  if x == len(arr[y]) - 1 {
    if !arr[y][0] {
      return 0
    }
    return 1
  }
  if !arr[y][x + 1] {
    return 0
  }
  return 1
}

func checkLeft(x, y int, arr [][]bool) int{
  if x == 0 {
    if !arr[y][len(arr[y]) - 1] {
      return 0
    }
    return 1
  }
  if !arr[y][x - 1] {
    return 0
  }
  return 1
}

func checkUp(x, y int, arr[][]bool) int {
  if y == 0 {
    if !arr[len(arr) - 1][x] {
      return 0
    }
    return 1
  }
  if !arr[y - 1][x] {
    return 0
  }
  return 1
}

func checkDown(x, y int, arr[][]bool) int {
  if y == len(arr) - 1 {
    if !arr[0][x] {
      return 0
    }
    return 1
  }
  if !arr[y + 1][x] {
    return 0
  }
  return 1
}

func checkBottomCorners(x, y int, arr[][]bool) int {
  counter := 0
  if y == len(arr) - 1 {
    counter += checkRight(x, 0, arr)
    counter += checkLeft(x, 0, arr)
    return counter
  }
  counter += checkRight(x, y + 1, arr)
  counter += checkLeft(x, y + 1, arr)
  return counter
}

func checkTopCorners(x, y int, arr [][]bool) int {
  counter := 0
  if y == 0 {
    counter += checkRight(x, len(arr) - 1, arr)
    counter += checkLeft(x, len(arr) - 1, arr)
    return counter
  }
  counter += checkRight(x, y - 1, arr)
  counter += checkLeft(x, y - 1, arr)
  return counter
}


