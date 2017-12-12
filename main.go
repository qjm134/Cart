package main

import (
	"net/http"
	"cart/controllers"
	"cart/models/redis"
	"cart/conf"
)

func main() {
	redis.Init(conf.DataSource)

	http.HandleFunc("/addCart", controllers.AddCart)
	http.ListenAndServe(":8018", nil)
}

