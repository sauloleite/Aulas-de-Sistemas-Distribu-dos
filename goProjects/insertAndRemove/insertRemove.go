package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Definição da estrutura de dados para um usuário
type User struct {
	Username string `json:"username"` // Campo para o nome de usuário
	Password string `json:"password"` // Campo para a senha do usuário
}

func main() {
	// Carrega os usuários do arquivo
	fmt.Println("Servidor rodando localmente na porta 8080")
	users, err := loadUsers()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Manipuladores para diferentes rotas HTTP

	// Manipulador para a rota de login

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// Extrai os valores do formulário HTTP para username e password
		username := r.FormValue("username")
		// Tente remover o usuário
		err := removeUser(username)
		if err != nil {
			// Ser houver erro, retorna messagem de erro
			fmt.Fprintf(w, "Nome de usuário não existe")
			return
		}
		fmt.Fprintf(w, "Usuário removido com sucesso!")
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		// Extrai os valores do formulário HTTP para username e password
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Verifica se o usuário e a senha fornecidos são válidos
		if !isValidUsername(users, username) {
			// Retorna status de não autorizado se não forem válidos
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Username já em uso.")
			return
		}
		newUser := User{Username: username, Password: password}
		users = append(users, newUser)
		err := saveUsers(users)
		if err != nil {
			fmt.Println(err)
			return

		}
		// Se o username for válido, inserir
		fmt.Fprintf(w, "Cadastrado com sucesso!")
	})

	// Manipulador para a rota padrão (login)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Lê o conteúdo do arquivo HTML de login e o envia como resposta HTTP
		data, err := ioutil.ReadFile("cadastro.html")
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(data)
	})
	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		// Lê o conteúdo do arquivo HTML
		data, err := ioutil.ReadFile("remove.html")
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(data)
	})

	// Inicia o servidor HTTP na porta 8080
	http.ListenAndServe(":8080", nil)
}

// Função para carregar os usuários do arquivo JSON
func loadUsers() ([]User, error) {
	// Lê os dados do arquivo de usuários
	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		return nil, err
	}

	// Decodifica os dados JSON em uma lista de usuários
	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Função para verificar se um usuário é válido
func isValidUsername(users []User, username string) bool {
	// Verifica se existe um usuário com o nome de usuário inserido
	for _, user := range users {
		if user.Username == username {
			return false
		}
	}

	return true
}

func saveUsers(users []User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("users.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func removeUser(username string) error {
	users, err := loadUsers()
	if err != nil {
		return err
	}
	var found bool
	for i, user := range users {
		if user.Username == username {
			// Remove o usuário da lista
			users = append(users[:i], users[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("Usuário não encontrado")
	}
	//Salva a lista atualizada de usuários no arquivo
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("users.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
