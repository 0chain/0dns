#!/bin/bash

set -e

root=$(pwd)

cd ../
zdns=$(pwd)

cd ./code/go/0dns.io
code=$(pwd)
cd $root

#fix docker network issue on MacOS

hostname=`ifconfig | grep "inet " | grep -Fv 127.0.0.1 | grep broadcast | awk '{print $2}'`


#clean all data
clean() {
    cd $root
    rm -rf ./data && echo "config/data are removed"
}

#install vscode debugger
install_debuggger() {
    [ -d ../code/go/0dns.io\/.vscode ] || mkdir -p ../code/go/0dns.io/.vscode
    sed "s/Hostname/$hostname/g" launch.json > ../code/go/0dns.io/.vscode/launch.json
    echo "debugbbers are installed"
}

start() {

    cd $root
    
    [ -d ./data/config ] || mkdir -p ./data/config

    cp -f ../docker.local/config/0dns.yaml ./data/config/
    find ./data/config -name "0dns.yaml" -exec sed -i '' 's/use_https: true/use_https: false/g' {} \;
    find ./data/config -name "0dns.yaml" -exec sed -i '' 's/use_path: true/use_path: false/g' {} \;


    cp -f ../../0chain/docker.local/config/b0magicBlock_4_miners_2_sharders.json ./data/config/magic_block.json

    find ./data/config -name "magic_block.json" -exec sed -i '' 's/198.18.0.71/127.0.0.1/g' {} \;
    find ./data/config -name "magic_block.json" -exec sed -i '' "s/198.18.0.72/127.0.0.1/g" {} \;
    find ./data/config -name "magic_block.json" -exec sed -i '' "s/198.18.0.73/127.0.0.1/g" {} \;
    find ./data/config -name "magic_block.json" -exec sed -i '' "s/198.18.0.74/127.0.0.1/g" {} \;
    find ./data/config -name "magic_block.json" -exec sed -i '' "s/198.18.0.81/127.0.0.1/g" {} \;
    find ./data/config -name "magic_block.json" -exec sed -i '' "s/198.18.0.82/127.0.0.1/g" {} \;
    

    cd $code

    cd ./zdnscore/zdns

    # Build bls with CGO_LDFLAGS and CGO_CPPFLAGS to fix `ld: library not found for -lcrypto`
    export CGO_LDFLAGS="-L/usr/local/opt/openssl@1.1/lib"
    export CGO_CPPFLAGS="-I/usr/local/opt/openssl@1.1/include"

    GIT_COMMIT=$GIT_COMMIT
    go build -o $root/data/zdns -v -tags "bn256 development" -ldflags "-X 0chain.net/core/build.BuildTag=$GIT_COMMIT"

    cd $root/data/
    ./zdns --deployment_mode 0 --magic_block $zdns/dev.local/data/config/magic_block.json --config_dir $zdns/dev.local/data/config --log_dir $zdns/dev.local/data/log
}


echo "
**********************************************
  Welcome to 0dns development CLI 
**********************************************

"

echo "Hostname: $hostname"


echo " "
echo "Please select what are you working on: "

select i in "start" "clean" "install debugers on .vscode/launch.json"; do
    case $i in
        "start"                                     ) start                     ;;
        "clean"                                     ) clean                     ;;
        "install debugers on .vscode/launch.json"   ) install_debuggger         ;;
    esac
done

