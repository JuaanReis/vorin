package collector

import (
	"strings"
	"github.com/JuaanReis/vorin/internal/model"
)

func IsSameContent(title string, titleAle string, cfg model.ParserConfigGet, htmlSize int, lines int, structureSize int, fakeStructureSize int, content string) bool {
	titleMatch := title == titleAle || title == cfg.FilterTitle
	sizeTooSmall := htmlSize <= cfg.FilterSize
	linesTooFew := lines <= cfg.FilterLine || lines == 0
	structureMatch := structureSize == fakeStructureSize
	erroBody := false
	if cfg.FilterBody != "" {
		erroBody = strings.Contains(strings.ToLower(content), strings.ToLower(cfg.FilterBody))
	}
	isSameContent := titleMatch || structureMatch || sizeTooSmall || linesTooFew || erroBody
	return isSameContent
}