# -*- mode: ruby -*-
# vi: set ft=ruby :

UBUNTU = "ubuntu/trusty64"
RHEL = "iamseth/rhel-7.3"


Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.

  config.vm.provider "virtualbox" do |v|
    v.memory = 4096
    v.cpus = 2
  end

  config.vm.define "machine-1" do |machine|
    # machine.ssh.username = 'root'
    # machine.ssh.password = 'vagrant'
    # machine.ssh.insert_key = 'true'
    machine.vm.box = UBUNTU
    # machine.vm.box_version = "1.0.0"
  end

end
