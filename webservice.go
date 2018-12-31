package main
import (
	"fmt"
	"net/http"
 	"io/ioutil"
	"encoding/xml"
	"io"
	_"github.com/lib/pq"
 	"database/sql"
	"strconv"
)


const (
        DB_USER     = "???"
        DB_PASSWORD = "???"
        DB_NAME     = "???"
	DB_HOST     = "localhost"
	CAMINHO = "/arquivamento/???/"
	CAMINHO_ERROR = "/arquivamento/error/"
    )


type Exames struct {
	CodigoClinica string  `xml:"codigoClinica"`
	Codigo string `xml:"codigo"`
	Descricao string `xml:"descricao"`
	CrmMedicoRealizante string `xml:"crmMedicoRealizante"`
	CrmUfMedicoRealizante string `xml:"crmUfMedicoRealizante"`
	NomeRealizante string `xml:"nomeRealizante"`
	Datahora string `xml:"datahora"`
	Arquivo string `xml:"arquivo"`
}



type Resultados struct {
	AccessionNumber string  `xml:"accessionNumber"`
 	PacienteId string  `xml:"pacienteid"`
	Nome string  `xml:"nome"`
	DataNascimento string  `xml:"dataNascimento"`
	Sexo string  `xml:"sexo"`
 	Xml string `xml:",innerxml"`
	Exame []Exames `xml:"exames>exame"`
}

type Retorno struct {
	XMLName xml.Name `xml:"RetornaProcedimentoResponse"`
	RetornaProcedimentoResponse bool `xml:RetornaProcedimento"`
	Xml string `xml:",innerxml"`
}


func gravar_arquivo(post []byte,arquivo string) {

	// output, err := xml.Marshal(&post)
	err := ioutil.WriteFile(CAMINHO + arquivo, post, 0644)
	if err != nil {
		fmt.Println("Error writing XML to file:", err)
//		return
	}
}

func gravar_banco(post Resultados) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		DB_USER, DB_PASSWORD, DB_NAME,DB_HOST)
        db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println("Erro ao conectar banco:",err)
		return
	}

	cd_exame,_ := strconv.Atoi(post.Exame[0].Codigo)
	_,err = db.Exec(`insert into aux_laudos_pacs(accession_number,
                        laudo_datahora,
                         paciente_id,
                         laudo_conteudo,
                         medico_responsavel_conselho_numero,
                         medico_responsavel_conselho_uf,
                         medico_responsavel_nome,
                         ds_paciente,
                         ds_procedimento,
                         cd_exame)
                        values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`,
		 post.AccessionNumber,
		 post.Exame[0].Datahora,
		 post.PacienteId,
		 post.Exame[0].Arquivo,
		 post.Exame[0].CrmMedicoRealizante,
		 post.Exame[0].CrmUfMedicoRealizante,
		 post.Exame[0].NomeRealizante,
	       	post.Nome,
		post.Exame[0].Descricao,
		cd_exame)

	if err != nil {
		panic(err)

	}

}


func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func integracao(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "World!")
	// ffmt.Println(r)
	// body, _ := ioutil.ReadAll(r.Body)
	// fmt.Println(string(body))
	decoder := xml.NewDecoder(r.Body)
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error decoding XML into tokens:", err)
			// return
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "Resultados" {
				var result Resultados
				decoder.DecodeElement(&result, &se)
				fmt.Println(result.Xml)
				gravar_arquivo([]byte(result.Xml),result.AccessionNumber)
				gravar_banco(result)

			}
		}
	}

//	var retorno Retorno;
	//retorno.RetornaProcedimentoResponse = true
	retorno := Retorno{RetornaProcedimentoResponse : true,}
	output,err  := xml.Marshal(&retorno)
	if err != nil {
		fmt.Println("Error marshalling to XML:", err)
	}

	err = ioutil.WriteFile("post-retorno.xml", output, 0644)
	if err != nil {
		fmt.Println("Error writing XML to file:", err)

	}
	fmt.Fprint(w,string(output))
	//xml.Unmarshal(body, &result)
//	fmt.Println(result)

}


func main() {
	/* server := http.Server{
		Addr: "127.0.0.1:8081",
	}
*/
	http.HandleFunc("/integracao/integracao.php", hello)
	http.HandleFunc("?????", integracao)
	http.ListenAndServe(":80", nil)
	// server.ListenAndServe()
}
