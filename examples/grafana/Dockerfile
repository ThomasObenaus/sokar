FROM  grafana/grafana:latest

COPY grafana.ini /etc/grafana/grafana.ini
COPY prometheus_ds.yaml /etc/grafana/provisioning/datasources

EXPOSE     3000
#VOLUME     [ "/prometheus" ]
#WORKDIR    /prometheus

USER grafana
WORKDIR /
ENTRYPOINT [ "/run.sh" ]
#CMD        [ "--config.file=/etc/prometheus/prometheus.yml"]


