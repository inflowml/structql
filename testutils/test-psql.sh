sudo -u postgres psql -c "CREATE USER StructqlUser PASSWORD 'StructqlPW';"
sudo -u postgres psql -c "CREATE DATABASE testdb;"
sudo service postgresql start