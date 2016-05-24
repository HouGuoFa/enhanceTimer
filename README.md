# about timer
timer is a golang packet, as high performance tick-time. Although the accuracy is not high(Resolution 1 second ) but when extensive use of the timer can improve performance.

# how to use
Use the same standard package (time packet and NewTimer function)

import ("timer")
func main() {
	tick := timer.NewTimer(1)
	defer tick.Stop()
	select {
	case <-tick.C:
	}
}
