package collector

import (
	"strings"
	"github.com/JuaanReis/vorin/internal/analyzer"
)

func DataTargetFake(bodyAle []byte) (int, string) {
	stringBodyAle := string(bodyAle)
	structureOnly := analyzer.CleanStructure(stringBodyAle)
	fakeStructureSize := len(structureOnly)
	titleAle := strings.TrimSpace(strings.ToLower(analyzer.GetTitle(stringBodyAle)))
	return fakeStructureSize, titleAle
}