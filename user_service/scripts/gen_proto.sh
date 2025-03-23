#!/bin/bash

CURRENT_DIR=$1

# Ensure CURRENT_DIR is provided
if [ -z "$CURRENT_DIR" ]; then
  echo "Error: Missing directory argument."
  exit 1
fi

# Remove old generated files
rm -rf "$CURRENT_DIR/genproto"/*

# Ensure GOPATH is correctly set (fallback to default location)
GOPATH=${GOPATH:-$HOME/go}
PROTOC_GEN_GO="$GOPATH/bin/protoc-gen-go"
PROTOC_GEN_GO_GRPC="$GOPATH/bin/protoc-gen-go-grpc"

# Verify that protoc-gen-go exists
if [ ! -x "$PROTOC_GEN_GO" ] || [ ! -x "$PROTOC_GEN_GO_GRPC" ]; then
  echo "Error: protoc-gen-go or protoc-gen-go-grpc not found or not executable."
  echo "Run the following to install them:"
  echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
  exit 1
fi

# Generate protobuf files
for x in "$CURRENT_DIR/protos"/*; do
  if [ -d "$x" ]; then
    protoc \
      --plugin="protoc-gen-go=$PROTOC_GEN_GO" \
      --plugin="protoc-gen-go-grpc=$PROTOC_GEN_GO_GRPC" \
      -I="$x" -I="$CURRENT_DIR/protos" -I="/usr/local/include" \
      --go_out="$CURRENT_DIR" \
      --go-grpc_out=require_unimplemented_servers=false:"$CURRENT_DIR" \
      "$x"/*.proto
  fi
done

echo "Proto files successfully generated!"
