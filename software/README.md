# Software
The server software is a program that will run within the network and handle file saving, encryption, internet syncing, and deletion. All of these will run independent from each other on timers.

## Capture
Capture will download the most recent server from all cameras on an interval. Each camera will have a folder for it's images, and the files will be the timestamp when the capture sequence was activated.

# Encrypting
On an interval, all current images will be saved into some encrypted format. All plaintext images will be cleared from the captures.

# Internet syncing
On an interval, all saved encrypted copies will be synced to other registered servers.

# Deletion
On an interval, all encrypted chunks older than a configured date are deleted.

## Configuration
config.json is used to define all configuration.

## Container
A docker-compose.yml file will be provided to allow easy deployment of the software.        