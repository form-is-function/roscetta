#!/bin/sh

# Initialize can0 interface
# Note: the `ip` binary from Busybox doesn't support CAN
ip link set can0 type can bitrate 500000
ip link set can0 up

# Start proxy
/opt/roscetta/proxy
