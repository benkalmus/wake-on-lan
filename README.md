# WOL http service

Web app which runs on one of your servers, e.g Raspberry PI consuming very little power.
When HTTP route is hit /wake/<device name or MAC> it will broadcast WakeOnLan magic packet to that device.
Ensure WOL is configured in UEFI/BIOS, in your OS (Ubuntu for me).

## Build the app

```sh
go build -o wol-http .
```

## To Run (HTTP service)

copy wol-http.service file to

```sh
# create systemd service
sudo cp wol-http.service /etc/systemd/system/wol-http.service
# copy the http server we just built:
sudo cp wol-http /usr/local/bin/wol-http
# create config dir and copy it
sudo mkdir /etc/wake-on-lan/
sudo cp config.json /etc/wake-on-lan/

sudo systemd daemon-reload
sudo systemd enable wol-http
sudo systemd start wol-http
```

ensure wol-config exists in the specified systemd service directory

## Enabling Wake on LAN

I created a script that will auto enable wake on lan on Ubuntu on startup and resume.
Copy it to your systemd directory and **change the network interface** (mine is `eno1`).

[https://help.ubuntu.com/community/WakeOnLan](https://help.ubuntu.com/community/WakeOnLan)

```sh
sudo cp wol-enable.service /etc/systemd/system/
sudo systemd enable wol-enable.service
sudo systemd start wol-enable.service
sudo systemd daemon-reload
```

### NOTE

Check that the following is set:

- `WOL_DISABLE=N` in the file `/etc/tlp.conf`
- change `eno1` to your interface (check with `ip a`)
