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
sudo cp  wol-http.service /etc/systemd/system/wol-http.service

sudo systemd daemon-reload
sudo systemd enable wol-http
sudo systemd start wol-http
```

ensure wol-config exists in the specified systemd service directory

## Enabling Wake on LAN

https://help.ubuntu.com/community/WakeOnLan

```sh
sudo cp wol-enable.service /etc/systemd/system/
sudo systemd enable wol-enable.service
sudo systemd start wol-enable.service
sudo systemd daemon-reload
```

### NOTE

Check that the following is set:

- `WOL_DISABLE=N` in the file `/etc/tlp.conf`
