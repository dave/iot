### Setting up Raspberry Pi headless with USB connection to Mac

# From: https://gist.github.com/superdodd/06b91ba03899e47beb43078b27dc601e

1) Flash Raspberry Pi OS Lite 64-bit version (released 2022-04-04) using Pi Imager to SD card
2) Remember to set hostname, ssh access and locale settings in advanced
3) In config.txt, add this as the last line of the file: `dtoverlay=dwc2`
4) In cmdline.txt, add this as a parameter, just after the rootwait parameter: `modules-load=dwc2,g_ether`
5) Connect Pi to Mac with USB cable, making sure to connect to USB port on Pi (not PWR)
6) Boot Pi
7) Check in Mac `System Preferences > Sharing` that `Internet Sharing` and `RNDIS/Ethernet Gadget` are ticked.
8) Ping Pi: `ping raspberrypi.local`, check that it has a `192.168.X.Y` address.
9) Pull the power out while it's pinging, check that it stops (proving that this is definitely our Pi and not another on the network!)
10) Boot Pi again
11) Log in with `ssh pi@raspberrypi.local` - if you correctly configured SSH keys in step 2, no password is needed.
12) Remember `/etc/network/interfaces` stops processing after the `source` command, so you need to put a file in `/etc/network/interfaces.d`
13) Add `/etc/network/interfaces.d/new_network`:

```
allow-hotplug usb0
iface usb0 inet static
  address 192.168.2.2
  netmask 255.255.255.0
  gateway 192.168.2.1
```

13) Edit `/etc/resolvconf.conf`, uncomment / change: `nameservers=192.168.2.1`
14) On Mac, configure the `RNDIS/Ethernet Gadget` interface with the following parameters: 

* Configure IPV4: `Manually`
* IP Address: `192.168.2.1`
* Subnet Mask: `255.255.255.0`
* Router: (none)
* Advanced -> DNS, `8.8.8.8` and `1.1.1.1`

15) Reboot the Pi.

### NOTE: for some reason at this point, the Pi has internet connectivity, but pings will fail.

Now transpose into: https://imti.co/iot-wifi/#legacy-instructions-the-manual-way

1) 