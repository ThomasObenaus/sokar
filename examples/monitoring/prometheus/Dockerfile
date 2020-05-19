FROM  quay.io/prometheus/prometheus:v2.10.0

COPY config.yml /etc/prometheus/prometheus.yml
COPY entrypoint.sh /entrypoint.sh

USER root
RUN chmod u+x /entrypoint.sh
# Needed to ensure that there are no weird line breaks when code was checked out on windows
RUN dos2unix /entrypoint.sh

EXPOSE     9090

# Delegate the start up to entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml"]