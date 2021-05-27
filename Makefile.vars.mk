arch ?= amd64

ifeq ("$(shell uname -s)", "Darwin")
os ?= darwin
else
os ?= linux
endif

curl_args ?= --location --fail --silent --show-error

test_dir ?= $(shell pwd)/.test

backup_dir ?= $(test_dir)/backup
restore_dir ?= $(test_dir)/restore

stats_url ?= http://localhost:8091

restic_version ?= 0.12.0
restic_path ?= $(test_dir)/restic
restic_pid ?= $(restic_path).pid
restic_url ?= https://github.com/restic/restic/releases/download/v$(restic_version)/restic_$(restic_version)_$(os)_$(arch).bz2
restic_password ?= repopw

minio_port ?= 9000
minio_host ?= localhost
minio_address = $(minio_host):$(minio_port)
minio_path ?= $(test_dir)/minio
minio_data ?= $(minio_path).d/data
minio_config ?= $(minio_path).d/config
minio_root_user ?= accesskey
minio_root_password ?= secretkey
minio_pid ?= $(minio_path).pid
minio_url ?= https://dl.min.io/server/minio/release/$(os)-$(arch)/minio

docker_tag ?= vshn/wrestic:snapshot
