
package speedtest

import (
	"fmt"
	"time"
)

func EachHour(latency, upload, download int) string {
	return fmt.Sprintf(`
‼️‼️ Starlink Speed Test ‼️‼️

Test taken at: %s
Latency: %d
Upload Speed: %d
Download Speed: %d
	`, time.Now().Format("01/02/2006 15:04:05 MST"), latency, upload, download)
}

func Daily(downloadAvg, downloadLow, downloadHigh, uploadAvg, uploadLow, uploadHigh, latencyAvg, latencyLow, latencyHigh, totalTests, failedTests int) string {
	return fmt.Sprintf(`
‼️‼️ Starlink Daily Speed Results ‼️‼️

Latency (Average, Highest, Lowest): %d, %d, %d
Upload (Average, Highest, Lowest): %d, %d, %d
Download (Average, Highest, Lowest): %d, %d, %d

Total Tests: %d
Failed Tests: %d
	`, latencyAvg, latencyHigh, latencyLow, uploadAvg, uploadHigh, uploadLow, downloadAvg, downloadHigh, downloadLow, totalTests, failedTests)
}
