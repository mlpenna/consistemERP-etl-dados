package elts

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/accurati-bi/csw-rpa-elt/config"
	_ "github.com/lib/pq"
	_ "github.com/nakagami/firebirdsql"
)

type sigmaPecas struct {
	pcCod      string
	pcDesc     string
	pcEstatual float64
	pcUni      string
	pcEstmax   float64
	pcEstmin   float64
}

func EltSigmaEstoque() {
	updatePecasTable(readPecasTable())
	// readPecasTable()
}

func updatePecasTable(pecasArray []sigmaPecas) {

	//Connect
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PsgHost, config.PsgPort, config.PsgUser, config.PsgPassword, config.PsgSigmaDbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// sqlStatement := `
	// DELETE FROM public.dwt1_sigma_estoque`

	// _, err = db.Exec(sqlStatement)
	// if err != nil {
	// 	log.Print(err)
	// }

	currentTime := time.Now()
	// fmt.Println(currentTime.Format("02/01/2006 3:4:5"))
	//Insert new TAG values
	for _, v := range pecasArray {
		sqlStatement := `
		INSERT INTO public.dwt1_sigma_estoque(data_carga,pc_cod,pc_desc,pc_estoque,pc_un,pc_max,pc_min)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err = db.Exec(sqlStatement, currentTime.Format("02/01/2006 3:4:5"), v.pcCod, v.pcDesc, v.pcEstatual, v.pcUni, v.pcEstmax, v.pcEstmin)
		if err != nil {
			log.Print(err)
		}
	}

	// return
}

func readPecasTable() []sigmaPecas {

	var pecasArray []sigmaPecas
	var currentPeca sigmaPecas

	conn, err := sql.Open("firebirdsql", config.SigmaConnectionString)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	//PC_CODIGO,PC_DESCRI,PC_ESTOATU,UNIDADE,PC_ESTOMAX,PC_ESTOMIN
	rows, err := conn.Query("SELECT PC_CODIGO,PC_DESCRIC,PC_ESTOATU,UNIDADE,PC_ESTOMAX,PC_ESTOMIN FROM PECAS")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		err := rows.Scan(&currentPeca.pcCod, &currentPeca.pcDesc, &currentPeca.pcEstatual, &currentPeca.pcUni, &currentPeca.pcEstmax, &currentPeca.pcEstmin)
		if err != nil {
			log.Print(err)
		}
		// fmt.Println(currentPeca)
		pecasArray = append(pecasArray, currentPeca)
	}
	// fmt.Println(len(pecasArray))
	return pecasArray
}
