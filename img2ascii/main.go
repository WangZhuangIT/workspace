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
	source := args[1]
	dst := args[2]
	size, _ := strconv.ParseFloat(args[3], 64)
	if size <= 0 {
		fmt.Println("缩放比例参数异常")
		return
	}
	fmt.Println("...转换中...")

	thumb := getThumb(source, size)
	png2ascii(thumb, dst)
	fmt.Println("...转换完成...")
}

func usage() {
	fmt.Println("请输入参数：1源文件  2目标文件  3缩放比例")
}

func getThumb(file string, size float64) *image.NRGBA {
	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())
	img, err := imaging.Open(file)
	if err != nil {
		panic(err)
	}
	x := int(float64(img.Bounds().Dx()) * size)
	y := int(float64(img.Bounds().Dy()) * size)
	thumb := imaging.Thumbnail(img, x, y, imaging.CatmullRom)
	return thumb
}

func png2ascii(m *image.NRGBA, dst string) {
	arr := []string{"M", "N", "H", "Q", "$", "O", "C", "?", "7", ">", "!", ":", "–", ";", "."}
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
