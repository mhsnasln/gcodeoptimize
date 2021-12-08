package lines

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mhsnasln/gocode/internal/models"
)

var Temp_z = 0.0

// Fonksiyon amacı gelen satırları daha özgür bir şekilde manipüle etmem için
func Compressor(line string, x_value float64) (*models.Point, error) {

	// Nokta oluşturuyorum
	point := models.Point{}

	// Gelen satırda ne kadar uzunlukta bir veri var eğer boş veriyse hemen işlem iptal ediyorum
	if len(line) == 0 || string(line[0]) == ";" || (string(line[0:2]) != "G0" && string(line[0:2]) != "G1") {
		return &point, nil
	}

	// Dosyadan gelen satırı boşluklarına göre böl (X12, Y12, Z123)
	split_line_with_space := strings.Split(line, " ")

	// Satırdaki verileri boşluğa göre parçalayıp herbirinde işlem yaptırıyorum
	for _, v := range split_line_with_space {

		// İlk karakter boş mu kontrol ediyorum
		if len(v) == 0 {
			continue
		}

		// İlk karakter boş mu kontrol ediyorum
		if string(v) == "" {
			continue
		}

		// Satır içindeki parçalardan her birinin ilk karakterine bakıyorum
		first_char := v[0:1]

		// Burada belki parçaların ilkinde olmasa da 3. satırda açıklama satırı olabilir bunu siliyorum
		if string(first_char) == ";" {
			for i := 0; i < len(split_line_with_space); i++ {
				continue
			}
		}

		// ilk satır G0 ise sonuna "M9"
		if string(v) == "G0" {
			point.G = 0
			point.M1 = 205
			point.M2 = 206
		}

		// ilk satır G1 ise sonuna "M8"
		if string(v) == "G1" {
			point.G = 1
			point.M1 = 215
			point.M2 = 216
		}

		// İlk karakter x ise artırma işlemini yapıyorum
		if string(first_char) == "X" {

			// x'in yanındaki değeri x kadar artırıyorum
			x_val_float, err := strconv.ParseFloat(v[1:len(string(v))], 64)
			if err != nil {
				log.Fatalf("Hata Oluştu %s", err)
			}

			point.X = x_val_float
			point.XInc = x_val_float + x_value

		}

		// İlk karakter y ise verimi alıyorum
		if string(first_char) == "Y" {

			// x'in yanındaki değeri x kadar artırıyorum
			y_val_float, err := strconv.ParseFloat(v[1:len(string(v))], 64)
			if err != nil {
				log.Fatalf("Y Değeri alınırken hata Oluştu %s", err)
			}

			point.Y = y_val_float

		}

		// İlk karakter z ise verimi alıyorum
		if string(first_char) == "Z" {

			// x'in yanındaki değeri x kadar artırıyorum
			z_val_float, err := strconv.ParseFloat(v[1:len(string(v))], 64)
			if err != nil {
				log.Fatalf("Z Değeri alınırken hata Oluştu %s", err)
			}

			Temp_z = z_val_float
			point.Z = Temp_z

		} else {
			point.Z = Temp_z
		}

	}

	// İşlem yaptığımız satırı kaydetmek için gönderiyorum
	return &point, nil

}

// Regülatör parçaları düzenlememize yarayacak
func Regulator(layer *models.Layer, export_name string) {

	var block = &models.Block{}
	var part = models.Part{}

	// layer uzunluğu
	var layer_length = len(layer.Items)

	// ilk karşılacağımız point'in z değerini alıyoruz
	var current_z float64 = layer.Items[0].Z

	// Her bir pointte yine işlem yapacağız
	for i := 0; i < layer_length; i++ {

		// Her bir point
		item := layer.Items[i]

		if item.Z != current_z {

			part.Z = layer.Items[i-1].Z
			block.Parts = append(block.Parts, part)
			current_z = item.Z

			part = models.Part{}
			part.Items = append(part.Items, item)

		} else {
			part.Items = append(part.Items, item)
		}

		if i == layer_length-1 {
			part.Z = layer.Items[i].Z
			block.Parts = append(block.Parts, part)
		}

	}

	// Tüm veriyi kaydedebilmek için bu isimde bir dosya oluşturuyoruz!
	wf, err := os.Create("./exports/optimized_" + export_name + ".gcode")
	if err != nil {
		log.Fatalf("Exports klasörü bulunamadı!")
	}

	datawriter := bufio.NewWriter(wf)

	// Klasik Başlangıç Değerleri
	_, err = datawriter.WriteString("G28 G21 G90" + "\n")
	if err != nil {
		log.Fatalf("G28 G21 G90 satırı projeye eklenmedi")
	}

	_, err = datawriter.WriteString("M203" + "\n")
	if err != nil {
		log.Fatalf("M203 satırı projeye eklenmedi")
	}

	// Her bir z grubu için
	for _, p := range block.Parts {

		// Her bir nokta için
		for index, j := range p.Items {

			headOne := fmt.Sprintf("G%d X%f Y%f Z%f M%f", j.G, j.X, j.Y, p.Z, j.M1)
			_, err = datawriter.WriteString(headOne + "\n")
			if err != nil {
				log.Fatalf("Kafa1 kodlarında bir sorun oluştu: %s", err)
			}

			// İlk başlığın noktarı bittikten hemen sonra 2. makinenin kafa bilgilerini kaydediyorum
			if index == len(p.Items)-1 {
				for _, a := range p.Items {
					headTwo := fmt.Sprintf("G%d X%f Y%f Z%f M%f", a.G, a.XInc, a.Y, p.Z, a.M2)
					_, err = datawriter.WriteString(headTwo + "\n")
					if err != nil {
						log.Fatalf("Kafa2 kodlarında bir sorun oluştu: %s", err)
					}
				}
			}
		}

	}

	// Klasik Bitiş Değerleri
	_, err = datawriter.WriteString("M213" + "\n")
	if err != nil {
		log.Fatalf("M213 satırı projeye eklenmedi")
	}

	_, err = datawriter.WriteString("M30" + "\n")
	if err != nil {
		log.Fatalf("M30 satırı projeye eklenmedi")
	}

	// Dosya işlemlerini kapatıyorum
	datawriter.Flush()
	wf.Close()

}
