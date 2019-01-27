# Phingester

Ingest your photos from attached devices automatically onto a disk.

## Features

* Scans a directory periodically for mounted devices that look like camera memory cards
* Copies files to target directory while trying to avoid collisions
* Optionally changes owner
* Minimal memory footprint (less than 5 MB)

## Requirements

* `rsync`
* `Go` (only for compiling, not necessary on target)

Build using `go build`. If you want to cross compile, use `GOOS=freebsd GOARCH=amd64 go build` (substitute freebsd for target operating system and amd64 for target architecture).

## Configuration

All configuration can be performed using environmental variables:

	PHINGESTER_SCANPATH= "/media"                 // the path that will be searched for camera memory cards
	PHINGESTER_DEST=     "$HOME/phingester_media" // target path where files will be copied to
	PHINGESTER_OWNER=    ""                       // optional, if set, copied files will belong to this user

## Platform Specific Remarks

### FreeBSD Setup

You'll need autofs for automatic mounting of USB devices. Ensure that you have `/media -media -nosuid` in `/etc/auto_master`. Ensure that autofs is enabled in `/etc/rc.conf`: `autofs_enable="YES"`.

Start autofs immediately by running:

	/etc/rc.d/automount start
	/etc/rc.d/automountd start
	/etc/rc.d/autounmountd start

For more details about autofs, see https://forums.freebsd.org/threads/freebsd-from-the-trenches-using-autofs-5-to-mount-removable-media.50831/

Note that media that is formatted using exFAT will need fuse-exfat. Install `sysutil/fuse-exfat`, enable the kernel module in `/boot/loader.conf`: `fusefs_load="YES"`. You enable it immediately using `kldload fuse.ko`.
