<!--
DocName: README.md
DocType: README
Keywords: metrics
-->
# Metrics README

The Metrics package provides facilities to add metrics and telemetry to golang daemons.

Metric objects are lower-level primitives which can be either used directly or are internally used to construct
task-specific higher level abstractions which are easier to use but more restricted in their flexibility. Think
of them as raw materials vs fully assembled components.

Which one to use then? It is recommended to first start out figuring if any of the specific abstractions will fit
your use case and use that. For example `telemetry` is created specifically to provide a quick way to
add toggable per-function telemetry (ex. total function calls, time taken, statistics, etc.).

# Telemetry

The telemetry object is intended to provide quick per-function statistics:

- Calls
- Total time (ms)
- Average time (ms)

These can be accessed in the form of a JSON:

```json
[{"Function":"rdlabs.hpecorp.net/yang-resolver/pipeline/notification_service.(*NotificationService).Process","Calls":1,"TotalTimeMs":0,"AverageTimeMs":0},{"Function":"rdlabs.hpecorp.net/yang-resolver/pipeline/subscription_service.(*SubscriptionService).subscribeStreamToView","Calls":1,"TotalTimeMs":13,"AverageTimeMs":13},{"Function":"rdlabs.hpecorp.net/yang-resolver/pipeline/subscription_service.(*SubscriptionService).processSubscribe","Calls":1,"TotalTimeMs":13,"AverageTimeMs":13}]
```

To visualize it nicely you can pipe it to `jq`:

```json
[
  {
    "Function": "rdlabs.hpecorp.net/yang-resolver/pipeline/notification_service.(*NotificationService).Process",
    "Calls": 1,
    "TotalTimeMs": 0,
    "AverageTimeMs": 0
  },
  {
    "Function": "rdlabs.hpecorp.net/yang-resolver/pipeline/subscription_service.(*SubscriptionService).subscribeStreamToView",
    "Calls": 1,
    "TotalTimeMs": 13,
    "AverageTimeMs": 13
  },
  {
    "Function": "rdlabs.hpecorp.net/yang-resolver/pipeline/subscription_service.(*SubscriptionService).processSubscribe",
    "Calls": 1,
    "TotalTimeMs": 13,
    "AverageTimeMs": 13
  }
]
```

Or use the `telemetry_json_to_html.sh` script which creates a single, self-contained, HTML with charts
(sorted in ascending order):

```bash
corralda@hpnsw4014:~/workspace/halon/halon-src/go-common/go-metrics/scripts$ ./telemetry_json_to_html.sh -f ./output.json
Input JSON: ./output.json
Output HTML: /users/corralda/workspace/halon/halon-src/go-common/go-metrics/scripts/telemetry.html
```

The resulting `telemetry.html` file can be opened directly in a web browser.

## Telemetry example

```golang
    // Include the telemetry package
    include(
        "time"
        "math/rand"

        "rdlabs.hpecorp.net/metrics/telemetry"
    )

    func taskA() {
        if telemetry.IsEnabled() {
            defer telemetry.IncreaseFunctionTracer(telemetry.CurrentFunctionName(), time.Now())
        }

        // Simulate some work
        time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
    }

    func taskB() {
        if telemetry.IsEnabled() {
            defer telemetry.IncreaseFunctionTracer(telemetry.CurrentFunctionName(), time.Now())
        }

        // Simulate some work
        time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
    }

    func main() {
        // Enable the global metrics
        // Note: this should typically be done via a UnixCTL as
        //       they should start out disabled
        telemetry.Enable()

        // Simulate some work
        taskA()
        taskA()
        taskA()
        taskB()
        taskB()

        fmt.Println(telemetry.GetMetricsJSON())
    }
```

# Raw metrics

When using raw metric structures, you must first define a map that will contain the `name` of the metric, as well as its `type`. You can choose any `name` for a metric and for its `type` it can be Int or Float.

## Raw metrics example

Let's create two raw integer Metrics named `success` and `error`:

```golang
    // Include the package
    include(
        gometrics "rdlabs.hpecorp.net/metrics/"
    )

    var globalMetrics *gometrics.Metrics

    func onSuccess() {
        globalMetrics.IncreaseMetricValue("success", 1)
    }

    func onError() {
        globalMetrics.IncreaseMetricValue("error", 1)

    }

    func main() {
        // Define a map with the Metric names and types
        var metricsNameType = map[string]interface{}{
            "success": interface{}(0),
            "error": interface{}(0),
        }

        // Create the metrics
        globalMetrics = &gometrics.NewMetrics(metricsNameType)

        // Simulate some successes and errors
        onSuccess()
        onSuccess()
        onSuccess()
        onError()
        onError()

        // Read them one by one -- undefined metrics will result in err != nil
        successes, err := globalMetrics.ReadMetric("success")
        fmt.Println("dcl succeses: ", successes)
        errors, err := globalMetrics.ReadMetric("error")
        fmt.Println("dcl errors: ", errors)

        // A more straightforward approach is to get a map with all metrics
        // and iterate through them
        for metric, value := range globalMetrics.GetAllMetrics() {
            fmt.Println("dcl", metric, ": ", value)
        }
    }
```
