package services

import (
	"log"
	"strconv"
	"strings"

	"github.com/mhsnasln/gocode/models"
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
