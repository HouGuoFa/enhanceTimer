## About timer

- timer is a standard golang packet, with high performance tick-timer, while the accuracy is not high(Resolution 1s). By using 'timer' in this package, the performance can be improved much.

## How to use

- Importing the same standard package (time packet and NewTimer function)
- Add these source files to your program, named by enhanceTimer or something

## Example
    
    import (
        "timer"  //the standard golang package;
        "enhanceTimer" //or something you named this package, you can simply add these files to your GOPATH without import at the same time;
    )

    func main() {

        tick := timer.NewTimer(1)
        defer tick.Stop()
    
        select {
            case <-tick.C:
            // todo something
        }
    }
   