#!bin/bash

echo "[+] Installing Vorin..."

git clone https://github.com/JuaanReis/vorin.git

cd vorin

go build -o vorin

sudo mv vorin /usr/local/bin/

echo "[✓] Vorin installed!"