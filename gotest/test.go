package main

import (
	"fmt"
	"strings"
)

func main() {
	dir := "/this/is/a/msg////////////////"
	l := strings.TrimRight(dir, "/")
	fmt.Println(l)
}
