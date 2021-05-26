#!/bin/bash

rsync -Pav -e "ssh -p 5555" --exclude 'sync.sh'  --exclude '/config' /Users/rvm/go/src/github.com/aau-network-security/sandbox/* ubuntu@130.226.98.173:/home/ubuntu/vlad/sec03/sandbox