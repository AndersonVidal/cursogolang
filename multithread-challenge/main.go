package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CEPAPIResponse interface {
	ViaCEPResponse | CEPBrasilAPIResponse
}

type CEPBrasilAPIResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEPResponse struct {
	CEP          string `json:"cep"`
	Street       string `json:"logradouro"`
	Complement   string `json:"complemento"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	IBGE         string `json:"ibge"`
	GIA          string `json:"gia"`
	DDD          string `json:"ddd"`
	SIAFI        string `json:"siafi"`
	Unidade      string `json:"unidade"`
	Estado       string `json:"estado"`
	Regiao       string `json:"regiao"`
}

type ResultWrapper struct {
	Source string      `json:"source"`
	Data   interface{} `json:"data"`
}

func main() {
	r := chi.NewRouter()
    r.Use(middleware.Logger)
	r.Get("/cep/{cep:^[0-9]{8}}", getCEPDetails)
    http.ListenAndServe(":8080", r)
}

func getCEPDetails(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	brasilAPIChan := make(chan *CEPBrasilAPIResponse)
	viaCEPChan := make(chan *ViaCEPResponse)

	go getCEPBrasilAPI(cep, brasilAPIChan)
	go getViaCEPAPI(cep, viaCEPChan)

	var resJson []byte

	select {
	case dataCEPBrasilAPI := <- brasilAPIChan:
		if dataCEPBrasilAPI != nil {
			data := ResultWrapper{
				Source: "BrasilAPI",
				Data: dataCEPBrasilAPI,
			}
			resJson, _ = json.MarshalIndent(data, "", "\t")
			w.Write(resJson)
		}
	case dataViaCEPAPI := <- viaCEPChan:
		if dataViaCEPAPI != nil {
			data := ResultWrapper{
				Source: "ViaCEP",
				Data: dataViaCEPAPI,
			}
			resJson, _ = json.MarshalIndent(data, "", "\t")
			w.Write(resJson)
		}
	case <- time.After(time.Second):
		http.Error(w, "Timeout!", http.StatusGatewayTimeout)
	}
}

func getCEPBrasilAPI(cep string, channel chan *CEPBrasilAPIResponse) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Print(err.Error())
		channel <- nil
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Erro ao fechar resposta: %v", err)
		}
	}()

	var result CEPBrasilAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Print(err.Error())
		channel <- nil
		return
	}

	channel <- &result
}

func getViaCEPAPI(cep string, channel chan *ViaCEPResponse) {
	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Print(err.Error())
		channel <- nil
		return
	}
	
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Erro ao fechar resposta: %v", err)
		}
	}()

	var result ViaCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Print(err.Error())
		channel <- nil
		return
	}

	channel <- &result
}
