#!/bin/bash
set -e

envsubst '$$NGINX_HOST $$NGINX_PORT $$PROXY_PASS_HOST $$PROXY_PASS_PORT' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

exec "$@"
