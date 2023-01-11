package main

import (
	"bufio"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/signintech/gopdf"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	for {
		PdfName := Scan(s, "保存するPDFファイル名を入力してください")
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
		ImageToPdf(PdfName, images)

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

	fmt.Println(imageCount, "秒後に撮影を開始します。")
	for i := 0; i < imageCount; i++ {
		for j := imageTime; j > 0; j-- {
			time.Sleep(time.Second)
			if j <= 3 {
				fmt.Println(j)
			}
		}
		bounds := screenshot.GetDisplayBounds(1)
		img, _ := screenshot.CaptureRect(bounds)
		fileName := fmt.Sprintf("tmp_%d.png", i)
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)
		images = append(images, fileName)
	}
	return images
}

func ImageToPdf(PdfName string, images []string) {
	Pdf := gopdf.GoPdf{}
	SIZE := gopdf.Rect{W: 1920, H: 1080}
	Pdf.Start(gopdf.Config{PageSize: SIZE})
	for _, image := range images {
		Pdf.AddPage()
		Pdf.Image(image, 0, 0, &SIZE)
	}
	Pdf.WritePdf(fmt.Sprint(PdfName, ".pdf"))
	DeleteTmpPNG(images)
}

func DeleteTmpPNG(images []string) {
	for _, image := range images {
		os.Remove(image)
	}
}
