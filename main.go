package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	var import_name string
	var export_name string

	fmt.Print("İşlem yapılacak dosyanın adını girin:")
	fmt.Scanf("%s", &import_name)
	fmt.Print("Yeni dosyanın ismini girin:")
	fmt.Scanf("%s", &export_name)

	// Dosyamızın yolunu belirtiyoruz
	file, err := os.ReadFile(import_name)
	if err != nil {
		log.Fatalf("İşlem yapılacak dosya bulunamadı: %s", err)
	}

	// Dosyamızın her bir satırını bize text olarak yansıtacak
	scanner := bufio.NewScanner(strings.NewReader(string(file)))

	response := make([]string, 10)

	for scanner.Scan() {

		// Dosyadaki her bir satırı alıyoruz
		line := scanner.Text()

		// Satır işlemlerini ayrıca bir fonksiyonda yapıyorum
		optimize, err := worker(line)
		if err != nil {
			panic("Birşeyler ters gitti.")
		}

		// Boş geldi dolayısı ile açıklama satırı veya hakkaten boş satır
		if len(optimize) == 0 {
			continue
		}

		response = append(response, optimize)

	}

	wf, err := os.Create(export_name)
	if err != nil {
		log.Fatalf("Bir sorun oluştu: %s", err)
	}

	datawriter := bufio.NewWriter(wf)

	for _, data := range response {
		if len(data) == 0 {
			continue
		}
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	wf.Close()

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
