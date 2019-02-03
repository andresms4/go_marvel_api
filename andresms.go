package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

//APIresponse primer nivel
type APIresponse struct {
	Data Data `json:"data"`
}

//Data segundo nivel
type Data struct {
	Results []Results `json:"results"`
}

//Results segundo nivel
type Results struct {
	ID          int       `json:"id"`
	Nombre      string    `json:"name"` //llenar todos los datos
	Descripcion string    `json:"description"`
	Modificado  string    `json:"modified"`
	Comics      Comics    `json:"comics"`
	Series      Series    `json:"series"`
	Historias   Historias `json:"stories"`
	Eventos     Eventos   `json:"events"`
}

//Comics tercer nivel
type Comics struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Series tercer nivel
type Series struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Historias tercer nivel
type Historias struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Eventos tercer nivel
type Eventos struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Items cuarto nivel
type Items struct {
	Nombre string `json:"name"`
}

func main() {

	ts := time.Now().Format("20060102150405") //"20060102150405" constantes de Go para el formato del timestamp
	publicKey, privateKey := "4de5c74cfbb660b526a8ce118904822d", "7dc5edeb782cdb834b0c66a54b68dff0eb6add8e"
	hash := md5.New()
	io.WriteString(hash, ts+privateKey+publicKey)
	hashhex := hex.EncodeToString(hash.Sum(nil))

	fmt.Println("Este programa consume datos de la API de Marvel")
	fmt.Println("Seleccione la opción deseada:\n1. Busqueda\n2. Listar")
	var input uint8
	fmt.Scanf("%o\n", &input)

	switch input {
	case 1:
		fmt.Println("¿Que personaje desea encontrar?:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		nombre := url.QueryEscape(text)
		url := fmt.Sprintf("http://gateway.marvel.com/v1/public/characters?name=%s&ts=%s&apikey=%s&hash=%s", nombre, ts, publicKey, hashhex)
		contents, err := GetRequest(url)
		if err != nil {
			fmt.Println(err)
		}
		var s APIresponse
		err = json.Unmarshal(contents, &s)
		if err != nil {
			fmt.Println(err)
		}
		b, err := json.MarshalIndent(s.Data.Results, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		os.Stdout.Write(b)

	case 2:

		url := fmt.Sprintf("http://gateway.marvel.com/v1/public/characters?ts=%s&apikey=%s&hash=%s", ts, publicKey, hashhex)
		contents, err := GetRequest(url)
		if err != nil {
			fmt.Println(err)
		}
		var s APIresponse
		err = json.Unmarshal(contents, &s)
		if err != nil {
			fmt.Println(err)
		}
		b, err := json.MarshalIndent(s.Data.Results, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		os.Stdout.Write(b)

	default:
		fmt.Println("La opción seleccionada no es válida")
	}

}

//GetRequest peticion
func GetRequest(url string) ([]byte, error) {
	resp, err := http.Get(string(url))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("Solicitando...")
	return []byte(contents), err

}
