#!/bin/bash -x
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/go/bin" >> /home/develop/.zshrc

echo "Waiting for services to start..."
sleep 5

make -f /home/develop/helium/Makefile postgres-setup POSTGRES_USER=postgres
echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.5/main/pg_hba.conf
echo "listen_addresses='*'" >> /etc/postgresql/9.5/main/postgresql.conf
/etc/init.d/postgresql restart

#cp /home/develop/helium/configs/helium.conf /etc/nginx/sites-available/default
#service nginx restart

sudo -u develop /home/develop/helium/bin/provision_user.sh