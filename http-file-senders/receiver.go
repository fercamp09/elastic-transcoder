package main

import "fmt"
import "net/http"
import "io"
import "os"

func main() {
	url := "http://localhost:3000/files/58926c92ba2dcd11b4d6c794" //se ingresa el id con el que se descargara el archivo

	//file_path := "sender.py"

	resp, err := http.Get(url)
  //check(err)
  defer resp.Body.Close()
  out, err := os.Create("files/tesasdasdasdast1.jpg")
  if err != nil {
    // panic?
  }
  defer out.Close()
  io.Copy(out, resp.Body)


	fmt.Println("Hello mundo")
}
