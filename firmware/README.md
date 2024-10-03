# Firmware
A simple stack for esp32-cam to supply an authenticated http endpoint for image output

# Camera Software

wifi-http-cam-server is an ESP-IDF program for the hardware that will expose an http endpoint `/capture` to download a jpeg image from the camera.

# How to build and deploy
TODO: Update this

# Recording Software

TBD details. This will be a Go based program that periodically pulls captures off the cameras and saves them to disk.