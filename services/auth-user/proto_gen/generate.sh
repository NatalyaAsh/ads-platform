#!/bin/bash

# Переходим в корень проекта
cd "$(dirname "$0")/.."

echo "Current directory: $(pwd)"

# Проверяем, существует ли файл
if [ ! -f "proto/auth_user/v1/auth_user.proto" ]; then
    echo "❌ Error: proto/auth_user/v1/auth_user.proto not found!"
    echo "Files in proto directory:"
    ls -la proto/auth_user/v1/
    exit 1
fi

# Создаём папку для сгенерированных файлов
mkdir -p proto/gen/go

# Генерируем код
protoc \
    --proto_path=proto \
    --go_out=proto/gen/go \
    --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen/go \
    --go-grpc_opt=paths=source_relative \
    proto/auth_user/v1/auth_user.proto

echo "✅ gRPC code generated successfully!"