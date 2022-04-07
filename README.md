# iot

### Raspberry Pi Zero 2 W setup

* Use Raspberry Pi software to burn "Lite" 64bit OS onto SD card
* Add `wpa_supplicant.conf` to root of card with contents:

```
country=GB
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
	scan_ssid=1
	ssid="XXX"
	psk="XXX"
}
```

* Remember it must be a 2.4Ghz network not 5Ghz
* Add empty file named `ssh` to the root of the card
* Boot
* Use a network scanner to find the DHCP assigned ip address (e.g. `sudo nmap -sn 192.168.1.0/24`)
* SSH into the pi: `ssh pi@192.168.X.X`, default password is `raspberry`
* Change password: `passwd`
* Update with: `sudo apt update` then `sudo apt upgrade`
* Add a second static ip by adding a file: `sudo pico /etc/dhcpcd.exit-hook` with the contents `ip address add 192.168.X.X/24 dev wlan0`
* Add second network interface to your Mac with a static ip in the same range as the static ip of the pi.
* Now you can SSH into the pi on the static ip.

IPs for my project:

```
192.168.77.10 - Laptop
192.168.77.11 - Raspberry Pi Zero (A)
192.168.77.12 - Raspberry Pi Zero (B)
...
```






