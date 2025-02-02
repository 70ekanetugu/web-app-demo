FROM nginx:1.27.3-perl

COPY ./docker/nginx/templates/default.conf.template /etc/nginx/conf.d/
COPY ./docker/nginx/entrypoint.sh /usr/bin/

RUN chmod +x /usr/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]
