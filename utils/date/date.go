package date

import "time"

func TimeStampSlug(timeStamp time.Time) string {
	date := timeStamp.Format("2006-01-02")
	hours := timeStamp.Format("15")
	minutes := timeStamp.Format("04")
	seconds := timeStamp.Format("05")
	return date + "_" + hours + "-" + minutes + "-" + seconds
}

// Time format to US format as string
func TimeStampToUSFormat(timestamp time.Time) string {
	// Use Format to convert to human-readable US format (January 02, 2006 15:04:05)
	usFormat := timestamp.Format("January 02, 2006 15:04:05")
	return usFormat
}

func TimeStamp() time.Time {
	return time.Now()
}
