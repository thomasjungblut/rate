# rate

Rate gives you the current rate (per second) of incoming lines over a pipe.

You can test with the checked-in bash script:

```
$ ./test_rate.sh | ./rate
Rate: 10.00/s
```


## Installation

With golang 1.12 (or higher) installed, just run

> go install github.com/thomasjungblut/rate

Make sure your PATH includes the $GOPATH/bin directory so your commands can be easily used:

> export PATH=$PATH:$GOPATH/bin
