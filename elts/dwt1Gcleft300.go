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

type Gcleft300Record struct {
	Gcl300CodPed     int
	Gcl300CodCliente int
	Gcl300DesCliente string
	Gcl300CodRep     string
	Gcl300DescRep    string
	Gcl300Valor      float64
	Gcl300Unidade    string
	Gcl300DataEmi    string
	Gcl300DataPrev   string
	Gcl300DataFat    string
	Gcl300Media      int
	Gcl300Justi      string
	Gcl300Obs        string
}

func EltGcleft300() {
	log.SetPrefix("EltGcleft300: ")
	utils.CswLogin()
	time.Sleep(5 * time.Second)
	utils.CswAbrirRotina(config.Gcleft300Nome)

	ExportCsvGcleft300()
	records := ReadGcleft300Csv()
	fmt.Println(len(records))
	ImportToDwtGcleft300(records)
	time.Sleep(5 * time.Second)
	utils.CswLogout()
	time.Sleep(5 * time.Second)

}

func ExportCsvGcleft300() {

	log.Println("Starting RPA for Gcleft300...")
	time.Sleep(2 * time.Second)
	robotgo.TypeStr("01122020")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("tab")
	time.Sleep(2 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(60 * time.Second)
	robotgo.TypeStr(config.Gcleft300Nome)
	time.Sleep(5 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(5 * time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	utils.CswReloadBrowser()
	time.Sleep(2 * time.Second)
	log.Println("Ending RPA for Gcleft300...")
}

func ReadGcleft300Csv() [][]string {
	filename := config.Gcleft300Nome + ".csv"
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
			// log.Println(err)
			record = append(record, " ")
			// break
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

func ImportToDwtGcleft300(records [][]string) {
	log.Println("Connecting to pgSQL db...")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PsgHost, config.PsgPort, config.PsgUser, config.PsgPassword, config.PsgDbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	//Clear Table
	sqlStatement := `
	DELETE FROM public.dwt1_csw_gcleft300` //103247 ultimo Item empresa 1

	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Println(err)
	}

	var r Gcleft300Record
	for i := range records {

		var t time.Time
		r.Gcl300CodPed, _ = strconv.Atoi(records[i][0])
		r.Gcl300CodCliente, _ = strconv.Atoi(records[i][1])
		r.Gcl300DesCliente = utils.ConvertW1252ToUTF8(records[i][2])
		r.Gcl300CodRep = utils.ConvertW1252ToUTF8(records[i][3])
		r.Gcl300DescRep = utils.ConvertW1252ToUTF8(records[i][4])
		r.Gcl300Valor, _ = strconv.ParseFloat(utils.NormalizeFloat(records[i][5]), 64)
		t, _ = time.Parse("02/01/2006", records[i][6])
		r.Gcl300DataEmi = t.Format("01/02/2006")
		t, _ = time.Parse("02/01/2006", records[i][7])
		r.Gcl300DataPrev = t.Format("01/02/2006")
		t, _ = time.Parse("02/01/2006", records[i][8])
		r.Gcl300DataFat = t.Format("01/02/2006")
		r.Gcl300Media, _ = strconv.Atoi(records[i][9])
		r.Gcl300Justi = utils.ConvertW1252ToUTF8(records[i][10])
		r.Gcl300Obs = utils.ConvertW1252ToUTF8(records[i][11])

		sqlStatement := `
			INSERT INTO public.dwt1_csw_gcleft300 (cod_ped,cod_cliente,desc_client,cod_rep,desc_rep,valor_ped,data_emi,data_prev,data_fat,media,justificativa,obs)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = db.Exec(sqlStatement, r.Gcl300CodPed, r.Gcl300CodCliente, r.Gcl300DesCliente, r.Gcl300CodRep, r.Gcl300DescRep, r.Gcl300Valor, r.Gcl300DataEmi, r.Gcl300DataPrev, r.Gcl300DataFat, r.Gcl300Media, r.Gcl300Justi, r.Gcl300Obs)
		if err != nil {
			log.Println(r)
			log.Print(err)
		}
	}
	log.Println("Finnished wrting to db")
}
