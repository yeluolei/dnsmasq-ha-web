# dnsmasq-ha-web
A high availability dnsmasq deployment with web UI for edit hosts file.

This is a common problem if you want to set up dns server all by yourselves.

This project contains two parts:
 1. web based edit tool for hosts file, it will save the hosts record in sqlite 3 data base and generate a hosts file at /etc/dnsmasq-ha-web/hosts
 2. a configuration based on lsyncd and rsync for remote hosts file synchronization, which make sure the hosts file for dnsmasq servers is on sync

You can either install the dnsmasq-ha-web server with a standalone instance or
install it with one of your dnsmasq instance. The script make sure there is no
conflict.

## Instalation
### dnsmasq server
1. Config no password sudo permission to remote server, because we need to restart dnsmasq remotely

       echo "ubuntu ALL=(ALL) NOPASSWD:ALL" | sudo EDITOR='tee -a' visudo

2. Install required packages

       sudo apt install dnsmasq rsync

3. Copy the dnsmasq config in `scripts/etc/dnsmasq.d/ha-web.conf` to `/etc/dnsmasq.d/ha-web.conf`
4. Copy the scripts from `scripts/local/bin/dnsmasq-rsync.sh` to `/local/bin/dnsmasq-rsync.sh`
5. Change the permission of the script file

       sudo chmod +x /local/bin/dnsmasq-rsync.sh

6. Create the dir to add hosts file


       sudo mkdir -p /etc/dnsmasq
       sudo touch /etc/dnsmasq/hosts


### dnsmasq-ha-web server
1. Install lsyncd and rsync

       sudo apt install -y lsyncd rsync

2. Generate ssh key access and upload to dns servers

       sudo mkdir -p /etc/lsyncd
       sudo ssh-keygen -t rsa -f /etc/lsyncd/id_rsa
       ssh-copy-id -i /etc/lsyncd/id_rsa.pub <your user>@<your remote server>

3. Copy the config file in `scripts/etc/lsyncd/lsyncd.conf.lua` to the destination `/etc/lsyncd/lsyncd.conf.lua`
4. Restart lsyncd service

       sudo service lsyncd restart

5. Copy the binary build in this repo to the destination `bin/dnsmasq-ha-web` -> `/usr/local/bin/dnsmasq-ha-web`
6. Copy the service config file in `scripts/etc/systemd/system/dnsmasq-ha-web.service` -> `/etc/systemd/system/dnsmasq-ha-web.service`
7. Start the service and check the status

       sudo service dnsmasq-ha-web start


## Build
1. Build the web frontend

       cd frontend
       yarn install
       yarn build

2. Build the go application, you need gcc installed, since this application is based
on sqlite3 and the go-sqlite3 package need CGO enabled

       make build