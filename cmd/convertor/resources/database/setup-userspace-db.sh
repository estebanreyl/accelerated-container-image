#!/bin/sh
# Reference Script for getting started with mysql DB for userspace convertor
sudo apt update
sudo apt install mysql-server
sudo service mysql start
sudo cat ./mysql.conf | mysql