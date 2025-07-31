package print

import (
	"github.com/JuaanReis/vorin/internal/flags"
	"github.com/JuaanReis/vorin/internal/output"
	"github.com/JuaanReis/vorin/internal/model"
	"fmt"
	"os"
)

func SaveJson(cfg flags.CLIConfig, resultadoJson any) {
	res, ok := resultadoJson.([]model.ResultadoJSON)
	if !ok {
		fmt.Println("[ERROR] Failed to convert resultJson to []model.ResultadoJSON")
		return
	}

	err := output.SaveJson(res, cfg.OutputFile)
	if err != nil {
		fmt.Printf("Error saving JSON to %s: %v\n", cfg.OutputFile, err)
		os.Exit(1)
	}
	fmt.Printf("Result saved in %s\n", cfg.OutputFile)
}