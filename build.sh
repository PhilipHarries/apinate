#!/usr/bin/env bash

if [[ $# -ne 1 ]];then
    echo "Please specify apinate version"
    exit 1
fi

if [[ "x`which fpm`x" == "xx" ]];then
    echo "Please install fpm."
    exit 2
fi

gom install

gom build -o $GOPATH/bin/apinate

build_dir=/tmp/fpm/apinate

mkdir -p ${build_dir}/etc/apinate
mkdir -p ${build_dir}/usr/share/apinate/templates
mkdir -p ${build_dir}/usr/bin

cp ./templates/plain.tmpl ${build_dir}/usr/share/apinate/templates
cp ${GOPATH}/bin/apinate ${build_dir}/usr/bin

fpm -s dir -t rpm -n apinate -v $1 -C ${build_dir}
fpm -s dir -t deb -n apinate -v $1 -C ${build_dir}
