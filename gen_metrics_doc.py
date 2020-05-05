
import os

metric_types = ["NewCounter", "NewGauge",
                "NewWrappedGaugeVec", "NewWrappedCounterVec", "NewHistogram"]
column_format = "| {:<50s} | {:<220s} | {:<20s} |"


def is_end_of_metric_definition(line):
    if "})" in line or "}," in line:
        return True
    else:
        return False


def is_start_of_metric_definition(line):
    for metric_type in metric_types:
        #print(metric_type, "->", line)
        if metric_type in line:
            return True
        else:
            continue
    return False


def extract_metric_component(line):
    "Takes a line(string) like 'Name: \"scale_events_total\",' and splits them into a key value pair"
    elements = line.split(":")
    key = elements[0].strip()
    value = elements[1].strip()

    # remove trailing ,
    value = value.rstrip(",")
    # remove trailing "
    value = value.rstrip("\"")
    # remove leading "
    value = value.lstrip("\"")

    return (key, value)


def extract_metric_type(line):
    if "NewCounter" in line:
        return "Counter"
    if "NewGauge" in line:
        return "Gauge"
    if "NewWrappedGaugeVec" in line:
        return "Labelled Gauge"
    if "NewWrappedCounterVec" in line:
        return "Labelled Counter"
    if "NewHistogram" in line:
        return "Histogram"
    return "Unkown"


def create_metric(m_block, metric_type):
    m = Metric()
    m.type = metric_type

    for element in m_block:
        key, value = extract_metric_component(element)
        if key == "Name":
            m.name = value
        elif key == "Namespace":
            m.namespace = value
        elif key == "Help":
            m.help = value
        elif key == "Subsystem":
            m.subsystem = value
        else:
            print("Error: key '", key, "' not known")
    return m


def metric_to_md(metric):
    full_name = ""
    if len(metric.namespace) > 0:
        full_name += metric.namespace + "_"
    if len(metric.subsystem) > 0:
        full_name += metric.subsystem + "_"
    full_name += metric.name

    return column_format.format(full_name, metric.help, metric.type)


class Metric:
    def __init__(self):
        self.type = ""
        self.name = ""
        self.help = ""
        self.namespace = ""
        self.subsystem = ""
        self.buckets = None


def find_metrics_files(directory):
    m_files = list()
    for r, _, f in os.walk(directory):
        for file in f:
            full_file_name = os.path.join(r, file)
            if "/metrics.go" in full_file_name and "vendor" not in full_file_name:
                m_files.append(full_file_name)
    return m_files


print("# Metrics")
print("")
# print metrics table header
print(column_format.format("Name", "Help", "Type"))
print("| {:-<50s} | {:-<220s} | {:-<20s} |".format(":", ":", ":"))

m_files = find_metrics_files(".")
for m_file in m_files:
    file = open(m_file, "r")

    # read as a array, where each line is one element
    file_content = file.readlines()

    metrics = list()
    metric_block = list()
    metric_block_open = False
    metric_type = ""
    for i, line in enumerate(file_content):
        if is_end_of_metric_definition(line) and metric_block_open:
            metric_block_open = False
            metric = create_metric(metric_block, metric_type)
            metrics.append(metric)

        if metric_block_open:
            metric_block.append(line)

        if is_start_of_metric_definition(line):
            metric_type = extract_metric_type(line)
            metric_block_open = True
            metric_block = list()

    # write the metrics to file
    for metric in metrics:
        print(metric_to_md(metric))
