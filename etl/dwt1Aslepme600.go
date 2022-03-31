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

type Aslepme600Record struct {
	Asle600TimeStamp   string
	Asle600CodEng      string
	Asle600CodProd     int
	Asle600DescEng     string
	Asle600CodInsumo   int
	Asle600DescInsumo  string
	Asle600Qtde        float64
	Asle600Un          string
	Asle600Status      string
	Asle600CentroCusto string
}

func EltAslepme600() {
	log.SetPrefix("EltAslepme600: ")
	utils.CswLogin()
	time.Sleep(5 * time.Second)
	utils.CswAbrirRotina(config.Aslepme600Nome)

	ExportCsvAslepme600()
	records := ReadAslepme600Csv()
	fmt.Println(len(records))
	ImportToDwtAslepme600(records)
	time.Sleep(5 * time.Second)
	utils.CswLogout()
	time.Sleep(5 * time.Second)

}

func ExportCsvAslepme600() {

	log.Println("Starting RPA for Aslepme600...")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	time.Sleep(120 * time.Second)
	robotgo.MoveMouse(450, 700)
	time.Sleep(2 * time.Second)
	robotgo.MouseClick("left", false)
	time.Sleep(30 * time.Second)
	robotgo.TypeStr(config.Aslepme600Nome)
	time.Sleep(5 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(5 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	utils.CswReloadBrowser()
	time.Sleep(2 * time.Second)
	log.Println("Ending RPA for Aslepme600...")
}

func ReadAslepme600Csv() [][]string {
	filename := config.Aslepme600Nome + ".csv"
	var records [][]string
	log.Println("Reading .csv file...")
	csvfile, err := os.Open(path.Join(config.DownloadFolderString, filename))
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
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
			// record = append(record, " ")
			break
		}
		fmt.Println(len(record))
		fmt.Println(record)
		records = append(records, record)
	}
	records = records[1:] //remove cabe
	// fmt.Println(records)
	log.Println("Finished reading .csv file. Read: ", len(records[1:]))
	return records
}

func ImportToDwtAslepme600(records [][]string) {
	log.Println("Connecting to pgSQL db...")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PsgHost, config.PsgPort, config.PsgUser, config.PsgPassword, config.PsgDbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()


	var r Aslepme600Record
	for i := range records {

		currentTime := time.Now()
		r.Asle600TimeStamp = currentTime.Format("02/01/2006 3:4:5")
		r.Asle600CodEng = utils.ConvertW1252ToUTF8(records[i][0])
		r.Asle600CodProd, _ = strconv.Atoi(records[i][1])
		r.Asle600DescEng = utils.ConvertW1252ToUTF8(records[i][2])
		r.Asle600CodInsumo, _ = strconv.Atoi(records[i][3])
		r.Asle600DescInsumo = utils.ConvertW1252ToUTF8(records[i][4])
		r.Asle600Qtde, _ = strconv.ParseFloat(utils.NormalizeFloat(records[i][5]), 64)
		r.Asle600Un = utils.ConvertW1252ToUTF8(records[i][6])
		r.Asle600Status = utils.ConvertW1252ToUTF8(records[i][9])
		r.Asle600CentroCusto = utils.ConvertW1252ToUTF8(records[i][10])

		sqlStatement := `
			INSERT INTO public.dwt1_csw_aslepme600 (data_carga,cod_eng,cod_prod,desc_eng,cod_insumo,desc_insumo,qtde,unidade,status,centro_custo)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err = db.Exec(sqlStatement, r.Asle600TimeStamp, r.Asle600CodEng, r.Asle600CodProd, r.Asle600DescEng, r.Asle600CodInsumo, r.Asle600DescInsumo, r.Asle600Qtde, r.Asle600Un, r.Asle600Status, r.Asle600CentroCusto)
		if err != nil {
			log.Println(r)
			log.Print(err)
		}
	}
	log.Println("Finnished wrting to db")
}
