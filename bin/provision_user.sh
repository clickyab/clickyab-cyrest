#!/bin/bash -x
set -euo pipefail

echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/cyrest/bin" >> /home/develop/.zshrc

cd /home/develop/cyrest

make conditional-restore
make codegen
make all
make migup
