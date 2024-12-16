package monitoring

import (
    "time"
    "sync"
)

type Metrics struct {
    mutex sync.RWMutex
    requestCount    int64
    responseTime    []time.Duration
    cacheHits      int64
    cacheMisses    int64
    errors         int64
}

var metrics = &Metrics{
    responseTime: make([]time.Duration, 0),
}

func RecordRequest() {
    metrics.mutex.Lock()
    defer metrics.mutex.Unlock()
    metrics.requestCount++
}

func RecordResponseTime(duration time.Duration) {
    metrics.mutex.Lock()
    defer metrics.mutex.Unlock()
    metrics.responseTime = append(metrics.responseTime, duration)
}

func RecordCacheHit() {
    metrics.mutex.Lock()
    defer metrics.mutex.Unlock()
    metrics.cacheHits++
}

func RecordCacheMiss() {
    metrics.mutex.Lock()
    defer metrics.mutex.Unlock()
    metrics.cacheMisses++
}

func RecordError() {
    metrics.mutex.Lock()
    defer metrics.mutex.Unlock()
    metrics.errors++
}

func GetMetrics() map[string]interface{} {
    metrics.mutex.RLock()
    defer metrics.mutex.RUnlock()
    
    var avgResponseTime time.Duration
    if len(metrics.responseTime) > 0 {
        total := time.Duration(0)
        for _, t := range metrics.responseTime {
            total += t
        }
        avgResponseTime = total / time.Duration(len(metrics.responseTime))
    }

    return map[string]interface{}{
        "total_requests": metrics.requestCount,
        "avg_response_time_ms": avgResponseTime.Milliseconds(),
        "cache_hits": metrics.cacheHits,
        "cache_misses": metrics.cacheMisses,
        "errors": metrics.errors,
    }
}