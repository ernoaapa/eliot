#!/bin/sh

# To test container attach, catch the SIGINT signal and 
# echo something before closing so it's clear that it's working
sigintHandler() {
    echo "SIGINT received! I will stop the process now..."
    exit
}

trap sigintHandler INT

while [ 1 ]; do
    echo "Hello world!"
    sleep 1
done