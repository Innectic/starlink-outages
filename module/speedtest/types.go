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

func (l LastData) AverageSpeeds() (int, int) {
	var totalUpload, totalDownload int = 0

	for speed, _ := range l.Speeds {
		totalUpload += speed.Upload
		totalDownload += speed.Download
	}

	return totalUpload / l.SuccessfulRuns, totalDownload / l.SuccessfulRuns
}

func (l LastData) Latency() (int, int, int) {
	var lowest, average, highest, total int = 0

	for speed, _ := range l.Speeds {
		total += speed.Latency

		if 
	}
}

func (l LastData) Failed() int {
	return l.TotalRuns - l.SuccessfulRuns
}

func (l *LastData) Result(upload, download int) {
	result := SpeedtestResult{
		Upload: upload,
		Download: download,
	}

	l.Speeds = append(l.Speeds, result)
}
