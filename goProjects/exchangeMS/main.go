package main

import (
	"encoding/json" // Pacote para manipular JSON
	"fmt"           // Pacote para formatação e impressão de texto
	"io/ioutil"     // Pacote para leitura e escrita de arquivos
	"net/http"      // Pacote para comunicação HTTP
)

type TaxasCambio struct {
	Base  string             `json:"base"`  // Moeda base para as taxas de câmbio
	Rates map[string]float64 `json:"rates"` // Mapa com as taxas de câmbio para outras moedas
}

type RequisicaoConversao struct {
	Valor        float64 `json:"valor"`        // Valor a ser convertido
	MoedaOrigem  string  `json:"moedaOrigem"`  // Moeda de origem do valor
	MoedaDestino string  `json:"moedaDestino"` // Moeda de destino para conversão
}

type RespostaConversao struct {
	ValorConvertido float64 `json:"valorConvertido"` // Valor convertido na moeda de destino
}

func main() {
	fmt.Println("Servidor escutando na porta 8080") // Imprime mensagem informando que o servidor está escutando na porta 8080

	// Manipulador para a rota "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Lê o conteúdo do arquivo HTML "conversor.html"
		data, err := ioutil.ReadFile("conversor.html")
		if err != nil {
			fmt.Println(err) // Imprime o erro na leitura do arquivo
			return
		}

		// Envia o conteúdo do arquivo HTML como resposta HTTP
		w.Write(data)
	})

	// Manipulador para a rota "/converter"
	http.HandleFunc("/converter", converter)

	// Inicia o servidor HTTP na porta 8080
	http.ListenAndServe(":8080", nil)
}

func converter(w http.ResponseWriter, r *http.Request) {
	var requisicao RequisicaoConversao // Variável para armazenar os dados da requisição

	// Obtém o corpo da requisição HTTP
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest) // Envia erro 400 (Bad Request)
		return
	}

	// Deserializa o JSON no corpo da requisição para a struct RequisicaoConversao
	err = json.Unmarshal(body, &requisicao)
	if err != nil {
		http.Error(w, "Erro ao deserializar dados da requisição", http.StatusBadRequest) // Envia erro 400 (Bad Request)
		return
	}

	// Obtém as taxas de câmbio
	taxasCambio, err := obterTaxasCambio()
	if err != nil {
		http.Error(w, "Erro ao obter taxas de câmbio", http.StatusInternalServerError) // Envia erro 500 (Internal Server Error)
		return
	}

	// Converte o valor
	valorConvertido := converterValor(requisicao.Valor, requisicao.MoedaOrigem, requisicao.MoedaDestino, taxasCambio)

	// Cria a struct RespostaConversao com o valor convertido
	resposta := RespostaConversao{
		ValorConvertido: valorConvertido,
	}

	// Serializa a struct RespostaConversao em JSON
	jsonData, err := json.Marshal(resposta)
	if err != nil {
		http.Error(w, "Erro ao serializar dados da resposta", http.StatusInternalServerError) // Envia erro 500 (Internal Server Error)
		return
	}

	// Define o cabeçalho Content-Type como application/json
	w.Header().Set("Content-Type", "application/json")

	// Envia o JSON com o valor convertido como resposta HTTP
	w.Write(jsonData)
}

func obterTaxasCambio() (TaxasCambio, error) {
	// URL da API de taxas de câmbio
	url := "https://api.exchangerate-api.com/v4/latest/BRL"

	// Faz uma requisição GET para a URL da API
	resp, err := http.Get(url)
	if err != nil {
		return TaxasCambio{}, err // Retorna struct vazia de TaxasCambio e o erro
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TaxasCambio{}, fmt.Errorf("Erro ao obter taxas de câmbio: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TaxasCambio{}, err
	}

	var taxasCambio TaxasCambio
	err = json.Unmarshal(body, &taxasCambio)
	if err != nil {
		return TaxasCambio{}, err
	}

	return taxasCambio, nil
}

func converterValor(valor float64, moedaOrigem string, moedaDestino string, taxasCambio TaxasCambio) float64 {
	// Validar moedas
	if _, ok := taxasCambio.Rates[moedaOrigem]; !ok {
		panic(fmt.Errorf("Moeda de origem inválida: %s", moedaOrigem))
	}

	if _, ok := taxasCambio.Rates[moedaDestino]; !ok {
		panic(fmt.Errorf("Moeda de destino inválida: %s", moedaDestino))
	}

	// Converter valor
	taxaOrigem := taxasCambio.Rates[moedaOrigem]
	taxaDestino := taxasCambio.Rates[moedaDestino]

	return valor * taxaDestino / taxaOrigem
}
