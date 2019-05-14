
import glob
import sys
import os

metric_types = ["NewCounter", "NewGauge"]
column_format = "|{:50s}|{:80s}|{:15s}|"


def isEndOfMetricDefinition(line):
    if "})" in line:
        return True
    else:
        return False


def isStartOfMetricDefinition(line):
    for metric_type in metric_types:
        if metric_type in line:
            return True
        else:
            return False


def extractMetricComponent(line):
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


def extractMetricType(line):
    if "NewCounter" in line:
        return "Counter"
    if "NewGauge" in line:
        return "Gauge"

    return "Unkown"


def createMetric(mBlock, metric_type):
    m = Metric()
    m.type = metric_type

    for element in mBlock:
        key, value = extractMetricComponent(element)
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


def metricToMD(metric):
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


def findMetricsFiles(dir):
    mFiles = list()
    for r, d, f in os.walk("."):
        for file in f:
            full_file_name = os.path.join(r, file)
            if "/metrics.go" in full_file_name and "vendor" not in full_file_name:
                mFiles.append(full_file_name)
    return mFiles


print("# Metrics")
print("")
# print metrics table header
print(column_format.format("Name", "Help", "Type"))
print("|{:->50s}|{:->80s}|{:->15s}|".format("", "", ""))

files = findMetricsFiles(".")
for file in files:
    file = open("sokar/metrics.go", "r")

    # read as a array, where each line is one element
    file_content = file.readlines()

    metrics = list()
    metric_block = list()
    metric_block_open = False
    metric_type = ""
    for i in range(0, len(file_content)):
        line = file_content[i]

        if isEndOfMetricDefinition(line) and metric_block_open:
            metric_block_open = False
            metric = createMetric(metric_block, metric_type)
            metrics.append(metric)

        if metric_block_open:
            metric_block.append(line)

        if isStartOfMetricDefinition(line):
            metric_type = extractMetricType(line)
            metric_block_open = True
            metric_block = list()

    # print metrics
    for metric in metrics:
        print(metricToMD(metric))
