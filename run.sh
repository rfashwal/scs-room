#!/bin/sh

echo "#############################################"
echo "Waiting for eureka"
echo "#############################################"

while ! `nc -z discovery 8761 `; do sleep 3; done

echo "#############################################"
echo "Waiting for gateway"
echo "#############################################"

while ! `nc -z gateway 8000 `; do sleep 3; done

echo "#############################################"
echo "Ready to rumble. Starting rooms service"
echo "#############################################"
./rooms
