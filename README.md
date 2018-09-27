# Wrestic

Wrapper for restic.

## Configuration
All the configuration needed can be done via environment variables:

* `HOSTNAME` overwrite the hostname used for the backup
* `KEEP_LAST` amount of last backups that should be kept
* `KEEP_HOURLY` amount of hourly backups that should be kept
* `KEEP_DAILY` amount of daily backups that should be kept
* `KEEP_WEEKLY` amount of weekly backups that should be kept
* `KEEP_MONTHLY` amount of monthly backups that should be kept
* `KEEP_YEARLY` amount of yearly backups that should be kept
* `KEEP_TAG` what tags should be kept **Not yet implemented**
* `PROM_URL` Prometheus push gateway url
* `BACKUP_DIR` directory that should get backed up, default: `/data`
* `STATS_URL` url where additional afte backup stats get pushed to
* `BACKUP_LIST_TIMEOUT` set the timeout for listing snapshots, default 30 secs

## Execution
First build the container:

```
cd cmd/wrestic
docker build -t wrestic/wrestic .
```

Then run the container and mount the folders you'd like to be backed up to `/data`:
```
docker run -e "HOSTNAME=test" -e "PROM_URL=http://192.168.1.43:9091" -v /path/to/back:/data/ wrestic/wrestic
```

Run a check of the repository:
```
docker run -e "HOSTNAME=test" -e "PROM_URL=http://192.168.1.43:9091" -v /path/to/back:/data/ wrestic/wrestic -check
```

Run a restore to disk:
```
docker run -e "HOSTNAME=test" -v /path/for/restore:/restore wrestic -restore -restoreType folder
```

Run a restore to disk with a filter:
```
docker run -e "HOSTNAME=test" -v /patch/for/restore:/restore wrestic -restore -restoreType folder -restoreFilter /var/mysql
```

Run a restore to S3:
```
docker run -e "HOSTNAME=test" -e "RESTORE_S3ENDPOINT=http://localhost:9000/bucketName" -e "RESTORE_ACCESSKEYID=1324" -e "RESTORE_SECRETACCESSKEY=secret" wrestic -restore -restoreType s3
```

The container will exit after the job is done. If a valid `PROM_URL` is provided it will push metrics there.
