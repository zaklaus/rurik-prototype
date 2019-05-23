#!/bin/sh

rm -rf package

# Ship the latest Windows build
echo Shipping for Windows...
make
./zip.sh game.exe
mv package/build.zip package/darkorbia-win64-x86_64.zip
# butler push package/darkorbia-win64-x86_64.zip zaklaus/darkorbia:darkorbia-windows

echo Done!