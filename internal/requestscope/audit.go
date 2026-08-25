package requestscope

import "sync"

type AuditSink struct {
	mu      sync.Mutex
	records []Snapshot
}

func (a *AuditSink) Record(snapshot Snapshot) {
	snapshot.Labels = append([]string(nil), snapshot.Labels...)
	a.mu.Lock()
	a.records = append(a.records, snapshot)
	a.mu.Unlock()
}

func (a *AuditSink) Records() []Snapshot {
	a.mu.Lock()
	defer a.mu.Unlock()
	result := make([]Snapshot, len(a.records))
	for i, record := range a.records {
		result[i] = record
		result[i].Labels = append([]string(nil), record.Labels...)
	}
	return result
}
