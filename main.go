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
  
  fmt.Println(lifeMap)
}
