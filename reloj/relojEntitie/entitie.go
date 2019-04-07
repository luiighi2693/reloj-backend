package relojEntitie

import "time"

type User struct {
	NroDoc    string
	CodUs     string
	Psw       string
	Nombres   string
	Apellido  string
	Operation string
}

type Reloj struct {
	CodUs  string
	Ing    time.Time
	Sale   time.Time
	Estado string
	Tiempo time.Time
}
