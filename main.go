package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	// Dosyamızın yolunu belirtiyoruz
	file, _ := os.Open("./import.gcode")

	// Dosyamızın her bir satırını bize text olarak yansıtacak
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		// Dosyadaki her bir satırı alıyoruz
		line := scanner.Text()

		// Satır işlemlerini ayrıca bir fonksiyonda yapıyorum
		optimze, _ := worker(line)

		// Boş geldi dolayısı ile açıklama satırı veya hakkaten boş satır
		if len(optimze) == 0 {
			continue
		}

		// Kontrol için ekrana yazdırıyoruz
		fmt.Println(optimze)

	}

}

func worker(line string) (string, error) {

	// ilk karakter ";"'ise satırı göz ardı et
	if string(line[0]) == ";" {
		return "", nil
	}

	// Dosyadan gelen satırı boşluklarına göre böl
	split_line_with_space := strings.Split(line, " ")

	// Satırın son kısımda M ile başlayan satır varsa onu sil
	lastline := split_line_with_space[len(split_line_with_space)-1]

	// Satırın sonında E veya e harfi gördüğünde silecek
	if string(lastline[0]) == "E" || string(lastline[0]) == "e" {
		split_line_with_space = append(split_line_with_space[:len(split_line_with_space)-1], split_line_with_space[len(split_line_with_space)-1+1:]...)
	}

	// ilk satır g0 ise sonuna "M9"
	if string(split_line_with_space[0]) == "G0" {
		split_line_with_space = append(split_line_with_space, "M9")
	}

	// ilk satır g1 ise sonuna "M8"
	if string(split_line_with_space[0]) == "G1" {
		split_line_with_space = append(split_line_with_space, "M8")
	}

	// Parçaladığım veriyi yeniden birleştiriyoruz
	join_current_line := strings.Join(split_line_with_space, " ")

	// Güncellenen satırı yeniden gönderiyorum
	return join_current_line, nil

}
