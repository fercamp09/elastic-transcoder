package main

import "fmt"
import "net/http"
import "io/ioutil"

func main() {
	url := "http://localhost:3000/files/5891488b6cc87d1ab55690d3"

	file_path := "sender.py"

	resp, err := http.Get(url)
  check(err)
  defer resp.Body.Close()
  out, err := os.Create("files/test1.jpg")
  if err != nil {
    // panic?
  }
  defer out.Close()
  io.Copy(out, resp.Body)


	fmt.Println("Hello mundo")
}
