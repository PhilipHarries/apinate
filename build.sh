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

[[ -d ${build_dir} ]] && rm -rf ${build_dir}
mkdir -p ${build_dir}/etc/apinate
mkdir -p ${build_dir}/usr/share/apinate/templates
mkdir -p ${build_dir}/usr/bin
mkdir -p ${build_dir}/usr/share/man/man1

cp ./templates/plain.tmpl ${build_dir}/usr/share/apinate/templates
cp ${GOPATH}/bin/apinate ${build_dir}/usr/bin
cp ./man/apinate.1 ${build_dir}/usr/share/man/man1
gzip ${build_dir}/usr/share/man/man1/apinate.1

fpm -s dir -t rpm -n apinate -v $1 -C ${build_dir}
fpm -s dir -t deb -n apinate -v $1 -C ${build_dir}
