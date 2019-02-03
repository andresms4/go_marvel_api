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

//APIresponse estructura primer nivel
type APIresponse struct {
	Data Data `json:"data"`
}

//Data estructura segundo nivel
type Data struct {
	Results []Results `json:"results"`
}

//Results estructura segundo nivel
type Results struct {
	ID          int       `json:"id"`
	Nombre      string    `json:"name"`
	Descripcion string    `json:"description"`
	Modificado  string    `json:"modified"`
	Comics      Comics    `json:"comics"`
	Series      Series    `json:"series"`
	Historias   Historias `json:"stories"`
	Eventos     Eventos   `json:"events"`
}

//Comics estructura tercer nivel
type Comics struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Series estructura tercer nivel
type Series struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Historias estructura tercer nivel
type Historias struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Eventos estructura tercer nivel
type Eventos struct {
	Disponible int     `json:"available"`
	Items      []Items `json:"items"`
}

//Items estructura cuarto nivel
type Items struct {
	Nombre string `json:"name"`
}

func main() {

	ts := time.Now().Format("20060102150405")                                                               //se obtiene el timestamp y se formatea
	publicKey, privateKey := "4de5c74cfbb660b526a8ce118904822d", "7dc5edeb782cdb834b0c66a54b68dff0eb6add8e" //llaves publica y privada
	hash := md5.New()
	io.WriteString(hash, ts+privateKey+publicKey) //se obtiene el md5
	hashhex := hex.EncodeToString(hash.Sum(nil))  //Formato hexadecimal para usar en la URL

	fmt.Println("Este programa consume datos de la API de Marvel")
	fmt.Println("Seleccione la opción deseada:\n1. Busqueda\n2. Listar")

	reader := bufio.NewReader(os.Stdin) //se usa bufio en lugar de fmt por existir en error en el momento de iniciar
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	switch char {
	case 49: //ascii para 1 y 2
		fmt.Println("¿Que personaje desea encontrar?:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() //se recibe el nombre a buscar
		text := scanner.Text()
		nombre := url.QueryEscape(text)                                                                                                      //si existe un espacio se da formato para evitar error en la solicitud
		url := fmt.Sprintf("http://gateway.marvel.com/v1/public/characters?name=%s&ts=%s&apikey=%s&hash=%s", nombre, ts, publicKey, hashhex) //URL lista para enviar la solicitud
		contents, err := GetRequest(url)                                                                                                     // se envia la solicitud
		if err != nil {
			fmt.Println(err)
		}
		var s APIresponse
		err = json.Unmarshal(contents, &s) //recibe los datos y almacena en el formato de la estructura
		if err != nil {
			fmt.Println(err)
		}
		b, err := json.MarshalIndent(s.Data.Results, "", "  ") //Se agrega indentacion y mejor presentacion a los datos
		if err != nil {
			fmt.Println(err)
		}
		os.Stdout.Write(b) //se imprime

	case 50: //ascii para 2

		url := fmt.Sprintf("http://gateway.marvel.com/v1/public/characters?orderBy=name&ts=%s&apikey=%s&hash=%s", ts, publicKey, hashhex) //url para la solicitud
		contents, err := GetRequest(url)                                                                                                  //envia la solicitud
		if err != nil {
			fmt.Println(err)
		}
		var s APIresponse
		err = json.Unmarshal(contents, &s) //recibe los datos y almacena en el formato de la estructura
		if err != nil {
			fmt.Println(err)
		}
		b, err := json.MarshalIndent(s.Data.Results, "", "  ") //Se agrega indentacion y mejor presentacion a los datos
		if err != nil {
			fmt.Println(err)
		}
		os.Stdout.Write(b) //se imprime

	default:
		fmt.Println("La opción seleccionada no es válida")
	}

}

//GetRequest peticion
func GetRequest(url string) ([]byte, error) {
	resp, err := http.Get(string(url)) //se envia la solicitud con la URL recibida
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()                    //Se cierra la conexion
	contents, err := ioutil.ReadAll(resp.Body) //lee los datos
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("Solicitando...")
	return []byte(contents), err //devuelve el contenido o el error

}
