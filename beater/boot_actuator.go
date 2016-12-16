package beater

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const METRICS_STATS = "/metrics"
const HEALTH_STATS = "/health"
const APP_INFO = "/info"

type HealthStats struct {
    Status string `json:"status"`
    DiskSpace struct {
        Status string `json:"status"`
        Total uint64 `json:"total"`
        Free uint64 `json:"free"`
        Threshold uint64 `json:"threshold"`
    } `json:"diskSpace"`
    DB struct {
        Status string `json:"status"`
        Database string `json:"database"`
        Hello uint64 `json:"hello"`
    } `json:"db"`
}

type MetricsStats struct {
    Mem struct {
        Total uint64 `json:"total"`
        Free uint64 `json:"free"`
    } `json:"mem"`
    Processors uint64 `json:"processors"`
    LoadAverage float64 `json:"load_average"`
    Uptime struct {
        Total uint64 `json:"total"`
        Instance uint64 `json:"instance"`
    } `json:"uptime"`
    Heap struct {
        Total uint64 `json:"total"`
        Committed uint64 `json:"committed"`
        Init uint64 `json:"init"`
        Used uint64 `json:"used"`
    } `json:"heap"`
    NonHeap struct {
        Total uint64 `json:"total"`
        Committed uint64 `json:"committed"`
        Init uint64 `json:"init"`
        Used uint64 `json:"used"`
    } `json:"non_heap"`
    Threads struct {
        Total uint64 `json:"total"`
        TotalStarted uint64 `json:"started"`
        Peak uint64 `json:"peak"`
        Daemon uint64 `json:"daemon"`
    } `json:"non_heap"`
    Classes struct {
        Total uint64 `json:"total"`
        Loaded uint64 `json:"loaded"`
        Unloaded uint64 `json:"unloaded"`
    } `json:"classes"`
    GC struct {
        Scavenge struct {
            Count uint64 `json:"count"`
            Time uint64 `json:"time"`
        } `json:"scavenge"`
        Marksweep struct {
            Count uint64 `json:"count"`
            Time uint64 `json:"time"`
        } `json:"marksweep"`
    } `json:"gc"`
    Http struct {
        SessionsMax int64 `json:"max_sessions"`
        SessionsActive uint64 `json:"active_sessions"`
    } `json:"http"`
    DataSource struct {
        PrimaryActive uint64 `json:"primary_active"`
        PrimaryUsage float64 `json:"primary_usage"`
    } `json:"data_source"`
    ResponseTime map[string]float64 `json:"response_time"`
    StatusCount map[string]map[string]float64 `json:"status_count"`
}

type RawMetricsStats struct {
    Mem uint64 `json:"mem"`
    MemFree uint64 `json:"mem.free"`
    Processors uint64 `json:"processors"`
    InstanceUptime uint64 `json:"instance.uptime"`
    Uptime uint64 `json:"uptime"`
    SystemloadAverage float64 `json:"systemload.average"`
    HeapCommitted uint64 `json:"heap.committed"`
    HeapInit uint64 `json:"heap.init"`
    HeapUsed uint64 `json:"heap.used"`
    Heap uint64 `json:"heap"`
    NonheapCommitted uint64 `json:"nonheap.committed"`
    NonheapInit uint64 `json:"nonheap.init"`
    NonheapUsed uint64 `json:"nonheap.used"`
    Nonheap uint64 `json:"nonheap"`
    ThreadsPeak uint64 `json:"threads.peak"`
    ThreadsDaemon uint64 `json:"threads.daemon"`
    ThreadsTotalStarted uint64 `json:"threads.totalStarted"`
    Threads uint64 `json:"threads"`
    Classes uint64 `json:"classes"`
    ClassesLoaded uint64 `json:"classes.loaded"`
    ClassesUnloaded uint64 `json:"classes.unloaded"`
    GCPsScavengeCount uint64 `json:"gc.ps_scavenge.count"`
    GCPsScavengeTime uint64 `json:"gc.ps_scavenge.time"`
    GCPsMarksweepCount uint64 `json:"gc.ps_marksweep.count"`
    GCPsMarksweepTime uint64 `json:"gc.ps_marksweep.time"`
    HttpSessionsMax int64 `json:"httpsessions.max"`
    HttpSessionsActive uint64 `json:"httpsessions.active"`
    DateSourcePrimaryActive uint64 `json:"datasource.primary.active"`
    DateSourcePrimaryUsage float64 `json:"datasource.primary.usage"`
}

