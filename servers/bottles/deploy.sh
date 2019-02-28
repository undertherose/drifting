#!/usr/bin/env bash
./build.sh

ssh -i ~/.ssh/MyPrivKey.pem ec2-user@18.217.182.145 << EOF
    
    docker rm -f bottles

    docker pull wecancodeit/bottles

    docker run -d \
    --name bottles \
    --network driftingNetwork \
    -e PORT=80 \
    wecancodeit/bottles
    
EOF