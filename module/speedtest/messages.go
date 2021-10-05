
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
