# Intent API devenv

This devenv stands up a `kube-apiserver` and a single-host `etcd` for local development / testing.

## kubeapi and etcd versions

You can switch to different kubeapi and etcd versions by changing the values in `.env`. Make sure that the images for those versions exist and can be pulled.

## Generating service account key and certificate

`kube-apiserver` requires TLS key and certificate in order to boot up. Those are used to authenticate requests from clients.

In order to generate those, you can use GNU openssl like so:

```sh
$ cd devenv/docker/blocks/intentapi
$ openssl req -new -newkey rsa:4096 -nodes -keyout sa-key.pem -out sa-csr.pem
$ openssl x509 -req -sha256 -in sa-csr.pem -signkey sa-key.pem -out sa-cert.pem
```
