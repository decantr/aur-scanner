# aur-scanner

Check AUR packages against compromsied package list.

## Build
```sh
go build
```

## Usage
```sh
./aur-scanner
./aur-scanner scan
./aur-scanner --file compromised-list.md
```

## Flags
- `--url`: override the compromised-package list source
- `--file`: read the compromised-package list from a local file
- `--quiet`: print only matching package names

## Default source

By default the scanner uses `https://md.archlinux.org/s/SxbqukK6IA/download`
