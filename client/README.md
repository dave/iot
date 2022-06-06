To restart mosquitto after an upgrade:
brew services restart mosquitto
Or, if you don't want/need a background service you can just run:
/opt/homebrew/opt/mosquitto/sbin/mosquitto -c /opt/homebrew/etc/mosquitto/mosquitto.conf
