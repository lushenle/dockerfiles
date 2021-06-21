#!/bin/bash
#
#openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout dtls.key -out dtls.crt -subj "/CN=dashboard.idoocker.io"
kubectl create secret generic dashboard-tls --from-file=dtls.crt --from-file=dtls.key -n kube-system
