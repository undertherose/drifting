#!/usr/bin/env bash
./build.sh

export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
export DB_NAME=auth
export MYSQL_ADDR=finalsqldb


ssh -i ~/.ssh/finalPrivKey.pem ec2-user@18.222.243.235 << EOF
    
    docker rm -f courses

    docker pull koolkids441/courses

    docker run -d \
    --name courses \
    --network queueNetwork \
    -e PORT=80 \
    koolkids441/courses

EOF