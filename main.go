package main

import (
  "fmt"
  // "sync"
  "strconv"
  "os"
)

func main(){
  if len(os.Args) < 3 {
    fmt.Println("Not enough args")
    return
  }
  width, err := strconv.Atoi(os.Args[1])
  if err != nil {
    fmt.Println("Invalid argument 1")
    return
  }
  height, err := strconv.Atoi(os.Args[2])
  if err != nil {
    fmt.Println("Invalid argument 1")
    return
  }
  var lifeMap = make([][]bool, height)

  for i := 0; i < len(lifeMap); i++ {
    lifeMap[i] = make([]bool, width)
  }
  
  
}

func processCell(x, y int, lifeMap [][]bool) {
  neighbors := getNeighbors(x, y, lifeMap)
  if neighbors == 3 {
    lifeMap[y][x] = true
  }
  if neighbors < 2 || neighbors > 3 {
    lifeMap[y][x] = false
  }
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
    if !arr[y][len(arr[y])] {
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
    if !arr[len(arr)][x] {
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


