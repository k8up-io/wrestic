#Test cases:

## Backup (PVC and stdin)
* backup files with correct metadata (namespace/PVC name)
* webhook on completion (success and failure)
  * namespace
  * mounted pvcs
  * backup stats (new/changed/unchanged/errors)
  * list of all snapshots (to be changed)
* prometheus push (success and failure)
* don't fail on simple file errors
* fail on repository issues
* fail on source issues (PVC not available or pod not available)
* add file-extension to stdin backup
## Restore
* restore files to s3 or folder correctly
* webhook on completion
  * restore s3 url
  * snapshot id
  * restored files
* no prometheus
* fail on single file if S3 restore due to tar corruption (obsolete with AMZE-1290)
* fail on repository issues (source and target)
## Archive
* trigger s3 restore for the newest snapshot for each PVC/stdin
* no webhook/prometheus
## Check
* fail pod on repo error
* no webhook/prometheus (yet)
## InitRepo
* create new repo if bucket doesn't exist
* timeout if repo not available
## Prune
* remove all locks before running
* clean up snapshots according to retention policy
* webhook
  * complete list of all snapshots
* no prometheus
## List Snapshots
* get list of all snapshots in the repository
* no snapshot/webhook
## Unlock
* unlock all stale locks on the repository
