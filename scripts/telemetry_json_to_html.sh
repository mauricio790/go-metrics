#!/bin/bash

# (c) Copyright 2022 Hewlett Packard Enterprise Development LP
#
# Confidential computer software. Valid license from Hewlett Packard
# Enterprise required for possession, use or copying.
#
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items
# are licensed to the U.S. Government under vendor's standard commercial
# license.

# This is the template for the main HTML and contains all charts
HTML_TEMPLATE='
<!DOCTYPE html>
<script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.9.1/chart.min.js"></script>
<html>
<body>
<h1>Telemetry metrics</h1>
CHARTS
</body>
</html>
'

HORIZONTAL_BAR_CHART_TEMPLATE="
<h2>METRIC</h2>
<canvas id="METRIC" style='width:100%'></canvas>

<script>
var telemetry_metrics = SORTED_JSON;

var labels = telemetry_metrics.map(function(e) {
   return e.Function;
});
var data = telemetry_metrics.map(function(e) {
   return e.METRIC;
});

new Chart('METRIC', {
  type: 'bar',
  data: {
    labels: labels,
		datasets: [
			{
				data: data,
				backgroundColor: 'rgba(0, 119, 204, 0.3)'
			}
		]
  },
  options: {
		indexAxis: 'y',
		skipNull: true,
    plugins: { legend: {display: false} }
  }
});
</script>"

# Check arguments
while getopts u:a:f: flag
do
    case "${flag}" in
        f) filename=${OPTARG};;
    esac
done

if [ -z "${filename}" ]; then
  echo "Usage: $0 -f <path_to_input_json>"
  exit -1
fi

# Got the filename at this point
echo "Input JSON: ${filename}"

# Now sort the JSON for every metric (ascending) and create a corresponding chart
CHARTS=""
for METRIC in TotalTimeMs AverageTimeMs Calls; do
  SORTED_JSON=`cat ${filename} | jq 'sort_by(.'${METRIC}') | reverse'`

  # Use string replacement (ex. METRIC --> ${METRIC)})
  # Other types of evaluation are... painful :)
  CHART=${HORIZONTAL_BAR_CHART_TEMPLATE//METRIC/${METRIC}}
  CHART=${CHART//SORTED_JSON/${SORTED_JSON}}

  # Add this chart to the main charts list
  CHARTS="${CHARTS}${CHART}"
done

# Output the final HTML file
OUTPUT_HTML=`pwd`"/telemetry.html"
echo "Output HTML: ${OUTPUT_HTML}"
echo ${HTML_TEMPLATE//CHARTS/${CHARTS}} > ${OUTPUT_HTML}