# ⚠️ This repository is no longer in development

wrestic is now part of K8up itself.
The codebase has been migrated to https://github.com/k8up-io/k8up, please open issues there.

# Wrestic

---

**PLEASE NOTE: THIS REPOSITORY HAS BEEN ARCHIVED.**

wrestic has been merged into [K8up](https://github.com/k8up-io/k8up) 2.0 and is available as `k8up restic` command.

---

Wrapper for restic.
This is the backup runner of k8up.

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
* `STATS_URL` url where additional after backup stats get pushed to
* `BACKUP_LIST_TIMEOUT` set the timeout for listing snapshots, default 30 secs
* `RESTORE_S3ENDPOINT` s3 endpoint where the tar.gz with all files should be uploaded, example `http://localhost:9000/bucketName`
* `RESTORE_ACCESSKEYID` s3 accesKeyID for the restore s3 endpoint
* `RESTORE_SECRETACCESSKEY` s3 secretAccessKey for the restore s3 endpoint
* `BACKUPCOMMAND_ANNOTATION` name of the backup command annotation, default: `k8up.syn.tools/backupcommand`
* `FILEEXTENSION_ANNOTATION` name of the file extension annotation, default: `k8up.syn.tools/file-extension`
* `RESTIC_OPTIONS` additional options, separated with a comma (i.e. `key=value,key2=value2,…`), to pass on to `restic`, [see `--option`](https://restic.readthedocs.io/en/stable/manual_rest.html)
* `RESTIC_BINARY` defines the `restic` binary to use, default: `/usr/local/bin/restic`

Configuration for the Restic repository also has to be provided via env variables, see [the official docs](https://restic.readthedocs.io/en/stable/040_backup.html#environment-variables).

## Execution

First build the container:

```bash
cd cmd/wrestic
docker build -t wrestic/wrestic .
```

Then run the container and mount the folders you'd like to be backed up to `/data`:

```bash
docker run -e "HOSTNAME=test" -e "PROM_URL=http://192.168.1.43:9091" -v /path/to/back:/data/ wrestic/wrestic
```

Run a check of the repository:

```bash
docker run -e "HOSTNAME=test" -e "PROM_URL=http://192.168.1.43:9091" -v /path/to/back:/data/ wrestic/wrestic -check
```

Run a restore to disk:

```bash
docker run -e "HOSTNAME=test" -v /path/for/restore:/restore wrestic -restore -restoreType folder
```

Run a restore to disk with a filter:

```bash
docker run -e "HOSTNAME=test" -v /patch/for/restore:/restore wrestic -restore -restoreType folder -restoreFilter /var/mysql
```

Run a restore to S3:

```bash
docker run -e "HOSTNAME=test" -e "RESTORE_S3ENDPOINT=http://localhost:9000/bucketName" -e "RESTORE_ACCESSKEYID=1324" -e "RESTORE_SECRETACCESSKEY=secret" wrestic -restore -restoreType s3
```

The container will exit after the job is done.
If a valid `PROM_URL` is provided, it will push metrics there.

## Nonroot Image

There is a variant of this image which runs as a non-root user. Build it with the following target:

```bash
docker build --target nonroot -t wrestic/wrestic .
```

## Integration Tests

To just run the integration tests, you can execute `make integration-test`.

If you want to run the integration tests in your IDE, you need to run `make integration-test-setup` first, and you need to define the following environment variables.

```dotenv
# adjust to where you checked out the wrestic source code
RESTIC_PATH=/home/<USER>/src/vshn/wrestic/.test/restic
RESTIC_BINARY=/home/<USER>/src/vshn/wrestic/.test/restic
BACKUP_DIR=/home/<USER>/src/vshn/wrestic/.test/backup/
RESTORE_DIR=/home/<USER>/src/vshn/wrestic/.test/restore/

RESTIC_PASSWORD=repopw
RESTIC_REPOSITORY=s3:http://localhost:9000/test
RESTORE_S3ENDPOINT=http://localhost:9000/restore
AWS_SECRET_ACCESS_KEY=secretkey
AWS_ACCESS_KEY_ID=accesskey
RESTORE_ACCESSKEYID=accesskey
RESTORE_SECRETACCESSKEY=secretkey
STATS_URL=http://localhost:8091
```

The stop all background services (like Minio) and remove all cached binaries, run `make clean`.
