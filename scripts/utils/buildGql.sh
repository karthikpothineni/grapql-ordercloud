#!/bin/bash
# source ./scripts/utils/func.sh
build_gql(){
  go get github.com/99designs/gqlgen/internal/imports@v0.17.20
  go get github.com/99designs/gqlgen/codegen/config@v0.17.20
  go get github.com/99designs/gqlgen@v0.17.20
  go run github.com/99designs/gqlgen generate
}
build_gql

