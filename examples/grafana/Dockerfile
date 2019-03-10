FROM  grafana/grafana:latest

#COPY config.yml /etc/prometheus/prometheus.yml

EXPOSE     3000
#VOLUME     [ "/prometheus" ]
#WORKDIR    /prometheus

USER grafana
WORKDIR /
ENTRYPOINT [ "/run.sh" ]
#CMD        [ "--config.file=/etc/prometheus/prometheus.yml"]


