
function logInfo() {
    echo "INFO:" $@ 1>&2
}

function main(){
    version=$(go run cmd/ec-tools.go -v |awk '{print $3}')
    logInfo "set main.Version to ${version}"  1>&2
    go build  -ldflags "-X main.Version=${version}" cmd/ec-tools.go
}
main
