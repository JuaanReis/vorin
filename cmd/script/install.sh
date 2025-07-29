#!bin/bash

git clone https://github.com/JuaanReis/vorin.git

cd vorin

go build -o vorin

./vorin -help