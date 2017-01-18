FROM registry.clickyab.ae/clickyab/ubuntu
MAINTAINER f0ruD <fzero@rubi.gd>

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y \
    wget git libevent1-dev libconfig-dev liblua5.1-0-dev lua5.1 libjansson-dev autoconf bison build-essential libssl-dev \
    && wget https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64.deb \
    && dpkg -i dumb-init_*.deb \
    && mkdir /app \
    && rm dumb-init_*.deb \
    && rm -rf /var/lib/apt/lists/*

# Make ssh dir
RUN mkdir /root/.ssh/
RUN chmod 700 /root/.ssh/
# Copy over private key, and set permissions
ADD github_rsa /root/.ssh/id_rsa
ADD github_rsa.pub /root/.ssh/id_rsa.pub
# Create known_hosts
RUN touch /root/.ssh/known_hosts
RUN chmod 600 root/.ssh/*
# Add bitbuckets key
RUN ssh-keyscan github.com > /root/.ssh/known_hosts  && ssh git@github.com || true
RUN git clone --recursive "git@github.com:fzerorubigd/tg.git" /tg && cd /tg && ls && ./configure && make

ADD ./init-script.sh /sbin/entry-point.sh
RUN chmod a+x /sbin/entry-point.sh

RUN useradd -r -m telegramd

VOLUME /home/telegramd/.telegram-cli

ENTRYPOINT ["/sbin/entry-point.sh"]

CMD ["tg"]
