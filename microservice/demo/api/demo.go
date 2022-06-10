package main

import (
	commonConfig "demo/common/config"
	"demo/microservice/demo/api/internal/config"
	"demo/microservice/demo/api/internal/handler"
	"demo/microservice/demo/api/internal/svc"
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "microservice/demo/api/etc/nacos.yaml", "the config file")

func main() {
	flag.Parse()

	var c = new(config.Config)
	ctx := svc.NewServiceContext(c, commonConfig.MustLoad(*configFile, c))
	server := rest.MustNewServer(c.RestConf, rest.WithCors(), rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)
	//swagger.RegisterSwagger(server)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