type AppInfo struct {
  App struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Port string `json:"port"`
    Environment string `json:"environment"`
  } `json:"app"`
}

func (bt *Springbeat) GetHealthStats(u url.URL) (*HealthStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + HEALTH_STATS)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP%s", res.Status)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    stats := &HealthStats{}
    err = json.Unmarshal([]byte(body), &stats)
    if err != nil {
        return nil, err
    }

    return stats, nil
}

func (bt *Springbeat) GetMetricsStats(u url.URL) (*MetricsStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + METRICS_STATS)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP%s", res.Status)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    raw_stats := &RawMetricsStats{}
    err = json.Unmarshal([]byte(body), &raw_stats)
    if err != nil {
        return nil, err
    }

    raw_json := make(map[string]interface{})
    err = json.Unmarshal([]byte(body), &raw_json)
    if err != nil {
        return nil, err
    }

    keys := make([]string, len(raw_json))
    i := 0
    for k := range raw_json {
        keys[i] = k
        i++
    }

    stats := &MetricsStats{}
    stats.Mem.Free = raw_stats.MemFree
    stats.Mem.Total = raw_stats.Mem
    stats.Processors = raw_stats.Processors
    stats.LoadAverage = raw_stats.SystemloadAverage
    stats.Uptime.Total = raw_stats.Uptime
    stats.Uptime.Instance = raw_stats.InstanceUptime
    stats.Heap.Total = raw_stats.Heap
    stats.Heap.Init = raw_stats.HeapInit
    stats.Heap.Committed = raw_stats.HeapCommitted
    stats.Heap.Used = raw_stats.HeapUsed
    stats.NonHeap.Total = raw_stats.Nonheap
    stats.NonHeap.Init = raw_stats.NonheapInit
    stats.NonHeap.Committed = raw_stats.NonheapCommitted
    stats.NonHeap.Used = raw_stats.NonheapUsed
    stats.Threads.Total = raw_stats.Threads
    stats.Threads.TotalStarted = raw_stats.ThreadsTotalStarted
    stats.Threads.Peak = raw_stats.ThreadsPeak
    stats.Threads.Daemon = raw_stats.ThreadsDaemon
    stats.Classes.Total = raw_stats.Classes
    stats.Classes.Loaded = raw_stats.ClassesLoaded
    stats.Classes.Unloaded = raw_stats.ClassesUnloaded
    stats.GC.Scavenge.Count = raw_stats.GCPsScavengeCount
    stats.GC.Scavenge.Time = raw_stats.GCPsScavengeTime
    stats.GC.Marksweep.Count = raw_stats.GCPsMarksweepCount
    stats.GC.Marksweep.Time = raw_stats.GCPsMarksweepTime
    stats.Http.SessionsActive = raw_stats.HttpSessionsActive
    stats.Http.SessionsMax = raw_stats.HttpSessionsMax
    stats.DataSource.PrimaryActive = raw_stats.DateSourcePrimaryActive
    stats.DataSource.PrimaryUsage = raw_stats.DateSourcePrimaryUsage
    stats.ResponseTime = make(map[string]float64)
    stats.StatusCount = make(map[string]map[string]float64)

    // Dynamicly add counters and response times
    for _, k := range keys {
        if strings.HasPrefix(k, "counter.status.") {
            suffix := strings.TrimPrefix(k, "counter.status.");
            items := strings.SplitN(suffix, ".", 2);
            statusMap := stats.StatusCount[items[0]];
            if statusMap == nil {
                statusMap = make(map[string]float64);
            }
            statusMap[strings.Replace(items[1], ".", "_", -1)] = raw_json[k].(float64);
            stats.StatusCount[items[0]] = statusMap;
        }
        if strings.HasPrefix(k, "gauge.response.") {
            suffix := strings.TrimPrefix(k, "gauge.response.");
            suffix = strings.Replace(suffix, ".", "_", -1);
            stats.ResponseTime[suffix] = raw_json[k].(float64);
        }
    }

    return stats, nil
}

func (bt *Springbeat) GetAppInfo(u url.URL) (*AppInfo, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + APP_INFO)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP%s", res.Status)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    stats := &AppInfo{}
    err = json.Unmarshal([]byte(body), &stats)
    if err != nil {
        return nil, err
    }
    return stats, nil
}
