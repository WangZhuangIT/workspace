package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func ReadFrom(filename string) (image.Image, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	fmt.Println(file)
	return image.Decode(file)
}

func SaveJPEG(m image.Image, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, m, &jpeg.Options{Quality: 100})
}

func main() {
	_, _, err := ReadFrom("./cc.jpg")
	if err != nil {
		fmt.Println("ReadFrom decode err")
		return
	}
	// bounds := m.Bounds()

	// mgrey := image.NewGray(bounds)
	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
	// 		r, g, b, _ := m.At(x, y).RGBA()
	// 		r = r >> 8
	// 		g = g >> 8
	// 		b = b >> 8
	// 		v := uint8((float32(r)*299 + float32(g)*587 + float32(b)*114) / 1000)
	// 		mgrey.Set(x, y, color.Gray{v})
	// 	}
	// }

	// SaveJPEG(mgrey, "./example1.jpg")
}
