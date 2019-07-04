package main

import (
	"fmt"
	"wcore/db"
	"wcore/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	// 初始化引擎
	engine := gin.Default()
	handler.InitHandler(engine)
	// 绑定端口，然后启动应用
	err := engine.Run(":8080")
	if err != nil {
		fmt.Printf("ListenAndServe err:%s", err.Error())
	}
}
