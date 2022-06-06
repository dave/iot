
# sudo make init
init:
	@echo 'Section "Extensions"\n    Option      "DPMS" "Disable"\nEndSection\n\nSection "ServerFlags"\n    Option "StandbyTime" "0"\n    Option "SuspendTime" "0"\n    Option "OffTime" "0"\n    Option "BlankTime" "0"\nEndSection\n' > /usr/share/X11/xorg.conf.d/disable-screensaver.conf


build:
	go build github.com/dave/iot/ui

run:
	sudo xinit "./ui"

screensaver:
	xset -dpms -display :0.0
	xset s off -display :0.0
	xset s noblank -display :0.0

query:
	xset q -display :0.0

#-geometry=1920x1080+0+0

all: build run
