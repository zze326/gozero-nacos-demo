package swagger

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/**
 * @Author: zze
 * @Date: 2022/5/30 15:21
 * @Desc: 初始化 Swagger 路由
 */

func RegisterSwagger(server *rest.Server) {
	server.AddRoutes(
		InitRoutes(),
	)
}

func InitRoutes() []rest.Route {
	swaggerJsonFilePath := os.Getenv("SWAGGER_JSON_FILE_PATH")
	if len(swaggerJsonFilePath) == 0 {
		swaggerJsonFilePath = "swagger.json"
	}
	swaggerFile, err := os.Open(swaggerJsonFilePath)
	if err != nil {
		log.Println(err)
	}
	defer swaggerFile.Close()
	SwaggerByte, err := ioutil.ReadAll(swaggerFile)
	if err != nil {
		log.Println(err)
	}

	env := os.Getenv("ENV")

	return []rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/swagger",
			Handler: Doc("/swagger", env),
		},
		{
			Method: http.MethodGet,
			Path:   "/swagger-json",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")
				_, err := writer.Write(SwaggerByte)
				if err != nil {
					httpx.Error(writer, err)
				}
			},
		},
	}
}
