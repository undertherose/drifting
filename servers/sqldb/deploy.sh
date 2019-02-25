#!/usr/bin/env bash
./build.sh

export MYSQL_ROOT_PASSWORD=sqldbpassword
export MYSQL_DATABASE=auth

ssh -i ~/.ssh/finalPrivKey.pem ec2-user@18.222.243.235 << EOF

    docker network rm queueNetwork
    docker network create queueNetwork
    
    docker rm -f finalsqldb

    docker run -d --name finalsqldb \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e MYSQL_DATABASE=$MYSQL_DATABASE \
    --network queueNetwork \
    koolkids441/finalsqldb

EOF