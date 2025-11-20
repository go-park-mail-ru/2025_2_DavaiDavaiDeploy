package metrics

type MetricsHTTP interface {
	IncreaseHits(string)
	IncreaseErr(string)
	ObserveResponseTime(int, string, float64)
}

type MetricsGrpc interface {
	IncreaseHits(string)
	IncreaseErr(string)
	ObserveResponseTime(int, string, float64)
}
