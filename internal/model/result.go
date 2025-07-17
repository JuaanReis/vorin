package model

import (
	"time"
)

type Resultado struct {
	Status   int
	Resposta string
	URL      string
	Title    string
	Text     int
	Size     int
	Lines    int
	Time     time.Duration
	Label    string
	Color    string
	User     string
	Pass     string
	Endereco string
}
