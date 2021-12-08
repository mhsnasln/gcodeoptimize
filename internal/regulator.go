package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Regülatör parçaları düzenlememize yarayacak
func Regulator(layer *internal.Layer, export_name string) {

	var block = &internal.Block{}
	var part internal.Part

	// layer uzunluğu
	var layer_length = len(layer.Items)

	// ilk karşılacağımız point'in z değerini alıyoruz
	var current_z float64 = -1.00

	// Her bir pointte yine işlem yapacağız
	for i := 0; i < layer_length; i++ {

		// Her bir point
		item := layer.Items[i]

		if i == 0 {
			fmt.Println("For i == 0")
			part = internal.Part{}
		} else if item.Z != current_z {
			part.Z = item.Z
			block.Parts = append(block.Parts, part)
			current_z = item.Z
			fmt.Println("For block")
			part = internal.Part{}
		} else {
			part.Items = append(part.Items, item)
		}

	}

	// Tüm veriyi kaydedebilmek için bu isimde bir dosya oluşturuyoruz!
	wf, err := os.Create("./exports/optimized_" + export_name + ".gcode")
	if err != nil {
		log.Fatalf("exports klasörü bulunamadı!")
	}

	datawriter := bufio.NewWriter(wf)

	// Klasik Başlangıç Değerleri
	_, _ = datawriter.WriteString("G28 G21 G90" + "\n")
	_, _ = datawriter.WriteString("M203" + "\n")

	// Her bir z grubu için
	for _, p := range block.Parts {

		// Her bir nokta için
		for index, j := range p.Items {

			headOne := fmt.Sprintf("G%d X%f Y%f Z%f M%f", j.G, j.X, j.Y, p.Z, j.M1)
			_, err = datawriter.WriteString(headOne + "\n")
			if err != nil {
				log.Fatalf("Bir sorun oluştu: %s", err)
			}

			// İlk başlığın noktarı bittikten hemen sonra 2. makinenin kafa bilgilerini kaydediyorum
			if index == len(p.Items)-1 {
				for _, a := range p.Items {
					headTwo := fmt.Sprintf("G%d X%f Y%f Z%f M%f", a.G, a.XInc, a.Y, p.Z, a.M2)
					_, err = datawriter.WriteString(headTwo + "\n")
					if err != nil {
						log.Fatalf("Bir sorun oluştu: %s", err)
					}
				}
			}
		}

	}

	// Klasik Bitiş Değerleri
	_, _ = datawriter.WriteString("M213" + "\n")
	_, _ = datawriter.WriteString("M30" + "\n")

	// Dosya işlemlerini kapatıyorum
	datawriter.Flush()
	wf.Close()

}
