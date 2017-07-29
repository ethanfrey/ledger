package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethanfrey/goc/ledger"
)

func main() {
	ledger, err := ledger.FindLedger()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}

	data := strings.Join(os.Args[1:], " ")
	fmt.Println("Sending", data)
	fmt.Println("")

	resp, err := ledger.Exchange([]byte(data), 100)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}
	fmt.Printf("Response: %X\n", resp)
}
