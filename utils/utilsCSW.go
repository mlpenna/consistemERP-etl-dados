package utils

import (
	"time"

	"github.com/go-vgo/robotgo"
)

func CswLogin() {
	robotgo.KeyTap("f5")
	robotgo.TypeStr("RPA")
	robotgo.KeyTap("tab")
	robotgo.TypeStr("123456")
	robotgo.MoveMouse(656, 496)
	robotgo.MouseClick("left", false)
}

func CswLogout() {
	robotgo.MoveMouse(1295, 111) //botao exportar CSV
	robotgo.MouseClick("left", false)
	robotgo.MoveMouse(1346, 303) //botao exportar CSV
	robotgo.MouseClick("left", false)
}

func CswAbrirRotina(rotina string) {
	robotgo.MoveMouse(145, 320) //botao execução direta
	robotgo.MouseClick("left", false)
	robotgo.TypeStr(rotina)
	robotgo.MoveMouse(820, 398) //botao execução direta
	robotgo.MouseClick("left", false)
}

func CswReloadBrowser() {
	robotgo.MoveMouse(87, 81) //botao atualiza navegador (rotina travada)
	robotgo.MouseClick("left", false)
}
