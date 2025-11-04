package utils

// Helper function to split content into lines
func SplitLines(content string) []string {
	lines := []string{}
	currentLine := ""
	for _, ch := range content {
		if ch == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else if ch != '\r' {
			currentLine += string(ch)
		}
	}
	lines = append(lines, currentLine)
	return lines
}
