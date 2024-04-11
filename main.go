package main

import (
	"fmt"

	"reprocess-gui/internal/utils"
)

func main() {
	key := "secret"
	x := &utils.PaginationToken{
		Offset: "1234fsadf234",
		Limit:  50,
	}

	token, err := utils.GeneratePaginationToken(x, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("Token:", token)

	parsed := &utils.PaginationToken{}
	err = utils.ParsePaginationToken(token, key, parsed)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", parsed)
}
