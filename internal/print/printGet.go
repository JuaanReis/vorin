package print

import (
	"fmt"
	"github.com/JuaanReis/vorin/internal/model"
)

func PrintGet(resultado []model.Resultado, verbose bool) {
	for _, v := range resultado {
		if !verbose {
			fmt.Printf("%s[%3d]%s  %-20s Words: %-6d Size: %-6dB Lines: %-5d %-6s %-11s\n",
				v.Color, v.Status, Reset,
				v.URL, v.Text, v.Size, v.Lines, v.Time, v.Label)
		} else {
			fmt.Printf("%s[%3d]%s  %4dw  %5dB  %4dL  %6s  %s\n",
				v.Color, v.Status, Reset,
				v.Text, v.Size, v.Lines, v.Time, v.Label)
			fmt.Printf(" ├─ URL     : %s\n", v.Endereco)
			fmt.Printf(" ├─ FUZZ    : %s\n", v.URL)
			if v.Title != "" {
				fmt.Printf(" └─ Title   : %s\n\n", v.Title)
			}
		}
	}
}