# rate

Rate gives you the current rate (per second) of incoming lines over a pipe.

You can test with the checked-in bash script:

```
$ ./test_rate.sh | ./rate
Rate: 10.00/s
```

There are two additional options to visualize the rate, plot and table.

Plot looks like this:

```
$ ./test_rate.sh | ./rate --plot

37.4¦
    ¦                                                                               ••
    ¦                                                                             ••  ••
    ¦                                                                          •••      •
    ¦                                                                       •••          ••
    ¦                                                                     ••               •
    R                                                  •••••••••••••••••••                  •
    a                                            ••••••                                      ••
    t                                     •••••••                                              •
    e                              •••••••                                                      ••
    /                        ••••••                                                               •
    s                    ••••
    ¦                ••••
    ¦            ••••
    ¦         •••
    ¦     ••••
    ¦  •••
    ¦••
0   ¦---------------------------------Relative time in seconds--------------------------------------
     -35.0                                                                                       0.0

```

whereas table is a simple output of the past measurements:

```
$ ./test_rate.sh | ./rate --table

Time                               Count     Rate/s
2019-10-01 20:50:55 +0200 CEST     24        4.000000
2019-10-01 20:51:00 +0200 CEST     60        11.200000
2019-10-01 20:51:05 +0200 CEST     85        16.200000
2019-10-01 20:51:10 +0200 CEST     86        17.200000

```

## Installation

With golang 1.12 (or higher) installed, just run

> go install github.com/thomasjungblut/rate

Make sure your PATH includes the $GOPATH/bin directory so your commands can be easily used:

> export PATH=$PATH:$GOPATH/bin
