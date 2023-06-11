package main

import (
	"bufio"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/signintech/gopdf"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	for {
		pdfName := Scan(s, "保存するPDFファイル名を入力してください")
		imageCount, err := strconv.Atoi(Scan(s, "何枚撮影するかを入力してください"))
		if err != nil {
			fmt.Println("撮影枚数は数字で入力してください。")
			continue
		}
		imageTime, err := strconv.Atoi(Scan(s, "何秒間隔で撮影するか入力してください"))
		if err != nil {
			fmt.Println("撮影間隔は数字で入力してください。")
			continue
		}
		images := ScreenToImage(imageCount, imageTime)
		ImageToPDF(pdfName, images)

		c := Scan(s, "まだPDFを作成しますか？y/N")
		if strings.ToLower(c) != "y" {
			break
		}
	}
}

func Scan(s *bufio.Scanner, text string) string {
	for {
		fmt.Print(text, ": ")
		s.Scan()
		input := s.Text()
		if input != "" {
			return input
		}
	}
}

func ScreenToImage(imageCount, imageTime int) []string {
	var images []string

	fmt.Println("3秒後に撮影を開始します。")
	time.Sleep(3 * time.Second)
	for i := 0; i < imageCount; i++ {
		bounds := screenshot.GetDisplayBounds(1)
		img, _ := screenshot.CaptureRect(bounds)
		fileName := fmt.Sprintf("tmp_%d.png", i)
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)
		images = append(images, fileName)
		time.Sleep(time.Duration(imageTime) * time.Second)
	}
	return images
}

func ImageToPDF(pdfName string, images []string) {
	defer DeleteTmpPNG(images)
	pdf := gopdf.GoPdf{}
	SIZE := gopdf.Rect{W: 1920, H: 1080}
	pdf.Start(gopdf.Config{PageSize: SIZE})
	for _, image := range images {
		pdf.AddPage()
		pdf.Image(image, 0, 0, &SIZE)
	}
	pdf.WritePdf(fmt.Sprint(pdfName, ".pdf"))
}

func DeleteTmpPNG(images []string) {
	for _, image := range images {
		os.Remove(image)
	}
}
