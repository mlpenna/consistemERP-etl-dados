package elts

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/accurati-bi/csw-rpa-elt/config"
	"github.com/accurati-bi/csw-rpa-elt/utils"

	"github.com/go-vgo/robotgo"
	_ "github.com/lib/pq"
)

type Ccese320Record struct {
	CC320IdFornecedor   int
	CC320DescFornecedor string
	CC320DataEntrada    string
	CC320DataEmissao    string
	CC320DataPrevista   string
	CC320NF             int
	CC320NFItem         int
	CC320Pedido         int
	CC320PedidoItem     int
	CC320IdMaterial     int
	CC320DescMaterial   string
	CC320QtdeEntrada    float64
}

func EltCcese320() {
	log.SetPrefix("EltCcese320: ")
	utils.CswLogin()
	time.Sleep(5 * time.Second)
	utils.CswAbrirRotina(config.Ccese320Nome)

	ExportCsvCcese320()
	records := ReadCcese320Csv()
	// fmt.Print(records)
	ImportToDwtCcese320(records)
	time.Sleep(4 * time.Second)
	utils.CswLogout()
	time.Sleep(2 * time.Second)
}

func ExportCsvCcese320() {

	time.Sleep(2 * time.Second)
	robotgo.TypeStr("01122020")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.TypeStr("2.01")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	robotgo.TypeStr("2.02")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("1")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("down")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("down")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	robotgo.MoveMouse(640, 434)
	time.Sleep(2 * time.Second)
	robotgo.MouseClick("left", false)
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(8 * time.Second)
	robotgo.TypeStr(config.Ccese320Nome)
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	utils.CswReloadBrowser()
}

func ReadCcese320Csv() [][]string {
	filename := config.Ccese320Nome + ".csv"
	var records [][]string
	csvfile, err := os.Open(path.Join(config.DownloadFolderString, filename))
	if err != nil {
		log.Println("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)
	r.Comma = ';'
	r.LazyQuotes = true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		records = append(records, record)
	}
	records = records[1:] //remove cabe
	// log.Println(len(records))
	return records
}

func ImportToDwtCcese320(records [][]string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PsgHost, config.PsgPort, config.PsgUser, config.PsgPassword, config.PsgDbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	sqlStatement := `
DELETE FROM public.dwt1_csw_ccese320 WHERE id > 88999`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Print(err)
	}

	var r Ccese320Record
	for i := range records {
		var t time.Time

		r.CC320IdFornecedor, err = strconv.Atoi(records[i][0])
		if err != nil {
			log.Print(err)
		}

		r.CC320DescFornecedor = utils.ConvertW1252ToUTF8(records[i][1])

		t, err = time.Parse("02/01/06", records[i][2])
		if err != nil {
			log.Print(err)
		}
		r.CC320DataEntrada = t.Format("01/02/2006")

		t, err = time.Parse("02/01/06", records[i][3])
		if err != nil {
			log.Print(err)
		}
		r.CC320DataEmissao = t.Format("01/02/2006")

		t, err = time.Parse("02/01/06", records[i][5])
		if err != nil {
			log.Print(err)
		}
		r.CC320DataPrevista = t.Format("01/02/2006")

		r.CC320NF, _ = strconv.Atoi(records[i][6])
		r.CC320NFItem, _ = strconv.Atoi(records[i][7])
		r.CC320Pedido, _ = strconv.Atoi(records[i][8])
		r.CC320PedidoItem, _ = strconv.Atoi(records[i][9])
		r.CC320IdMaterial, _ = strconv.Atoi(records[i][12])
		r.CC320DescMaterial = utils.ConvertW1252ToUTF8(records[i][13])
		r.CC320QtdeEntrada, err = strconv.ParseFloat(utils.NormalizeFloat(records[i][14]), 64)
		if err != nil {
			log.Print(err)
		}
		// fmt.Println(r)
		sqlStatement := `
		INSERT INTO public.dwt1_csw_ccese320 (id_fornecedor,desc_fornecedor,data_entrada,data_emissao,data_prevista,nf,nf_item,
		pedido,pedido_item,cod_material,desc_material,qtde_entrada)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = db.Exec(sqlStatement, r.CC320IdFornecedor, r.CC320DescFornecedor, r.CC320DataEntrada, r.CC320DataEmissao, r.CC320DataPrevista,
			r.CC320NF, r.CC320NFItem, r.CC320Pedido, r.CC320PedidoItem, r.CC320IdMaterial, r.CC320DescMaterial, r.CC320QtdeEntrada)
		if err != nil {
			log.Println(r)
			log.Println(err)
		}
	}
}
