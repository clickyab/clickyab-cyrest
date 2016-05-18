#!/bin/bash -x
set -euo pipefail

echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/helium/bin" >> /home/develop/.zshrc
echo -e "alias psql='sudo su - postgres -c psql'"  >> /home/develop/.zshrc
echo -e "alias pgcli='sudo su - postgres -c pgcli'"  >> /home/develop/.zshrc

cd /home/develop
ln -sf /srv/http/phppgadmin

cd /home/develop/helium

make protobuf
sudo make install-protobuf

make all
make migup
make hook
