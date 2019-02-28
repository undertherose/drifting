#!/usr/bin/env bash
./build.sh

export MYSQL_ROOT_PASSWORD=sqldbpassword
export MYSQL_DATABASE=auth

ssh -i ~/.ssh/MyPrivKey.pem ec2-user@18.217.182.145 << EOF

    docker network rm driftingNetwork
    docker network create driftingNetwork
    
    docker rm -f sqldb

    docker run -d --name sqldb \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e MYSQL_DATABASE=$MYSQL_DATABASE \
    --network driftingNetwork \
    wecancodeit/sqldb

EOF