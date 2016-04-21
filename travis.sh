export GOVENDOREXPERIMENT="1"

go get -u -d -t -v github.com/Masterminds/glide
go install github.com/Masterminds/glide
glide install
# go test $(glide novendor)
