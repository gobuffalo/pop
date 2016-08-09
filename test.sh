#!/bin/bash
set -e
clear

go build -o tsoda ./soda

function test {
  echo "Testing $1"
  export SODA_DIALECT=$1
  echo ./tsoda -v
  echo $SODA_DIALECT
  ! ./tsoda drop -e $SODA_DIALECT -c ./database.yml
  ! ./tsoda create -e $SODA_DIALECT -c ./database.yml
  ./tsoda migrate -e $SODA_DIALECT -c ./database.yml
  go test ./...
}

test "postgres"
test "sqlite"
# test "mysql"

rm tsoda
