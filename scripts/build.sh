
function logInfo() {
    echo `date +%F-%T` "INFO" $0 $@ 1>&2
}

function main(){
    version=$(go run cmd/ec-tools.go -v |awk '{print $3}')
    if [[ -z $version ]]; then
        echo "ERROR" "version is null"
        exit 1
    fi
    logInfo "设置 main.Version=${version}"
    mkdir -p dist
    logInfo "开始构建"
    go build  -ldflags "-X main.Version=${version}" -o dist/ cmd/ec-tools.go
    if [[ $? -eq 0 ]]; then
        logInfo "构建完成"
    else
        exit 1
    fi
}
main
