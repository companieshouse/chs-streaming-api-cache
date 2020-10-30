#!/bin/bash
#
# Start script for chs-streaming-api-cache

APP_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ -z "${MESOS_SLAVE_PID}" ]]; then
    source ~/.chs_env/private_env
    source ~/.chs_env/global_env
    source ~/.chs_env/chs-streaming-api-cache/env

    PORT="${CHS_STREAMING_API_CACHE_PORT:=6001}"
else
    PORT="$1"
    CONFIG_URL="$2"
    ENVIRONMENT="$3"
    APP_NAME="$4"

    source /etc/profile

    echo "Downloading environment from: ${CONFIG_URL}/${ENVIRONMENT}/${APP_NAME}"
    wget -O "${APP_DIR}/private_env" "${CONFIG_URL}/${ENVIRONMENT}/private_env"
    wget -O "${APP_DIR}/global_env" "${CONFIG_URL}/${ENVIRONMENT}/global_env"
    wget -O "${APP_DIR}/app_env" "${CONFIG_URL}/${ENVIRONMENT}/${APP_NAME}/env"
    source "${APP_DIR}/private_env"
    source "${APP_DIR}/global_env"
    source "${APP_DIR}/app_env"
fi

# Read brokers from environment and split on comma
IFS=',' read -ra BROKERS <<< "${KAFKA_STREAMING_BROKER_ADDR}"

# Ensure we only populate the broker address via application arguments
unset KAFKA_STREAMING_BROKER_ADDR

exec "${APP_DIR}/chs-streaming-api-cache" "-bind-address=:${PORT}" $(for broker in "${BROKERS[@]}"; do echo -n "-cache-broker-url=${broker} "; done)
