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
	"strings"
	"time"

	"github.com/accurati-bi/csw-rpa-elt/config"
	"github.com/accurati-bi/csw-rpa-elt/utils"

	"github.com/go-vgo/robotgo"
	_ "github.com/lib/pq"
)

type Ccesu665Record struct {
	Cc665CodPedido       int
	Cc665DataCadastro    string
	Cc665ItemPedido      int
	Cc665CodProduto      int
	Cc665DescProduto     string
	Cc665SituacaoItem    string
	Cc665QtdeEstoque     float64
	Cc665UnEstoque       string
	Cc665QtdeCompra      float64
	Cc665UnCompra        string
	Cc665QtdeAtendida    float64
	Cc665QtdeCancelada   float64
	Cc665QtdeSaldo       float64
	Cc665CodFornecedor   int
	Cc665DescFornecedor  string
	Cc665SituacaoPedido  string
	Cc665DataLiberacao   string
	Cc665DataEnvio       string
	Cc665DataAceite      string
	Cc665DataPrevInici   string
	Cc665DataPreviAtual  string
	Cc665Comprador       int
	Cc665NomeComprador   string
	Cc665DataRecebimento string
}

func ReadCcesu665Csv() [][]string {
	filename := config.Ccesu665Nome + ".csv"
	var records [][]string
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
	records = records[1:]
	return records
}

func ImportToDwtCcesu665(records [][]string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PsgHost, config.PsgPort, config.PsgUser, config.PsgPassword, config.PsgDbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	var r Ccesu665Record
	for i := range records {
		var t time.Time

		r.Cc665CodPedido, err = strconv.Atoi(records[i][0])
		if err != nil {
			log.Print(err)
		}

		if idx := strings.IndexByte(records[i][1], ' '); idx >= 0 {
			records[i][1] = records[i][1][:idx]
		}
		t, err = time.Parse("02/01/2006", records[i][1])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DataCadastro = t.Format("01/02/2006")
		r.Cc665ItemPedido, _ = strconv.Atoi(records[i][2])
		r.Cc665CodProduto, _ = strconv.Atoi(records[i][3])
		r.Cc665DescProduto = utils.ConvertW1252ToUTF8(records[i][4])
		r.Cc665SituacaoItem = utils.ConvertW1252ToUTF8(records[i][5])
		r.Cc665QtdeEstoque, err = strconv.ParseFloat(utils.NormalizeFloat(records[i][6]), 64)
		if err != nil {
			log.Print(err)
		}
		r.Cc665UnEstoque = utils.ConvertW1252ToUTF8(records[i][7])
		r.Cc665QtdeCompra, err = strconv.ParseFloat(utils.NormalizeFloat(records[i][8]), 64)
		if err != nil {
			log.Print(err)
		}
		r.Cc665UnCompra = records[i][9]
		r.Cc665QtdeAtendida, err = strconv.ParseFloat(utils.NormalizeFloat(records[i][10]), 64)
		if err != nil {
			log.Print(err)
		}
		r.Cc665QtdeCancelada, err = strconv.ParseFloat(utils.NormalizeFloat(records[i][11]), 64)
		if err != nil {
			log.Print(err)
		}
		r.Cc665QtdeSaldo, err = strconv.ParseFloat(utils.NormalizeFloat(records[i][12]), 64)
		if err != nil {
			log.Print(err)
		}
		r.Cc665CodFornecedor, err = strconv.Atoi(records[i][30])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DescFornecedor = utils.ConvertW1252ToUTF8(records[i][31])
		r.Cc665SituacaoPedido = utils.ConvertW1252ToUTF8(records[i][32])

		if idx := strings.IndexByte(records[i][33], ' '); idx >= 0 {
			records[i][33] = records[i][33][:idx]
		}
		t, err = time.Parse("02/01/2006", records[i][33])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DataLiberacao = t.Format("01/02/2006")

		if idx := strings.IndexByte(records[i][34], ' '); idx >= 0 {
			records[i][34] = records[i][34][:idx]
		}
		t, err = time.Parse("02/01/2006", records[i][34])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DataEnvio = t.Format("01/02/2006")

		if idx := strings.IndexByte(records[i][35], ' '); idx >= 0 {
			records[i][35] = records[i][35][:idx]
		}
		t, err = time.Parse("02/01/2006", records[i][35])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DataAceite = t.Format("01/02/2006")

		if idx := strings.IndexByte(records[i][36], ' '); idx >= 0 {
			records[i][36] = records[i][36][:idx]
		}
		t, err = time.Parse("02/01/2006", records[i][36])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DataPrevInici = t.Format("01/02/2006")

		if idx := strings.IndexByte(records[i][37], ' '); idx >= 0 {
			records[i][37] = records[i][37][:idx]
		}
		t, err = time.Parse("02/01/2006", records[i][37])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DataPreviAtual = t.Format("01/02/2006")

		r.Cc665Comprador, err = strconv.Atoi(records[i][41])
		if err != nil {
			log.Print(err)
		}
		r.Cc665NomeComprador = utils.ConvertW1252ToUTF8(records[i][42])

		if idx := strings.IndexByte(records[i][53], ' '); idx >= 0 {
			records[i][53] = records[i][53][:idx]
		}
		t, err = time.Parse("02/01/2006", records[i][53])
		if err != nil {
			log.Print(err)
		}
		r.Cc665DataRecebimento = t.Format("01/02/2006")
		// fmt.Println(r)

		sqlStatement := `
		INSERT INTO public.dwt1_csw_ccesu665 (cod_pedido,data_cadastro,item_pedido,cod_produto,desc_produto,situacao_item,
			qtde_estoque,un_estoque,qtde_compra,un_compra,qtde_compra_atendida,qtde_compra_cancelada,
			qtde_compra_saldo,cod_fornecedor,desc_fornecedor,situacao_pedido,data_liberacao,data_envio,
			data_aceite,previsao_inicial,previsao_atual,comprador,nome_comprador,data_recebimento)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
				$18, $19, $20, $21, $22, $23, $24)`
		_, err = db.Exec(sqlStatement, r.Cc665CodPedido, r.Cc665DataCadastro,
			r.Cc665ItemPedido, r.Cc665CodProduto, r.Cc665DescProduto, r.Cc665SituacaoItem,
			r.Cc665QtdeEstoque, r.Cc665UnEstoque, r.Cc665QtdeCompra, r.Cc665UnCompra, r.Cc665QtdeAtendida, r.Cc665QtdeCancelada,
			r.Cc665QtdeSaldo, r.Cc665CodFornecedor, r.Cc665DescFornecedor, r.Cc665SituacaoPedido, r.Cc665DataLiberacao,
			r.Cc665DataEnvio, r.Cc665DataAceite, r.Cc665DataPrevInici, r.Cc665DataPreviAtual, r.Cc665Comprador,
			r.Cc665NomeComprador, r.Cc665DataRecebimento)
		if err != nil {
			fmt.Println(r)
			fmt.Println(err)
		}
	}
}
