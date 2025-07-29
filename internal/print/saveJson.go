package print

import (
	"github.com/JuaanReis/vorin/internal/flags"
	"github.com/JuaanReis/vorin/internal/output"
	"github.com/JuaanReis/vorin/internal/model"
	"fmt"
	"os"
)

func SaveJson(cfg flags.CLIConfig, resultadoJson any) {
	if cfg.OutputFile != "" {
		err := output.SaveJson([]model.ResultadoJSON{}, cfg.OutputFile)
		if err != nil {
			fmt.Printf("Error saving JSON to %s: %v\n", cfg.OutputFile, err)
			os.Exit(1)
		}
		fmt.Printf("Results saved to %s\n", cfg.OutputFile)
	}
}