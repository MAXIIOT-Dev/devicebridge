.PHONY: build generate clean swagger
GO_EXTRA_BUILD_ARGS=-a -installsuffix cgo
build: swagger generate
	# todo 
	@echo "starting vbasebridge complie"
	@mkdir -p build 
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${GO_EXTRA_BUILD_ARGS} -o build/vbasebridge cmd/vbasebridge/main.go
	@echo "complete vbasebride compile"

generate:
	@echo "generate migrate"
	@go generate ./storage/db.go

clean:
	@rm -f ./storage/migrations_gen.go
	@rm -rf ./docs

swagger:
	@echo "generate swagger api docs"
	@swag init --generalInfo routers/routers.go 