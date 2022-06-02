# Install

### Getting `ui` (a Gio app) running on a Raspberry Pi Zero 2 (with Lite OS version) with a BOOX Mira e-ink display.

Comment this line out in config.txt to get BOOX Mira display working:
```
dtoverlay=vc4-kms-v3d
```

Install Go (update version to latest):
```
wget https://dl.google.com/go/go1.18.3.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.18.3.linux-armv6l.tar.gz
```

Add `go` command to the path:
```
pico ~/.profile
```

... add:
```
PATH=$PATH:/usr/local/go/bin
GOPATH=$HOME/src
```

Required for Gio:
```
sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev
```

Required for starting Gio app without desktop environment:
```
sudo apt install xinit
```

To start app:
```
go build github.com/dave/iot/ui && sudo xinit ./ui $* -- :1
```


Notes (not needed): 

Maybe: `https://snapcraft.io/install/wayland/raspbian`
Maybe: `https://www.paulligocki.com/add-gui-to-raspberry-pi-os-lite/`
Maybe: `https://forums.raspberrypi.com/viewtopic.php?t=265090`
THIS: `https://linuxconfig.org/how-to-run-x-applications-without-a-desktop-or-a-wm`

