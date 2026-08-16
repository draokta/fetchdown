# fetchdown

Unduh release asset GitHub langsung dari terminal. Tanpa browser, tanpa API token untuk repo publik.

## Install

```sh
go install github.com/draokta/fetchdown@latest
```

## Usage

```sh
fetchdown owner/repo                # semua asset release terbaru
fetchdown owner/repo "windows"      # hanya asset yang namanya mengandung "windows"
```

Contoh:

```sh
$ fetchdown cli/cli "windows_amd64"
✓ gh_2.62.0_windows_amd64.zip (11274 KB)
```

## Cara kerja

- GET `repos/{owner}/{repo}/releases/latest` via GitHub API publik.
- Filter asset berdasarkan substring nama (opsional).
- Stream download ke file lokal.

Tanpa dependency eksternal — stdlib `net/http` saja.

## License

MIT
