package print

import (
	"fmt"
	"github.com/JuaanReis/vorin/internal/model"
)

func PrintStatusOnly(resultado []model.Resultado, method string) {
	for _, v := range resultado {
		if method == "POST" {
			fmt.Printf("%s[%3d]%s user=%s pass=%s\n", v.Color, v.Status, Reset, v.User, v.Pass)
		} else {
			fmt.Printf("%s[%3d]%s %-26s\n", v.Color, v.Status, Reset, v.URL)
		}
	}
}