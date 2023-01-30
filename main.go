package main

import (
	"fmt"
	"github.com/herzorf/filestroe-server/route"
)

func main() {
	router := route.Router()
	err := router.Run(":8080")
	if err != nil {
		fmt.Println("gin run err", err)
	}
}
