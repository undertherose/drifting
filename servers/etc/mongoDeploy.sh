#!/usr/bin/env bash

ssh -i ~/.ssh/finalPrivKey.pem ec2-user@18.222.243.235 << EOF

    docker rm -f mongoDB

    docker run -d --name mongoDB \
    --network driftingNetwork \
    mongo

    exit
EOF
