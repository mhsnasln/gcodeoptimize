package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mhsnasln/gocode/internal/lines"
	"github.com/mhsnasln/gocode/internal/models"
)

func main() {

	// İşlem yapılacak dosyanın adını giriyor
	var file_name string

	// X değerlerinin artırılma x kadar artırılacak
	var x_value float64

	// Kullanıcıdan dosya adı istiyoruz!
	fmt.Print("İşlem yapılacak dosyanın adı:")
	fmt.Scanf("%s", &file_name)

	// X'ler ne kadar artırılacağını soruyoruz
	fmt.Print("X Değerini Girin:")
	fmt.Scanf("%f", &x_value)

	// Kullanıcıdan aldığımız dosya adını okumaya başlıyoruz.
	file, err := os.ReadFile(fmt.Sprintf("%s.gcode", file_name))
	if err != nil {
		log.Fatalf("İşlem yapılacak dosya bulunamadı")
	}

	// Dosyamızın her bir satırını işlemek için alıyoruz
	scanner := bufio.NewScanner(strings.NewReader(string(file)))

	// Boş layerlar oluşturuyorum
	layer := &models.Layer{}
	layer.Items = []models.Point{}

	// Dosya bitene kadar her bir satırda işlem tekrar edeiyor
	for scanner.Scan() {

		// Dosyadaki her bir satırı alıyoruz
		line := scanner.Text()

		// Her bir ihtiyacımız olan satır geliyor
		optimized_point, err := lines.Compressor(line, x_value)
		if err != nil {
			panic("Compressorde sorun var!")
		}

		// x ve y değerleri boş ise bu satırı atlıyoruz
		if optimized_point.X == 0 && optimized_point.Y == 0 {
			continue
		}

		// Başarılı olan satırlar layer.Items'a katılıyor
		layer.Items = append(layer.Items, *optimized_point)

	}

	// her bir satırı Regülatör'e gönderiyorum. z değerlerine göre gruplandırması ve kayıt işlemlerinin yapılması için
	lines.Regulator(layer, file_name)

	// Uygulama burada mesaj verip kapanıyor!
	fmt.Println("Dosya başarıyla oluşturuldu!")

}
