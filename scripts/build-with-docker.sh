realPath=$(realpath scripts/build-with-docker.sh)
scriptPath=$(dirname ${realPath})
projectPath=$(dirname ${scriptPath})

function getContainerCmd(){
    local containerCmd=""
    for cmd in podman docker
    do
        which ${cmd} > /dev/null 2>&1
        if [[ $? -eq 0 ]]; then
            containerCmd="${cmd}"
            break
        fi
    done
    if [[ -z $containerCmd ]]; then
        echo "ERROR" "仅支持 podman 或 docker"
        exit 1
    fi
    echo $containerCmd
}

function main(){
    local containerCmd=$(getContainerCmd)

    echo "INFO" "使用命令 ${containerCmd}"
    echo "INFO" "项目地址" $projectPath

    cd ${scriptPath}
    ${containerCmd} build -v ${projectPath}:/root/ec-tools \
        --target Centos7Base \
        -t ec-tools:base \
        -f centos7.Dockerfile \
        ./
    if [[ $? -ne 0 ]]; then
        echo "ERROR" "基础镜像构建失败"
        exit 1
    fi
    ${containerCmd} build -v ${projectPath}:/root/ec-tools ./ \
        --target Centos7Builder \
        --cache-from ec-tools:base \
        -t ec-tools:builder \
        -f centos7.Dockerfile
    if [[ $? -eq 0 ]]; then
        ${containerCmd} rmi ec-tools:builder
    fi
}

main $*
