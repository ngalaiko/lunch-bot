## Local certificates

* Step: 1
  Install mkcert tool - macOS; you can see the mkcert repo for details
```
$ brew install mkcert
```

* Step: 2
  Install nss (only needed if you use Firefox)
```
$ brew install nss
```

* Step: 3
  Setup mkcert on your machine (creates a CA)
```
$ mkcert -install
```

* Step: 4 (Final)
  at the project root directory run the following command
```
$ mkdir -p .cert && mkcert -key-file ./.cert/key.pem -cert-file ./.cert/cert.pem 'localhost'
```
