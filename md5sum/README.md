## Get fm-package MD5

- build 

```bash
$ GO111MODULE=on go build -o pkgmd5
```

- fm-package md5

```bash
$ pkgmd5 -d ~/workspace/fm-package
```

- directory md5

```bash
$ pkgmd5 -d ~/workspace/fm-package/fingerprint-web -merge
```

- single file

```bash
$ pkgmd5 -f ~/workspace/fm-package/backend/backend.jar
```

