package main

import (
	"bufio"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

//重写reader用于实现显示下载进度
//type Reader struct {
//	io.Reader
//	//bufio.Reader
//	Total   int64
//	Current int64
//}
//
//func (r *Reader) Read(p []byte) (n int, err error) {
//	n, err = r.Reader.Read(p)
//	r.Current += int64(n)
//	fmt.Printf("\r当前进度 %.2f\n", float64((r.Current*10000)/r.Total)/100)
//	return
//}

var wg sync.WaitGroup

func main() {
	//从命令行中获取网站
	//baseUrl := os.Args[1]
	baseUrl := "https://www.bizhi88.com/"
	dirPath := "./temp/"
	//若目标图片没有后缀格式  使用图片格式png
	imgFormat := ".png"
	imgUrl, _ := getImgUrl(baseUrl)
	//imgUrl := []string{"http://img.sccnn.com/bimg/340/03821.jpg"}
	_ = createDir(dirPath)
	for index, url := range imgUrl {
		wg.Add(1)
		go downloadFile(url, dirPath, imgFormat, index)
	}
	wg.Wait()
}

//获取图片地址
func getImgUrl(baseUrl string) (imgUrl []string, err error) {
	doc, err := htmlquery.LoadURL(baseUrl)
	if err != nil {
		panic(err)
		return nil, err
	}

	list := htmlquery.Find(doc, "//img")

	for _, item := range list {
		//imgUrl = append(imgUrl, htmlquery.SelectAttr(item, "src"))
		//对于https://www.bizhi88.com/     图片真是地址位于标签data-original中
		if temp := htmlquery.SelectAttr(item, "data-original"); temp != "" {
			//避免某些图片地址位于标签src对于的值    取data-original则会为空值
			imgUrl = append(imgUrl, temp)
		}
	}
	return
}

func downloadFile(url string, dirPath string, imgFormat string, index int) {
	fileName := path.Base(url)
	//使用（保留）目标图片固有格式
	if strings.Contains(fileName, ".jpg") || strings.Contains(fileName, ".png") || strings.Contains(fileName, ".jpeg") || strings.Contains(fileName, ".gif") || strings.Contains(fileName, ".bmp") {
		imgFormat = ""
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(url)
		fmt.Println(index)
		panic(err)
		return
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}()
	// 获得get请求响应的reader对象
	bfReader := bufio.NewReaderSize(res.Body, 32*1024)

	//重写reader用于实现显示下载进度
	//reader := &Reader{
	//	Reader: bfReader,
	//	Total:  res.ContentLength,
	//}

	file, err := os.Create(dirPath + strconv.Itoa(index) + fileName + imgFormat)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	_, _ = io.Copy(writer, bfReader)

	//重写reader用于实现显示下载进度
	//_, _ = io.Copy(writer, reader)

	wg.Done()
}

//创建文件夹
func createDir(path string) (err error) {
	_, err = os.Stat(path) //   "./temp/"
	if err != nil {
		//文件夹不存在 创建文件夹
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return
}
