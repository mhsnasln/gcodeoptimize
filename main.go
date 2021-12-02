package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	// İşlem yapılacak dosyanın adını giriyor
	var import_name string

	// X değerlerinin artırılma x kadar artırılacak
	var x_value float64

	// Kullanıcıdan dosya adı istiyoruz!
	fmt.Print("İşlem yapılacak dosyanın adını girin:")
	fmt.Scanf("%s", &import_name)

	// X'ler ne kadar artırılacağını soruyoruz
	fmt.Print("X Değerini Girin:")
	fmt.Scanf("%f", &x_value)

	// Dosyamızın yolunu belirtiyoruz
	file, err := os.ReadFile(import_name)
	if err != nil {
		log.Fatalf("İşlem yapılacak dosya bulunamadı")
	}

	// Dosyamızın her bir satırını işlemek için alıyoruz
	scanner := bufio.NewScanner(strings.NewReader(string(file)))

	// Her satır dışarıya dönebilmesi için döngü dışında bir değişken tanımladım. İşlenmiş tüm satırlar burada döngünün tamamlanmasını bekleyeyip değişkene girecek
	response := make([]string, 10)

	// Dosya bitene kadar her bir satırda işlem tekrar edeiyor
	for scanner.Scan() {

		// Dosyadaki her bir satırı alıyoruz
		line := scanner.Text()

		// Satır işlemleraini ayrıca bir fonksiyonda yapıyorum
		optimize, err := worker(line, x_value)
		if err != nil {
			panic("Birşeyler ters gitti.")
		}

		// Burada worker fonksiyonundan gelen veriyi kontrol ediyorum boş geliyorsa ihtiyacımız olmadan veri göndermiştir. Boş, Açıklama satırı vs.
		if len(optimize) == 0 {
			continue
		}

		// Her satır işlendikten sonra yeniden kaydedilmesi için bir diziye aktarılıyor
		response = append(response, optimize)

	}

	// Buraya kadar geldiysek tüm satırlar işlenmiş ve kaydedilmeye hazır
	wf, err := os.Create("optimize_" + import_name)
	if err != nil {
		log.Fatalf("Bir sorun oluştu: %s", err)
	}

	datawriter := bufio.NewWriter(wf)

	// Yine satır satır, açıklama ve boş satırlar olmadan yeni kaydediyorum
	for _, data := range response {
		if len(data) == 0 {
			continue
		}
		_, err = datawriter.WriteString(data + "\n")
		if err != nil {
			log.Fatalf("Bir sorun oluştu: %s", err)
		}
	}

	datawriter.Flush()
	wf.Close()

	fmt.Println("Dosya sorunsuz dönüştürüldü!")

}

// bu fonksiyon gelen her satırı işleyip işlenmiş satırı geri gönderiyor
func worker(line string, x_value float64) (string, error) {

	// Gelen satırda ne kadar uzunlukta bir veri var eğer boş veriyse hemen işlem iptal ediyorum
	if len(line) == 0 {
		return "", nil
	}

	// Yine gelen veride açıklama satırı varsa doğrudan siliyorum
	if string(line[0]) == ";" {
		return "", nil
	}

	// Dosyadan gelen satırı boşluklarına göre böl
	split_line_with_space := strings.Split(line, " ")

	// Satırdaki verileri boşluğa göre parçalayıp herbirinde işlem yaptırıyorum
	for index, v := range split_line_with_space {

		// Yine aynı şekilde her satır parçası boş mu değil mi kontrol ediyorum
		if len(v) == 0 {
			continue
		}

		// Satır içindeki parçalardan her birinin ilk karakterine bakıyorum
		first_char := v[0:1]

		// İlk karakter boş mu kontrol ediyorum
		if string(first_char) == "" {
			split_line_with_space[index] = ""
		}

		// Burada belki parçaların ilkinde olmasa da 3. satırda açıklama satırı olabilir bunu siliyorum
		if string(first_char) == ";" {
			split_line_with_space[index] = ""
		}

		// İlk karakter x ise artırma işlemini yapıyorum
		if string(first_char) == "X" {

			// x'in yanındaki değeri x kadar artırıyorum
			x_val_float, err := strconv.ParseFloat(v[1:len(string(v))], 64)
			if err != nil {
				log.Fatalf("Hata Oluştu %s", err)
			}

			// Ve kaydediyorum yine aynı yerine
			float64ToString := fmt.Sprint(x_val_float + x_value)
			split_line_with_space[index] = "X" + float64ToString
		}

		// Satırın sonında E veya e harfi gördüğünde silecek
		if string(first_char) == "E" || string(first_char) == "e" {
			split_line_with_space[index] = ""
		}

		// ilk satır G0 ise sonuna "M9"
		if string(v) == "G0" {
			split_line_with_space = append(split_line_with_space, "M9")
		}

		// ilk satır G1 ise sonuna "M8"
		if string(v) == "G1" {
			split_line_with_space = append(split_line_with_space, "M8")
		}

	}

	// Paçalanan veri yeniden birleştiriyoruz
	join_current_line := strings.Join(split_line_with_space, " ")

	// İşlem yaptığımız satırı kaydetmek için gönderiyorum
	return join_current_line, nil

}
