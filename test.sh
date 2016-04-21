#!/bin/bash

set +e

clear

echo "postgres"
SODA_DIALECT=postgres go test $(glide novendor)
echo "--------------------"
echo "mysql"
SODA_DIALECT=mysql go test $(glide novendor)
# echo "--------------------"
# echo "sqlite"
# SODA_DIALECT=sqlite go test ./...
