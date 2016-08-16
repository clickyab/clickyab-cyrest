# -*- mode: ruby -*-
# vi: set ft=ruby :
VAGRANTFILE_API_VERSION = "2"
# set docker as the default provider
ENV['VAGRANT_DEFAULT_PROVIDER'] = 'docker'
# disable parallellism so that the containers come up in order
ENV['VAGRANT_NO_PARALLEL'] = "1"
ENV['FORWARD_DOCKER_PORTS'] = "1"

DOCKER_HOST_NAME = "dockerhostgo"
DOCKER_HOST_VAGRANTFILE = "docker-machine/Vagrantfile"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  # define a Postgres Vagrant VM with a docker provider
    config.vm.provider "docker" do |d|
      # use a vagrant host if required (OSX & Windows)
      d.vagrant_vagrantfile = "#{DOCKER_HOST_VAGRANTFILE}"
      d.vagrant_machine = "#{DOCKER_HOST_NAME}"
      # Build the container & run it
      d.image = "docker.clickyab.com/clickyab/baseimage-go:latest"
      d.name = "clickyabconsole-go"
      d.has_ssh = true
      d.cmd = ["/bin/bash", "/home/develop/cyrest/bin/init.sh"]
      d.expose = [5432]
      d.email = "dev@clickyab.com"
      d.username = "clickyab"
      d.password = "bita123"
      d.auth_server = "docker.clickyab.com"
    end
    config.vm.synced_folder ".", "/home/develop/cyrest",
      owner: "develop",
      group: "develop",
      mount_options: ["dmode=775,fmode=664"],
      create: true

    config.ssh.username = "develop"
    config.ssh.password = "bita123"
    config.vm.network "forwarded_port", guest: 80,      host: 80      # nginx
    config.vm.network "forwarded_port", guest: 8025,    host: 8025    # mailHog
    config.vm.network "forwarded_port", guest: 5432,    host: 5432    # postgres
    config.vm.provision "shell" do |s|
        s.path   = "bin/provision.sh"
        s.args   = [%x(ip addr | grep inet | grep docker0 | awk -F" " '{print $2}'| sed -e 's/\\/.*$//')]
    end


end

