# AGNC (Aviso triggered GRIB to NetCDF convertor)

Meteorological data is best served fresh. At MET Norway, we have always strived to make ECMWF prognoses available as soon as possible for downstream processing and utilization. In 2016, we have developed a system broadcasting an event every time we receive a model time step from the dissemination system, triggering further actions. Dissemination goes through SFTP and we use inotify to detect a change. With Aviso, this task has become easier. Today, Aviso provides these events, and ECMWF disseminates the data via S3.

We use the Aviso client in a proof-of-concept application, which converts GRIB files to NetCDF and serves them using an OPeNDAP server (Thredds). Systemd is the orchestrator used to run Aviso and schedule conversion jobs using Faktory (https://github.com/contribsys/faktory/). The notifications include a full s3 path, which is reachable from wherever the conversion system is running.

We make use of templating in the configuration file to structure commands for the processing system. This feature makes it possible to parametrise post-processing programs and doesn’t pose any requirements on the languages or libraries used. Any callable binary will work here.

## Systemd files

`%h` is used for ~ or $HOME in systemd service files. The whole setup runs in a userspace systemd, i.e., `systemd --user`

### `~/.config/systemd/user/aviso.service`

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

### `~/.config/systemd/user/faktory-worker.service`
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

### `~/aviso/config.yaml`
```
listeners:
  - event: dissemination
    request:
      destination: <ECPDS destination>
    triggers:
      - type: command
        working_dir: $HOME
        command: ./bin/client pt --ncDir /netcdf/${request.target} --gribDir /dev/shm/ --configDir fimex-config/${request.target}/etc/ ${location}
        environment:
          FAKTORY_URL: tcp://:<faktory password>@localhost:7419
```

To leverage the flexibility of S3, we are planning to add a conversion to [Zarr](https://github.com/zarr-developers) and distribute the data further by pushing it back to the object storage.
