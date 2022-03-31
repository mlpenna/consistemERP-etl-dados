package utils

import (
	"time"

	"github.com/go-vgo/robotgo"
)

func CswLogin() {

	// CSW Login
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("f5")
	time.Sleep(5 * time.Second)
	robotgo.TypeStr("RPA")
	time.Sleep(3 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(1 * time.Second)
	robotgo.TypeStr("123456")
	time.Sleep(5 * time.Second)
	robotgo.MoveMouse(656, 496)
	time.Sleep(5 * time.Second)
	robotgo.MouseClick("left", false)
	time.Sleep(10 * time.Second)
}

func CswLogout() {
	robotgo.MoveMouse(1295, 111) //botao exportar CSV
	time.Sleep(2 * time.Second)
	robotgo.MouseClick("left", false)
	time.Sleep(2 * time.Second)
	robotgo.MoveMouse(1346, 303) //botao exportar CSV
	time.Sleep(2 * time.Second)
	robotgo.MouseClick("left", false)
}

func CswAbrirRotina(rotina string) {
	time.Sleep(1 * time.Second)
	robotgo.MoveMouse(145, 320) //botao execução direta
	time.Sleep(1 * time.Second)
	robotgo.MouseClick("left", false)
	time.Sleep(3 * time.Second)
	robotgo.TypeStr(rotina)
	time.Sleep(1 * time.Second)
	robotgo.MoveMouse(820, 398) //botao execução direta
	time.Sleep(1 * time.Second)
	robotgo.MouseClick("left", false)
}

func CswReloadBrowser() {
	time.Sleep(1 * time.Second)
	robotgo.MoveMouse(87, 81) //botao atualiza navegador (rotina travada)
	time.Sleep(1 * time.Second)
	robotgo.MouseClick("left", false)
	time.Sleep(1 * time.Second)
}
