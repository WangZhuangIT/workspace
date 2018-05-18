package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"runtime"
	"strconv"

	"github.com/disintegration/imaging"
)

//该工具支持将图片色彩反转，图片灰化，图片转为字符画。
//author iccboy 2017-9-2
func main() {
	args := os.Args //获取用户输入的所有参数
	if args == nil || len(args) != 4 {
		usage()
		return
	}
	fmt.Println("...转换中...")
	source := args[1]
	dst := args[2]
	size, _ := strconv.ParseFloat(args[3], 64)

	tmp := resizepng(source, size)

	png2ascii(tmp, dst)
	fmt.Println("转换完成...")
}

func usage() {
	fmt.Println("请输入参数： 源文件	目标文件	缩放比例")
}

func resizepng(file string, size float64) string {
	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	img, err := imaging.Open(file)
	if err != nil {
		panic(err)
	}
	x := int(float64(img.Bounds().Dx()) * size)
	y := int(float64(img.Bounds().Dy()) * size)
	thumb := imaging.Thumbnail(img, x, y, imaging.CatmullRom)

	// save the combined image to file
	tmp := "tmp.png"
	err = imaging.Save(thumb, tmp)
	if err != nil {
		panic(err)
	}
	return tmp
}

func png2ascii(source string, dst string) {
	arr := []string{"M", "N", "H", "Q", "$", "O", "C", "?", "7", ">", "!", ":", "–", ";", "."}
	ff, _ := os.Open(source)
	m, _, _ := image.Decode(ff)
	x := m.Bounds().Dx()
	y := m.Bounds().Dy()
	file, _ := os.Create(dst)
	defer file.Close()
	for j := 0; j < y; j += 2 {
		for i := 0; i < x; i++ {
			colorRgb := m.At(i, j)
			r, g, b, _ := colorRgb.RGBA()
			r = r & 0xFF
			g = g & 0xFF
			b = b & 0xFF
			gray := 0.299*float64(r) + float64(g)*0.587 + float64(b)*0.114
			temp := fmt.Sprintf("%.0f", gray*float64(len(arr)+1)/255)
			index, _ := strconv.Atoi(temp)
			if index >= len(arr) {
				file.WriteString(" ")
			} else {
				file.WriteString(arr[index])
			}
		}
		file.WriteString("\n")
	}
}
