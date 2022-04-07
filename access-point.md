# This doesn't work


* To use Pi as a wifi access point:
* `sudo apt-get install hostapd`
* `sudo apt-get install dnsmasq`
* `sudo systemctl stop hostapd`
* `sudo systemctl stop dnsmasq`
* `sudo mv /etc/dnsmasq.conf /etc/dnsmasq.conf.orig`
* `sudo pico /etc/dnsmasq.conf` add:

```
interface=wlan0
  dhcp-range=192.168.77.100,192.168.77.200,255.255.255.0,24h
```

* `sudo pico /etc/hostapd/hostapd.conf` add:

```
interface=wlan0
# bridge=br0
hw_mode=g
channel=7
wmm_enabled=0
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
ssid=NETWORK
wpa_passphrase=PASSWORD
```

* `sudo pico /etc/default/hostapd`, find `DAEMON_CONF`, uncomment and add filename: `DAEMON_CONF="/etc/hostapd/hostapd.conf"`

* `sudo systemctl unmask hostapd`
* `sudo systemctl enable hostapd`
* `sudo systemctl start hostapd`
* `sudo systemctl start dnsmasq` 