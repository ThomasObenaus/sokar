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
export AM=<full path to the checked out prometheus/alertmanager repo> (see: [https://github.com/MatthiasScholz/cos](https://github.com/prometheus/alertmanager))
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
