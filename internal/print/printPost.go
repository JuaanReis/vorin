package print

import (
	"fmt"
	"github.com/JuaanReis/vorin/internal/model"
)

func PrintPost(resultado []model.Resultado, verbose bool) {
	for _, v := range resultado {
		if !verbose {
			fmt.Printf("%s[%3d]%s user=%-10s pass=%-10s Size: %-6dB Lines: %-5d %-6s %-11s\n",
				v.Color, v.Status, Reset,
				v.User, v.Pass, v.Size, v.Lines, v.Time, v.Label)
		} else {
			fmt.Printf("%s[%3d]%s  %4dw  %5dB  %4dL  %6s  %s\n",
				v.Color, v.Status, Reset,
				v.Text, v.Size, v.Lines, v.Time, v.Label)
			fmt.Printf(" ├─ FUZZ    : %s | %s\n", v.User, v.Pass)
			if v.Title != "" {
				fmt.Printf(" └─ Title   : %s\n\n", v.Title)
			}
		}
	}
}