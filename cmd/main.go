package main

import (
	"fmt"
	"github.com/skvoch/burlington-backend/tree/master/internal/qr"
)

func main() {
	img, _ := qr.Generate("1")
	qr.ConsolePngWriter(img)
	fmt.Println("hello world")
}



