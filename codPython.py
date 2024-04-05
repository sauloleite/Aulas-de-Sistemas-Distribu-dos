import requests
base_url = "http://localhost:8080"
def login(username, password):
    url = base_url + "/login"
    data = {"username": username, "password": password}
    response = requests.post(url, data=data)
    return response.text

#Exemplo de uso
if __name__ == "__main__":
    print("Login:", login("admin2","123"))