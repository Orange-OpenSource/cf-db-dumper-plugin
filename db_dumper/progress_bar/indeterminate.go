package progress_bar

import "fmt"

type IndeterminateProgressBar struct {
	state          string
	loadingMessage string
}

func NewIndeterminateProgressBar(loadingMessage string) *IndeterminateProgressBar {
	ipb := &IndeterminateProgressBar{"|", loadingMessage}
	ipb.write()
	return ipb
}
func (this *IndeterminateProgressBar) Next() {
	var nextState string
	switch (this.state) {
	case "|":
		nextState = "/"
		break;
	case "/":
		nextState = "-"
		break;
	case "-":
		nextState = "\\"
		break;
	default:
		nextState = "|"
	}
	this.state = nextState
	this.write()
}
func (this *IndeterminateProgressBar) write() {
	fmt.Print(this.state + " " + this.loadingMessage + "\r")
}
