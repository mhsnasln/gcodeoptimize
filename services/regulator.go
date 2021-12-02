package services

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/mhsnasln/gocode/models"
)

// Regülatör parçaları düzenlememize yarayacak
func Regulator(layer *models.Layer, export_name string) {

	var block = &models.Block{}
	var part = &models.Part{}

	// Her bir pointte yine işlem yapacağız
	for i := 0; i < len(layer.Items); i++ {

		item := layer.Items[i]

		// X ve Y'si 0 ise işimize yaramaz
		if item.X == 0 && item.Y == 0 {
			continue
		}

		// Z değeri var mı
		if item.Z != 0 {

			// z değerine özel gruplandırıyorum noktaları (bknz: models.point)
			if len(part.Items) > 0 {
				block.Parts = append(block.Parts, *part)
			}
			part.Z = item.Z
		}

		// Sonuncu ya geldinde kalan noktaları alıyoruz
		if i == len(layer.Items)-1 {
			block.Parts = append(block.Parts, *part)
		}

		// Tüm bunlarla bir katman oluşturuyoruz
		part.Items = append(part.Items, item)
	}

	// Tüm veriyi kaydedebilmek için bu isimde bir dosya oluşturuyoruz!
	wf, err := os.Create("./exports/optimized_" + export_name + ".gcode")
	if err != nil {
		log.Fatalf("Bir sorun oluştu: %s", err)
	}

	datawriter := bufio.NewWriter(wf)

	// Klasik Başlangıç Değerleri
	_, _ = datawriter.WriteString("M203" + "\n")
	_, _ = datawriter.WriteString("M203" + "\n")

	// Her bir z grubu için
	for _, p := range block.Parts {

		// Her bir nokta için
		for index, j := range p.Items {

			headOne := fmt.Sprintf("G%d X%f Y%f Z%f M%f", j.G, j.X, j.Y, p.Z, j.M1)
			_, _ = datawriter.WriteString(headOne + "\n")

			// İlk başlığın noktarı bittikten hemen sonra 2. makinenin kafa bilgilerini kaydediyorum
			if index == len(p.Items)-1 {
				for _, a := range p.Items {
					headTwo := fmt.Sprintf("G%d X%f Y%f Z%f M%f", a.G, a.XInc, a.Y, p.Z, a.M2)
					_, _ = datawriter.WriteString(headTwo + "\n")
				}
			}
		}

	}

	// Klasik Bitiş Değerleri
	_, _ = datawriter.WriteString("M213" + "\n")
	_, _ = datawriter.WriteString("M213" + "\n")

	// Dosya işlemlerini kapatıyorum
	datawriter.Flush()
	wf.Close()

}
