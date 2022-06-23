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

type Gcleese610Record struct {
	GclData          string
	GclCodProduto    int
	GclDesProduto    string
	GclCodNatureza   int
	GclDescNatureza  string
	GclQtdeConsumida float64
	GclUnidade       string
}

func EltGclese610() {
	utils.CswLogin()
	ExportCsvGcleese610()
	records := ReadGcleese610Csv()
	ImportToDwt(records)
	utils.CswLogout()
}

func ReadGcleese610Csv() [][]string {
	filename := config.Gcleese610Nome + ".csv"
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
			break
		}
		records = append(records, record)
	}
	records = records[1:] //remove cabe
	log.Println("Finished reading .csv file. Read: ", len(records[1:]))
	return records
}

func ImportToDwt(records [][]string) {
	log.Println("Connecting to pgSQL db...")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PsgHost, config.PsgPort, config.PsgUser, config.PsgPassword, config.PsgDbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	
	var r Gcleese610Record
	for i := range records {
		var t time.Time
		t, _ = time.Parse("02/01/2006", records[i][0])
		r.GclData = t.Format("01/02/2006")
		r.GclCodProduto, _ = strconv.Atoi(records[i][1])
		r.GclDesProduto = utils.ConvertW1252ToUTF8(records[i][2])
		r.GclCodNatureza, _ = strconv.Atoi(records[i][3])
		r.GclDescNatureza = utils.ConvertW1252ToUTF8(records[i][4])
		r.GclQtdeConsumida, _ = strconv.ParseFloat(utils.NormalizeFloat(records[i][5]), 64)
		r.GclUnidade = records[i][6]
		sqlStatement := `
		INSERT INTO public.dwt1_consumo_mp_emb (data,cod_produto,desc_produto,cod_natureza,desc_natureza,qtde_consumo,unidade_consumo)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = db.Exec(sqlStatement, r.GclData, r.GclCodProduto, r.GclDesProduto, r.GclCodNatureza, r.GclDescNatureza, r.GclQtdeConsumida, r.GclUnidade)
		if err != nil {
			log.Println(r.GclDesProduto + " " + r.GclDescNatureza)
			log.Print(err)
		}
	}
	log.Println("Finnished wrting to db")
}
