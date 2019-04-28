package colorful

const (
	colorOff    = string("\033[0m")
	colorRed    = string("\033[0;31m")
	colorGreen  = string("\033[0;32m")
	colorOrange = string("\033[0;33m")
	colorBlue   = string("\033[0;34m")
	colorPurple = string("\033[0;35m")
	colorCyan   = string("\033[0;36m")
	colorGray   = string("\033[0;37m")
)

func paint(data string, color string) string {
	return color + data + colorOff
}

// Red draw red on data
func Red(data string) string {
	return paint(data, colorRed)
}

// Green draw Green on data
func Green(data string) string {
	return paint(data, colorGreen)
}

// Orange draw Orange on data
func Orange(data string) string {
	return paint(data, colorOrange)
}

// Blue draw Blue on data
func Blue(data string) string {
	return paint(data, colorBlue)
}

// Purple draw Purple on data
func Purple(data string) string {
	return paint(data, colorPurple)
}

// Cyan draw Cyan on data
func Cyan(data string) string {
	return paint(data, colorCyan)
}

// Gray draw Gray on data
func Gray(data string) string {
	return paint(data, colorGray)
}
