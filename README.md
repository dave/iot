# Setup

1) Flash Raspberry Pi OS Lite 32-bit version using Pi Imager to SD card.
2) Remember to set hostname, ssh access and locale settings in advanced.
3) *Mac* Try to ssh in with the hostname configured in the Imager settings: `ssh pi@davespi.local`
4) *Mac* If that fails, use network scanner to find the DHCP assigned ip address: `sudo nmap -sn 192.168.1.0/24 | grep -B 2 Raspberry`
5) *Mac* Log in to Pi with on 192.168.X.X address: `ssh pi@192.168.X.X` 
6) We will follow [these instructions](https://imti.co/iot-wifi/) to install docker and pull image, but we skip the `disable wpa_supplicant` step until the end (because it will disable wifi):
7) *Pi* Install Docker: `curl -sSL https://get.docker.com | sh`
8) *Pi* Update permissions: `sudo usermod -aG docker pi`
9) *Pi* Reboot: `sudo reboot`
10) *Pi* Pull the Docker image: `docker pull cjimti/iotwifi`
11) *Pi* Create `wificfg.json`:

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

12) *Pi* Start Docker container: `docker run --restart=unless-stopped -d --privileged --net host -v $(pwd)/wificfg.json:/cfg/wificfg.json cjimti/iotwifi`
13) Pi will freeze and stop pinging for a few seconds here, but give it time and the SSH connection will start responding again.
14) Enable forwarding:
15) *Pi* `sudo sysctl net.ipv4.ip_forward=1`
16) *Pi* `sudo iptables -t nat -A POSTROUTING -o wlan0 -j MASQUERADE`
17) *Pi* `sudo iptables -A FORWARD -i wlan0 -o uap0 -m state --state RELATED,ESTABLISHED -j ACCEPT`
18) *Pi* `sudo iptables -A FORWARD -i uap0 -o wlan0 -j ACCEPT`
19) *Pi* For `iptables` changes to persist across reboots, we need to install and configure this: `sudo apt install iptables-persistent`
20) Now we can return to the `disable wpa_supplicant` step:
21) *Pi* `sudo systemctl mask wpa_supplicant.service`
22) *Pi* `sudo mv /sbin/wpa_supplicant /sbin/no_wpa_supplicant`
23) *Pi* `sudo pkill wpa_supplicant`
24) SSH to Pi will freeze because the Pi's wifi will be disconnected.
25) Reboot Pi
26) *Mac* Make a shortcut to the airport tool so it's easier to use: `sudo ln -s /System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport /usr/local/bin/airport`
27) *Mac* Check you have the wifi network: `airport scan | grep SSID_IN_HERE`
28) *Mac* Connect to the internal wifi newtork (you won't have internet).
29) Now you can use the internal address of the Pi to configure it's external wifi:
30) *Mac* Install `jq` tool for displaying json on command line if you don't have it: `brew install jq`
31) *Mac* Get Pi to return status: `curl -w "\n" http://192.168.27.1:8080/status | jq`
32) *Mac* Get Pi to scan wifi networks: `curl -w "\n" http://192.168.27.1:8080/scan | jq`
33) *Mac* Get Pi to connect to a network: `curl -w "\n" -d '{"ssid":"SSID", "psk":"PASS"}' -H "Content-Type: application/json" -X POST http://192.168.27.1:8080/connect | jq`
34) Internet should now be working on Mac
