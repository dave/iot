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

12) *Pi* Start Docker container: `docker run --restart=unless-stopped --privileged --net host -v $(pwd)/wificfg.json:/cfg/wificfg.json cjimti/iotwifi`
13) *Mac* Check you have the wifi network: `/System/Library/PrivateFrameworks/Apple80211.framework/Versions/A/Resources/airport scan | grep SSID_IN_HERE`
14) Now we can return to the `disable wpa_supplicant` step:
15) *Pi* `sudo systemctl mask wpa_supplicant.service`
16) *Pi* `sudo mv /sbin/wpa_supplicant /sbin/no_wpa_supplicant`
17) *Pi* `sudo pkill wpa_supplicant`
18) The wifi connection will be disabled so you'll lose the SSH connection.
19) Reboot the Pi
20) The wifi network should come up, check as before. When it appears, connect to it (internet will not be available yet)
21) *Mac* Get Pi to return status: `curl -w "\n" http://192.168.27.1:8080/status`
22) *Mac* Get Pi to scan wifi networks: `curl -w "\n" http://192.168.27.1:8080/scan`
23) *Mac* Get Pi to connect to a network: `curl -w "\n" -d '{"ssid":"SSID", "psk":"PASS"}' -H "Content-Type: application/json" -X POST http://192.168.27.1:8080/connect`
24) Now SSH to Pi using internal IP: `ssh pi@192.168.27.1`
25) Enable forwarding:
26) *Pi* `sudo sysctl net.ipv4.ip_forward=1`
27) *Pi* `sudo iptables -t nat -A POSTROUTING -o wlan0 -j MASQUERADE`
28) *Pi* `sudo iptables -A FORWARD -i wlan0 -o uap0 -m state --state RELATED,ESTABLISHED -j ACCEPT`
29) *Pi* `sudo iptables -A FORWARD -i uap0 -o wlan0 -j ACCEPT`
30) Internet should now be working on Mac
