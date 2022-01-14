gen:
	oapi-codegen -generate="types,chi-server,spec" -package=v1 api/openapi/v1.yaml  > pkg/api/openapi/v1.gen.go