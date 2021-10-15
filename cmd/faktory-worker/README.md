# AGNC (Aviso triggered GRIB to NetCDF convertor)

## Systemd files

### `.config/systemd/user/aviso.service`

```ini
[Unit]
Description=Aviso

[Service]
WorkingDirectory=%h
ExecStart=/opt/anaconda3/bin/aviso listen %h/aviso/config.yaml
RestartSec=20
Restart=always

[Install]
WantedBy=default.target
```

### `.config/systemd/user/faktory-worker.service`
```ini
[Unit]
Description=Faktory worker

[Service]
WorkingDirectory=%h
Environment="FAKTORY_URL=tcp://:<faktory password>@localhost:7419"
Environment="MMS_API_KEY=<MMS API key>"
Environment="MMS_PRODUCTION_HUB=https://<production hub URL>:<port>"
ExecStart=%h/bin/faktory-worker
Restart=always
RestartSec=30

[Install]
WantedBy=default.target

```