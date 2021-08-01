package softwareupdate

import "fmt"

func GetMessage(time, old, new string) string {
	return fmt.Sprintf(`
‼️‼️ Starlink Update Detected ‼️‼️

When: %s

Old Version: %s
New Version: %s
	`, time, old, new)
}
