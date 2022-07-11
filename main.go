package main

import (
    "fmt"
)

func main() {
    body := get_json("golang")
    fmt.Println(string(body))
}
