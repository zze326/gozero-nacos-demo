SWAGGER_JSON_FILE=swagger.json

api-gen:
	goctl api go -api *.api -dir . -style goZero

swagger-gen-json:
	goctl api plugin -plugin goctl-swagger="swagger -filename ${SWAGGER_JSON_FILE}" -api ci.api -dir .
