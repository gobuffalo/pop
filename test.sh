#!/bin/bash

########################################################
# test.sh is a wrapper to execute integration tests for
# pop.
########################################################

set -e
clear

VERBOSE=""
DEBUG='NO'

for i in "$@"
do
case $i in
    -v)
    VERBOSE="-v"
    shift
    ;;
    -d)
    DEBUG='YES'
    shift
    ;;
    *)
      # unknown option
    ;;
esac
done

function cleanup {
  echo "Cleanup resources..."
  docker-compose down
  rm tsoda
  find ./sql_scripts/sqlite -name *.sqlite* -delete
}
# defer cleanup, so it will be executed even after premature exit
trap cleanup EXIT

docker-compose up -d
sleep 4 # Ensure mysql is online

go build -v -tags sqlite -o tsoda ./soda

export GO111MODULE=on

function test {
  echo "!!! Testing $1"
  export SODA_DIALECT=$1
  echo ./tsoda -v
  echo "Setup..."
  ./tsoda drop -e $SODA_DIALECT -c ./database.yml -p ./testdata/migrations
  ./tsoda create -e $SODA_DIALECT -c ./database.yml -p ./testdata/migrations
  ./tsoda migrate -e $SODA_DIALECT -c ./database.yml -p ./testdata/migrations
  echo "Test..."
  go test -race -tags sqlite $VERBOSE ./... -count=1
}

function debug_test {
    echo "!!! Debug Testing $1"
    export SODA_DIALECT=$1
    echo ./tsoda -v
    echo "Setup..."
    ./tsoda drop -e $SODA_DIALECT -c ./database.yml -p ./testdata/migrations
    ./tsoda create -e $SODA_DIALECT -c ./database.yml -p ./testdata/migrations
    ./tsoda migrate -e $SODA_DIALECT -c ./database.yml -p ./testdata/migrations
    echo "Test and debug..."
    dlv test github.com/gobuffalo/pop
}

dialects=("postgres" "cockroach" "mysql" "sqlite")

for dialect in "${dialects[@]}" ; do
  if [ $DEBUG = 'NO' ]; then
  test ${dialect}
  else
  debug_test ${dialect}
  fi
done
