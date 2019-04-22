# Scratch

## Integration of Alerts

### Push vs. Pull

| .    | Advantages                                                                                                                                                                                                                                                                                                                                                 | Disadvantages                                                                                                                                                                                                                                                                        |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| pull | <ul><li>Interval/ number of incoming data is defined by sokar</li><li>Alertmanager is independent from SK</li><li>Complete configuration about alerts is on the SK side (service side).</li></ul>                                                                                                                                                          | <ul><li>SK has to know AWS SNS</li><li>SK has to know alertmanager</li><li>SK has to get incoming access to alertmanager</li><li>multiple SK's will scrape from alertmanager</li><li>Alerts are event based, a pull -> poll would lead to more traffic or less information</li></ul> |
| push | <ul><li>High resolution as possible, low resolution as needed (no needless calls).</li><li> Works with alertmanager (this is the intended way)</li><li>Works with AWS CW over AWS SNS</li><li>SK does not have to know alertmanager</li><li>Can be configured very specific (via routes). It can be defined which alerts are routed to which SK.</li></ul> | <ul><li>SK has (still) to know AWS SNS</li><li>alertmanager has to know all SK's (as alerting target)</li><li>Configuration is on alertmanger side (not on service-side).</li></ul>                                                                                                  |

- **Push wins**

## Local Setup

### Env Vars

```bash
export COS=<full path to the checked out cos repo> (see: https://github.com/MatthiasScholz/cos)
export SCRATCH=<full path to location for configuration/ testing files>
export AM=<full path to the checked out prometheus/alertmanager repo> (see: https://github.com/prometheus/alertmanager)
export PM=<full path to the checked out prometheus repo> (see: https://github.com/prometheus/prometheus)
export SK=<full path to the checked out sokar repo> (see: https://github.com/ThomasObenaus/sokar)
export LOCAL_IP=<ip of your host pc>

```

### Nomad

Start nomad locally

```bash
cd $COS/examples/devmode/ && ./devmode.sh $LOCAL_IP "public-services"
export NOMAD_ADDR=http://$LOCAL_IP:4646 && export CONSUL_HTTP_ADDR=http://LOCAL_IP:8500 && export IGRESS_ADDR=http://LOCAL_IP:9999

#Open ui
xdg-open $NOMAD_ADDR
```

Deploy fabio

- First inside `$COS/examples/devmode/fabio_docker.nomad` the `data_center` and `host-ip-address` have to be replaced with according values.

```bash
nomad run $COS/examples/devmode/fabio_docker.nomad
```

Deploy sample job

```bash
nomad run $SK/examples/multi-group.nomad
```

### Sokar

Build

```bash
cd $SK
make build
```

One Shot

```bash
$SK/sokar-bin -oneshot -nomad-server-address=$NOMAD_ADDR -job-name="fail-service" -scale-by=1
```

Deploy Sokar in Nomad (dev mode)

```bash
nomad run $SK/examples/sokar.nomad
```

### Alertmanager

Run it locally

```bash
# CMD
$AM/alertmanager --config.file=$SK/examples/alertmanager/config.yaml --log.level=debug

# Open ui
xdg-open http://localhost:9093
```

Fire alert

```bash
$AM/amtool alert --alertmanager.url=http://localhost:9093 add alertname=foo node=bar test=bla
```

Close all

```bash
pkill consul && sudo pkill nomad && pkill alertmanager
```

### Alerts via Curl

If fabio is deployed in a local setup (cos-devmode) the url would be: `http://${LOCAL_IP}:9999/sokar/api/alerts`.

```bash
curl -X POST \
  http://localhost:11000/api/alerts \
  -d '{
  "receiver": "PM",
  "status": "firing",
  "alerts": [
    {
      "status": "firings",
      "labels": {
        "alertname": "AlertA",
        "alert-type": "scaling",
        "scale-type": "up"
      },
      "annotations": {
        "description": "Scales the component XYZ UP"
      },
      "startsAt": "2019-02-23T12:00:00.000+01:00",
      "endsAt": "2019-02-23T12:05:00.000+01:00",
      "generatorURL": "http://generator_url"
    },
    {
      "status": "firings",
      "labels": {
        "alertname": "AlertB",
        "alert-type": "scaling",
        "scale-type": "down"
      },
      "annotations": {
        "description": "Scales the component XYZ DOWN"
      },
      "startsAt": "2019-02-23T12:00:00.000+01:00",
      "endsAt": "2019-02-23T12:05:00.000+01:00",
      "generatorURL": "http://generatorURL"
    }
  ],
  "groupLabels": {},
  "commonLabels": { "alertname": "AlertA" },
  "commonAnnotations": {},
  "externalURL": "http://externalURL",
  "version": "4",
  "groupKey": "{}:{}"
}
'
```

### Prometheus + Grafana

#### Start

```bash
cd $SK/examples/monitoring && docker-compose up -d

# open grafana ui
xdg-open http://localhost:3000

# open prometheus ui
xdg-open http://localhost:9090
```

#### Stop

```bash
cd $SK/examples/monitoring && docker-compose down
```
