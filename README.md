# wol http service

## Build the app

```sh
go build -o wol-http .
```

## To Run

copy wol-http.service file to

```sh
/etc/systemd/system/wol-http.service

sudo systemd daemon-reload
sudo systemd enable wol-http
sudo systemd start wol-http
```

ensure wol-config exists in specified directory
