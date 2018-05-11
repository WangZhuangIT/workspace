package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var filelist []string = make([]string, 0)

func main() {
	filepath.Walk("/tmp/zzdir/v3", walkfunc)

	if len(filelist) != 0 {
		fmt.Println("已经同步的文件：")
		for _, data := range filelist {
			fmt.Printf("%s \n", data)
		}
	} else {
		fmt.Println("本次不需要更新")
	}

}

func walkfunc(path string, info os.FileInfo, err error) error {
	oldpath := strings.Replace(path, "v3", "v2", -1)
	_, err = os.Stat(oldpath)
	f1, _ := os.Stat(path)

	if err != nil {
		//文件不存在，是文件夹直接创建，是文件直接创建读写文件
		// fmt.Println("文件不存在" + path)
		if f1.IsDir() {
			os.Mkdir(oldpath, 0777)
		} else {
			var file *os.File
			file, err = os.OpenFile(path, os.O_RDWR, 0)
			if err != nil {
				panic(err)
			}

			var oldfile *os.File
			oldfile, err = os.Create(oldpath)

			if err != nil {
				panic(err)
			}
			read2Old(file, oldfile)
			filelist = append(filelist, path)
		}

		return nil
	} else {
		//文件存在,是文件夹不操作，是文件对比md5
		// fmt.Println("文件存在" + oldpath)
		if !f1.IsDir() {
			md5File(path, oldpath, true)
		}
	}
	return nil
}

func md5File(path string, oldpath string, isExist bool) {

	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	var oldfile *os.File
	oldfile, err = os.OpenFile(oldpath, os.O_RDWR, 0)

	defer file.Close()
	defer oldfile.Close()

	if err != nil {
		panic(err)
	}

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		panic(err)
	}

	h2 := md5.New()
	_, err = io.Copy(h2, oldfile)

	if err != nil {
		panic(err)
	}

	fileMD5 := hex.EncodeToString(h.Sum(nil))
	oldfileMD5 := hex.EncodeToString(h2.Sum(nil))

	isEqual := strings.EqualFold(fileMD5, oldfileMD5)

	if isEqual {
		// fmt.Printf("%v 和 %v 的文件md5值相同：%x \n", path, oldpath, fileMD5)
	} else {
		// fmt.Println("*****************************二者出现不同的md5值*****************************")
		// fmt.Printf("%v 的文件md5值：%x \n", path, fileMD5)
		// fmt.Printf("%v 的文件md5值：%x \n", oldpath, oldfileMD5)
		read2Old(file, oldfile)
		filelist = append(filelist, path)
	}
}

func read2Old(file *os.File, oldfile *os.File) {
	buf := make([]byte, 1024)
	var context string
	file.Seek(0, 0)
	//将文件清空之后，需要移动光标到文件起始位置
	oldfile.Truncate(0)
	oldfile.Seek(0, 0)
	for {
		n, _ := file.Read(buf)
		if 0 == n {
			break
		}
		context += string(buf[:n])
	}
	_, err := oldfile.WriteString(context)

	if err != nil {
		panic(err)
	}

	// fmt.Println(strconv.Itoa(result) + "文件同步完毕")
}
