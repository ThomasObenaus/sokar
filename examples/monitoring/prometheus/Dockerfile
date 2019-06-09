FROM  quay.io/prometheus/prometheus:v2.10.0

COPY config.yml /etc/prometheus/prometheus.yml

EXPOSE     9090
#VOLUME     [ "/prometheus" ]
#WORKDIR    /prometheus
ENTRYPOINT [ "/bin/prometheus" ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml"]


