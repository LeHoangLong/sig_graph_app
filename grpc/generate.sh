protoc --go_out=../backend/internal/grpc --go_opt=paths=source_relative \
    --go-grpc_out=../backend/internal/grpc --go-grpc_opt=paths=source_relative $1

protoc --dart_out=grpc:../frontend/lib/grpc  $1
