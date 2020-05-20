#!/usr/bin/env bash

HOSTNAME=$(hostname)

IP=$(grep "${HOSTNAME}" /etc/hosts | awk '{print $1}')

/usr/local/bin/imx -registry "${REGISTRY}" -rpc-endpoint "$IP":9321
