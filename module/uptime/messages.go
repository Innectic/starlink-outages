package uptime

import "fmt"

func GetMessage(startTime string, duration int64, cause DishyStatus, downtimeEvents int, downtimeTotal int64, extra map[string]interface{}) string {
	switch cause {
	case DishyNoSatellites:
		return fmt.Sprintf(`
‼️‼️ Starlink Outage Detected ‼️‼️

When: %s
Outage Duration: %d seconds
Outage Cause: %s
Time waited until next satellite: %s seconds

Amount of downtime events today: %d
Total downtime today: %d
		`, startTime, duration, cause, extra["time_until_next_satellite"], downtimeEvents, downtimeTotal)
	default:
		return fmt.Sprintf(`
‼️‼️ Starlink Outage Detected ‼️‼️

When: %s
Outage Duration: %d seconds
Outage Cause: %s

Amount of downtime events today: %d
Total downtime today: %d seconds
		`, startTime, duration, cause, downtimeEvents, downtimeTotal)
	}
}
