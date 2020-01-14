#!/bin/bash

echo "Get Postgres"
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -sc)-pgdg main" > /etc/apt/sources.list.d/PostgreSQL.list'
sudo apt-get -y update
sudo apt-get -y install postgresql-10

echo "Testing Install"
psql -V

echo "Run testutils/test-srv.sh to start the StructQL test server"
