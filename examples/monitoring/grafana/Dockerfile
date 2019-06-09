FROM  grafana/grafana:6.2.2

COPY grafana.ini /etc/grafana/grafana.ini
COPY prometheus_ds.yaml /etc/grafana/provisioning/datasources
COPY dashboard_config.yaml   /etc/grafana/provisioning/dashboards/
COPY dashboards/ /var/lib/grafana/dashboards/

EXPOSE     3000

USER grafana
WORKDIR /
ENTRYPOINT [ "/run.sh" ]

