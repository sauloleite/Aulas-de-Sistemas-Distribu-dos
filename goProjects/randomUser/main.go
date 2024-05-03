package main

import (
	// Importa pacotes necessários
	// Para trabalhar com dados JSON
	"encoding/json"
	"fmt"       // Para imprimir na saída padrão
	"io/ioutil" // Para leitura e escrita de arquivos
	"log"       // Para registro de mensagens
	"net/http"  // Para criar servidor HTTP
)

// Define estrutura para armazenar dados da resposta da API
type RandomUserResponse struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			City    string `json:"city"`
			State   string `json:"state"`
			Country string `json:"country"`
		} `json:"location"`
		Email string `json:"email"`
	} `json:"results"`
}

func main() {
	// Imprime mensagem de inicialização do servidor
	fmt.Println("Servidor escutando na porta 8080")

	// Define o diretório atual para servir arquivos estáticos
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Cria endpoint para obter usuário aleatório
	http.HandleFunc("/getRandomUser", getRandomUser)

	// Inicia o servidor e registra erro caso haja falha
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func getRandomUser(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://randomuser.me/api/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// Lê o corpo da resposta da API
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var randomUserResp RandomUserResponse
	err = json.Unmarshal(body, &randomUserResp)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomUserResp)
	err = ioutil.WriteFile("randomuser.json", body, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
