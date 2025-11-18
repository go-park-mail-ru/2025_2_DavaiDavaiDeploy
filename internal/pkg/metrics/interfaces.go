package metrics

type MetricsHTTP interface {
	IncreaseHits(string, string)
	IncreaseErr(string, string, string)
}
