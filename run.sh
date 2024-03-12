#!/usr/bin/with-contenv bashio
set -e

mkdir -p /etc/wiegand2mqtt

cat > /etc/wiegand2mqtt/config.yaml <<EOF
mqtt:
  broker: $(bashio::config 'mqtt_ip')
  port: $(bashio::config 'mqtt_port')
  user: $(bashio::config 'mqtt_user')
  password: $(bashio::config 'mqtt_password')
loglevel: $(bashio::config 'loglevel')
EOF

/bin/wiegand2mqtt
