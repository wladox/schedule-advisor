package metrics

type Storage struct {
	Bandwidth          map[string]int
	CheckpointInterval int
	IngestionRate      int
	OperatorStats      map[string]int
}

func (s *Storage) GetOperatorStats(taskSubtaskId string) int {
	stats, ok := s.OperatorStats[taskSubtaskId]
	if !ok {
		return -1
	}
	return stats
}

func (s *Storage) GetIngestionRate() int {
	return s.IngestionRate
}

func (s *Storage) GetCheckpointInterval() int {
	return s.CheckpointInterval
}

func (s *Storage) GetBandwidth(link string) int {
	bw, ok := s.Bandwidth[link]
	if !ok {
		return -1
	}
	return bw
}
