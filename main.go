package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mhsnasln/gocode/models"
	"github.com/mhsnasln/gocode/services"
)

// İşlem yapılacak dosyanın adını giriyor

func main() {

	var file_name string

	// X değerlerinin artırılma x kadar artırılacak
	var x_value float64

	// Kullanıcıdan dosya adı istiyoruz!
	fmt.Print("İşlem yapılacak dosyanın adı:")
	fmt.Scanf("%s", &file_name)

	// X'ler ne kadar artırılacağını soruyoruz
	fmt.Print("X Değerini Girin:")
	fmt.Scanf("%f", &x_value)

	// Dosyamızın yolunu belirtiyoruz
	new_name := fmt.Sprintf("%s.gcode", file_name)
	file, err := os.ReadFile(new_name)
	if err != nil {
		log.Fatalf("İşlem yapılacak dosya bulunamadı")
	}

	// Dosyamızın her bir satırını işlemek için alıyoruz
	scanner := bufio.NewScanner(strings.NewReader(string(file)))

	layer := &models.Layer{}

	layer.Items = []models.Point{}

	// Dosya bitene kadar her bir satırda işlem tekrar edeiyor
	for scanner.Scan() {

		// Dosyadaki her bir satırı alıyoruz
		line := scanner.Text()

		// Kafa1
		optimized_point, err := services.Compressor(line, x_value)
		if err != nil {
			panic("Birşeyler ters gitti.")
		}

		if optimized_point.M1 != 0 {
			layer.Items = append(layer.Items, *optimized_point)
		}

	}

	services.Regulator(layer, file_name)

	fmt.Println("Dosya başarıyla oluşturuldu!")

}
