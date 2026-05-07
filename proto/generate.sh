#!/bin/bash

# Переходим в корень проекта
cd "$(dirname "$0")/.."

echo "Current directory: $(pwd)"

# Создаём папку для сгенерированных файлов
mkdir -p proto/gen/go

# ========== AUTH-USER SERVICE ==========
if [ -f "proto/auth_user/v1/auth_user.proto" ]; then
    echo "📦 Generating auth-user protos..."
    protoc \
        --proto_path=proto \
        --go_out=proto/gen/go \
        --go_opt=paths=source_relative \
        --go-grpc_out=proto/gen/go \
        --go-grpc_opt=paths=source_relative \
        proto/auth_user/v1/auth_user.proto
    echo "✅ auth-user protos generated"
else
    echo "⚠️  proto/auth_user/v1/auth_user.proto not found, skipping..."
fi

# ========== AD-SEARCH SERVICE ==========
if [ -f "proto/ad_search/v1/ad_search.proto" ]; then
    echo "📦 Generating ad-search protos..."
    protoc \
        --proto_path=proto \
        --go_out=proto/gen/go \
        --go_opt=paths=source_relative \
        --go-grpc_out=proto/gen/go \
        --go-grpc_opt=paths=source_relative \
        proto/ad_search/v1/ad_search.proto
    echo "✅ ad-search protos generated"
else
    echo "⚠️  proto/ad_search/v1/ad_search.proto not found, skipping..."
fi

echo "✅ All gRPC code generated successfully!"