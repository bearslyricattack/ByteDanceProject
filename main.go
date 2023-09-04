package main

import (
	api "douyin/kitex_gen/api/user"
	"log"
)

func main() {
	svr := api.NewServer(new(UserImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
