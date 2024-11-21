module github.com/NikolosHGW/auth

go 1.23.0

require (
	github.com/NikolosHGW/platform-common v0.0.0-20241107155759-c32fb5ee4060
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v1.1.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

require github.com/jmoiron/sqlx v1.4.0 // indirect

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
)
