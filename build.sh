#!/bin/bash

# build
cd src
go build -o ../build/img-sort
cd ..

# copy bin
mkdir -p build/bin
cp -r src/bin/* build/bin
