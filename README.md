# DirSync

**DirSync** is a command-line tool and Go package for synchronizing the contents of one directory with another.
It allows to copy files from a source to a target directory while optionally removing files in the target that no longer exist in the source.

---

## Features

- Recursive copy of all regular files
- Preserves permissions and modification times
- Skips files that are unchanged (based on content + timestamp)
- Optional deletion of files that are missing in the source
- Fully testable with a pluggable filesystem (`FS` interface)

---

## Usage

```bash
dirsync [--delete-missing] <source> <target>
```
```
# Basic copy
$ dirsync ./my-assets ./backup-assets

# Copy and delete anything in target not present in source
$ dirsync --delete-missing ./my-assets ./backup-assets
```
---
## Example output
```
[INFO] source path: /full/path/my-assets
[INFO] target path: /full/path/backup-assets
[INFO] Copied: /full/path/my-assets/image.png → /full/path/backup-assets/image.png
[INFO] Skipping identical file: /full/path/backup-assets/logo.svg
[INFO] Removing missing: /full/path/backup-assets/old.txt
```


