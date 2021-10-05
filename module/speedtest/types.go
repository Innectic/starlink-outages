package speedtest

type SpeedtestResult struct {
	Latency int
	Upload int
	Download int
}

type LastData struct {
	TotalRuns int
	SuccessfulRuns int

	Speeds []SpeedtestResult
}

func defaultLastData() LastData {
	return LastData{
		TotalRuns: 0,
		SuccessfulRuns: 0,
		Speeds: make([]SpeedtestResult, 0),
	}
}

func (l LastData) Download() (int, int, int) {
	var lowest, average, highest, total int = 0, 0, 0, 0

	for _, speed := range l.Speeds {
		total += speed.Download

		if speed.Download < lowest {
			lowest = speed.Download
		}

		if speed.Download > highest {
			highest = speed.Download
		}
	}

	average = total / len(l.Speeds)
	return average, lowest, highest
}

func (l LastData) Upload() (int, int, int) {
	var lowest, average, highest, total int = 0, 0, 0, 0

	for _, speed := range l.Speeds {
		total += speed.Upload

		if speed.Upload < lowest {
			lowest = speed.Upload
		}

		if speed.Upload > highest {
			highest = speed.Upload
		}
	}

	average = total / len(l.Speeds)
	return average, lowest, highest
}

func (l LastData) Latency() (int, int, int) {
	var lowest, average, highest, total int = 0, 0, 0, 0

	for _, speed := range l.Speeds {
		total += speed.Latency

		if speed.Latency < lowest {
			lowest = speed.Latency
		}

		if speed.Latency > highest {
			highest = speed.Latency
		}
	}

	average = total / len(l.Speeds)
	return average, lowest, highest
}

func (l LastData) Failed() int {
	return l.TotalRuns - l.SuccessfulRuns
}

func (l *LastData) Result(latency, upload, download int) {
	result := SpeedtestResult{
		Upload: upload,
		Download: download,
		Latency: latency,
	}

	l.Speeds = append(l.Speeds, result)
}

func (l *LastData) Reset() {
	l.TotalRuns = 0
	l.SuccessfulRuns = 0
	l.Speeds = make([]SpeedtestResult, 0)
}
