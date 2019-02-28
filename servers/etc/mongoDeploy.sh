#!/usr/bin/env bash

ssh -i ~/.ssh/MyPrivKey.pem ec2-user@18.217.182.145 << EOF

    docker rm -f mongoDB

    docker run -d --name mongoDB \
    --network driftingNetwork \
    mongo

    exit
EOF
