#!/bin/bash -x
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/go/bin" >> /home/develop/.zshrc

echo "Waiting for services to start..."
sleep 5

make -f /home/develop/cyrest/Makefile mysql-setup POSTGRES_USER=postgres

sudo -u develop /home/develop/cyrest/bin/provision_user.sh