# Setup

`docker run -d --name homeassistant --privileged --restart=unless-stopped -e TZ=Europe/London -v /home/pi:/config --network=host ghcr.io/home-assistant/home-assistant:stable`

* Install `docker-compose`:
* `sudo apt-get install libffi-dev libssl-dev`
* `sudo apt install python3-dev`
* `sudo apt-get install -y python3 python3-pip`
* `sudo pip3 install docker-compose`
* `sudo systemctl enable docker`

* create `docker-compose.yml`:
```
version: '3'
services:
  homeassistant:
    container_name: homeassistant
    image: "ghcr.io/home-assistant/home-assistant:stable"
    volumes:
      - /home/pi/config/homeassistant:/config
      - /etc/localtime:/etc/localtime:ro
    restart: unless-stopped
    privileged: true
    network_mode: host
  esphome:
    container_name: esphome
    image: esphome/esphome
    volumes:
      - /home/piconfig/esphome:/config
      - /etc/localtime:/etc/localtime:ro
    restart: unless-stopped
    privileged: true
    network_mode: host
```

# `docker-compose up -d`