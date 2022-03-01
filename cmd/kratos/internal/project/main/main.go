package main

import (
	"fmt"
	"path/filepath"
)

func main(){
	a := "/Users/ningzi/workspace/kratos/cmd/kratos/internal/project/main/main.go"
	a = filepath.Join(a, "..")
	fmt.Println(filepath.Join(a, ".."))
}