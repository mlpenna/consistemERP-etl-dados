package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func ConvertW1252ToUTF8(s string) string {
	rInUTF8 := transform.NewReader(strings.NewReader(s), charmap.Windows1252.NewDecoder())
	decBytes, _ := ioutil.ReadAll(rInUTF8)
	decS := string(decBytes)
	return decS
}

func NormalizeFloat(old string) string {
	s := strings.Replace(old, ".", "", -1)
	return strings.Replace(s, ",", ".", -1)
}

func CheckMousePos() {
	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(robotgo.GetMousePos())
	}

}
