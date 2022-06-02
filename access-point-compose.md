# Setup

I've been configuring an IOT network, which has a Raspberry Pi as a server, and several iot
nodes. They all communicate over wifi. During development, they all need internet access.

However, I change location quite frequently so reconfiguring them all with static IPs and new
wifi credentials every time I change wifi location was a real pain.

These instructions will configure a Raspberry Pi Zero 2 W as a wifi access point so all the
iot nodes can connect to it to join the local network. That way they can all have a predictable
internal IP address range, and a stable set of wifi credentials.

Not only is the Raspberry Pi server an access point, it also connects to an external wifi network
using standard DHCP, and routes internet traffic between the internal iot network and the external
public network.

Additionally, these instructions are assuming you don't have a keyboard and screen attached to
the Pi. It's all configured in a headless SSH session. We will approximately follow
[these instructions](https://imti.co/iot-wifi/) to install docker and pull image (but we skip
the `disable wpa_supplicant` step until the end, because it will disable wifi and lock you out
of the Pi):

## Things you might need on your Mac (optional)
* Jq tool for displaying json on command line: `brew install jq`
* Nmap network scanner `brew install nmap`
* Make a shortcut to the airport tool so it's easier to use: `sudo ln -s /System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport /usr/local/bin/airport`

## Instructions
1) We'll be using the [Pi Imager](https://www.raspberrypi.com/software/) tool to write the Pi OS to an SD card.
2) Remember to set hostname, ssh access and locale settings in the `Advanced options` panel before you flash.
3) Flash Raspberry Pi OS Lite 32-bit version. I've tested with the 2022-04-04 OS version.
4) *Mac* Try to ssh in with the hostname configured in the Imager settings: `ssh pi@XXX.local`
5) *Mac* If that fails, scan local network for the Pi (here XXX is the local range that your Mac is on): `sudo nmap -sn 192.168.XXX.0/24 | grep -B 2 Raspberry`, then ssh to Pi on local address: `ssh pi@192.168.XXX.XXX`
6) *Pi* `sudo apt update`
7) *Pi* `sudo apt upgrade -y`
8) *Pi* Install Docker: `curl -sSL https://get.docker.com | sh`
9) *Pi* Update permissions: `sudo usermod -aG docker pi`
10) Install `docker-compose`:
11) *Pi* `sudo apt-get install libffi-dev libssl-dev`
12) *Pi* `sudo apt install python3-dev`
13) *Pi* `sudo apt-get install -y python3 python3-pip`
14) *Pi* `sudo pip3 install docker-compose`
15) *Pi* `sudo systemctl enable docker`
16) *Pi* Reboot: `sudo reboot`
17) SSH will disconnect, so reconnect once it's rebooted.
18) *Pi* Create `docker-compose.yml`:
```
version: '3'
services:
  wifi:
    container_name: wifi
    image: "cjimti/iotwifi"
    volumes:
      - /home/pi/config/wifi:/cfg
      - /etc/localtime:/etc/localtime:ro
    restart: unless-stopped
    privileged: true
    network_mode: host
```
19) *Pi* `mkdir config && mkdir config/wifi`
20) *Pi* Create `config/wifi/wificfg.json` config file and edit credentials:
```
{
    "dnsmasq_cfg": {
      "address": "/#/192.168.27.1",
      "dhcp_range": "192.168.27.100,192.168.27.150,1h",
      "vendor_class": "set:device,IoT"
    },
    "host_apd_cfg": {
       "ip": "192.168.27.1",
       "ssid": "INTERNAL_SSID_IN_HERE",
       "wpa_passphrase":"INTERNAL_PASS_IN_HERE",
       "channel": "6"
    },
      "wpa_supplicant_cfg": {
        "cfg_file": "/etc/wpa_supplicant/wpa_supplicant.conf"
    }
}
```
12) *Pi* Start docker: `docker-compose up -d`
13) Pi will freeze and stop pinging for a few seconds here, but give it time and the SSH connection will start responding again.
14) Enable forwarding:
15) *Pi* `sudo sysctl net.ipv4.ip_forward=1`
16) *Pi* `sudo iptables -t nat -A POSTROUTING -o wlan0 -j MASQUERADE`
17) *Pi* `sudo iptables -A FORWARD -i wlan0 -o uap0 -m state --state RELATED,ESTABLISHED -j ACCEPT`
18) *Pi* `sudo iptables -A FORWARD -i uap0 -o wlan0 -j ACCEPT`
19) *Pi* Enable and configure `iptables` persistence: `sudo apt install iptables-persistent`
20) Now we can return to the `disable wpa_supplicant` step:
21) *Pi* `sudo systemctl mask wpa_supplicant.service`
22) *Pi* `sudo mv /sbin/wpa_supplicant /sbin/no_wpa_supplicant`
23) *Pi* `sudo pkill wpa_supplicant`
24) SSH to Pi will freeze because the Pi's wifi will be disconnected.
25) Reboot Pi
26) *Mac* Check you have the Pi's internal wifi network: `airport scan | grep SSID_IN_HERE`
27) *Mac* Connect to the internal wifi newtork (you won't have internet yet).
28) Once you're connected to this network you'll be able to ssh to the Pi on `192.168.27.1`, but we shouldn't need to do that right now.
29) Now you can use curl to configure the Pi's external wifi:
30) *Mac* Get Pi to return status: `curl -w "\n" http://192.168.27.1:8080/status | jq`
31) *Mac* Get Pi to scan wifi networks: `curl -w "\n" http://192.168.27.1:8080/scan | jq`
32) *Mac* Get Pi to connect to a network: `curl -w "\n" -d '{"ssid":"SSID_IN_HERE", "psk":"PASSWORD_IN_HERE"}' -H "Content-Type: application/json" -X POST http://192.168.27.1:8080/connect | jq`
33) Internet should now be working on Mac
