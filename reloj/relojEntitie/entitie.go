package relojEntitie

import "time"

type User struct {
	NroDoc    string
	CodUs     string
	Psw       string
	Nombres   string
	Apellido  string
	Operation string
	In1       string
	In2       string
	Flexible  int
}

type Area struct {
	Supervisor string
	Encargado  string
}

type Reloj struct {
	CodUs  string
	Ing    time.Time
	Sale   time.Time
	Estado string
	Tiempo time.Time
}
