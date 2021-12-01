package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	// İşlem yapılacak dosyanın adını giriyor
	var import_name string

	// İşlem yapıldıktan sonra yeni dosyanın ismi
	var export_name string

	// X değerlerinin artırılma x kadar artırılacak
	var x_value float64

	// Kullanıcıdan dosya adı istiyoruz!
	fmt.Print("İşlem yapılacak dosyanın adını girin:")
	fmt.Scanf("%s", &import_name)

	// Dosya işlem bittikten sonra hangi isimle kaydedilecek
	fmt.Print("Yeni dosyanın ismini girin:")
	fmt.Scanf("%s", &export_name)

	// X'ler ne kadar artırılacağını soruyoruz
	fmt.Print("X Değerini Girin:")
	fmt.Scanf("%s", &x_value)

	// Dosyamızın yolunu belirtiyoruz
	file, err := os.ReadFile(import_name)
	if err != nil {
		log.Fatalf("İşlem yapılacak dosya bulunamadı: %s", err)
	}

	// Dosyamızın her bir satırını işlemek için alıyoruz
	scanner := bufio.NewScanner(strings.NewReader(string(file)))

	// Her satır dışarıya dönebilmesi için döngü dışında bir değişken tanımladım. İşlenmiş tüm satırlar burada döngünün tamamlanmasını bekleyecek
	response := make([]string, 10)

	// Dosya bitene kadar her bir satırda işlem tekrar edeiyor
	for scanner.Scan() {

		// Dosyadaki her bir satırı alıyoruz
		line := scanner.Text()

		// Satır işlemleraini ayrıca bir fonksiyonda yapıyorum
		optimize, err := worker(line)
		if err != nil {
			panic("Birşeyler ters gitti.")
		}

		// Kontrol ediyorum satır boş veya açıklama satırı ise atlamasını söylüyorum
		if len(optimize) == 0 {
			continue
		}

		// Her satır işlendikten sonra yeniden kaydedilmesi için bir diziye aktarılıyor
		response = append(response, optimize)

	}

	// Buraya kadar geldiysek tüm satırlar işlenmiş ve kaydedilmeye hazır
	wf, err := os.Create(export_name)
	if err != nil {
		log.Fatalf("Bir sorun oluştu: %s", err)
	}

	datawriter := bufio.NewWriter(wf)

	// Yine satır satır, açıklama ve boş satırlar olmadan yeni kaydediyorum
	for _, data := range response {
		if len(data) == 0 {
			continue
		}
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	wf.Close()

	fmt.Println("Dosya başarıyla dönüştürüldü!")

}

// bu fonksiyon gelen her satırı işleyip işlenmiş satırı geri gönderiyor
func worker(line string) (string, error) {

	// Dosyadan gelen satırı boşluklarına göre böl
	split_line_with_space := strings.Split(line, " ")

	// Kuralımız işlem yapabilmem için ilk karakter her zaman G
	if string(line[0]) == ";" {
		return "", nil
	}

	// Satırın son kısımda M ile başlayan satır varsa onu sil
	lastline := split_line_with_space[len(split_line_with_space)-1]

	// X'ler belirtildiği kadar artırılacak
	for _, v := range split_line_with_space {

		first_char := v[0]

		if string(first_char) != "X" {
			continue
		}

		// TODO: Bu alanda kullanıcıdan aldığım x değerini her x satırına ekleyeceğim

	}

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

	// Paçalanan veri yeniden birleştiriyoruz
	join_current_line := strings.Join(split_line_with_space, " ")

	// İşlem yaptığımız satırı kaydetmek için gönderiyorum
	return join_current_line, nil

}
