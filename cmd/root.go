package cmd

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/ahmetalpbalkan/go-cursor"
	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
	"io"
	"os"
	"runtime"
	"time"
)

type Bucket struct {
	bucketTime time.Time
	count      *uint64
	rate       *float64
}

var (
	Plot                   bool
	Table                  bool
	BucketSeconds          int
	StorageDurationSeconds int
)

var rootCmd = &cobra.Command{
	Use:   "rate",
	Short: "This tool will give you the rate of incoming lines on stdin.",
	Run:   awaitPipeCommands,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of rate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.1")
	},
}

func awaitPipeCommands(cmd *cobra.Command, args []string) {
	list := list.New()

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	isTerminal := (info.Mode()&os.ModeDevice == os.ModeDevice || runtime.GOOS == "windows") &&
		info.Mode()&os.ModeCharDevice == os.ModeCharDevice
	if isTerminal {
		fmt.Println("rate requires a pipe as input")
		return
	}

	bucketSecondsFloat := float64(BucketSeconds)
	reader := bufio.NewReader(os.Stdin)
	lineCount := 0
	lastPrintedTime := time.Now().Add(time.Duration(-5) * time.Second)
	for {
		_, isPrefix, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}

		lineCount++

		// we skip over partial lines
		if isPrefix {
			continue
		}

		now := time.Now()
		nowSecondResolution := now.Round(time.Second)
		nowBucketedTime := now.Round(time.Duration(BucketSeconds) * time.Second)
		diffPercent := float64(nowBucketedTime.Second()) / bucketSecondsFloat

		currentElement := list.Back()
		if currentElement == nil || currentElement.Value.(Bucket).bucketTime != nowBucketedTime {
			list.PushBack(Bucket{bucketTime: nowBucketedTime, count: new(uint64), rate: new(float64)})
			currentElement = list.Back()
		}

		bucket := currentElement.Value.(Bucket)
		*bucket.count++

		// only sample every 1k lines or after 2s
		if lineCount%1000 == 0 || (nowSecondResolution.Add(time.Duration(-2) * time.Second).After(lastPrintedTime)) {

			// cleanup old items
			cutoffTime := nowSecondResolution.Add(time.Duration(-StorageDurationSeconds) * time.Second)
			for list.Front().Value.(Bucket).bucketTime.Before(cutoffTime) {
				list.Remove(list.Front())
			}

			lastPrintedTime = nowSecondResolution
			*bucket.rate = float64(*bucket.count) / bucketSecondsFloat
			smoothedRate := 0.0

			if list.Len() > 1 {
				lastBucket := currentElement.Prev().Value.(Bucket)
				// this is doing a smart average over the last two buckets to smooth out incomplete buckets.
				// it is weighted on how much of the bucket already passed.
				smoothedRate = (1.0-diffPercent)*(float64(*lastBucket.count)/bucketSecondsFloat) + diffPercent*float64(*bucket.count)/bucketSecondsFloat
			}

			if Plot {
				tm.Clear()
				tm.MoveCursor(0, 0)

				chart := tm.NewLineChart(100, 20)

				data := new(tm.DataTable)
				data.AddColumn("Relative time in seconds")
				data.AddColumn("Rate/s")

				for e := list.Front(); e != nil; e = e.Next() {
					i := e.Value.(Bucket)
					data.AddRow(float64(i.bucketTime.Unix()-bucket.bucketTime.Unix()), *i.rate)
				}

				_, _ = tm.Println(chart.Draw(data))
				tm.Flush()
			} else if Table {
				tm.Clear()
				tm.MoveCursor(0, 0)

				totals := tm.NewTable(0, 10, 5, ' ', 0)
				_, _ = fmt.Fprintf(totals, "Time\tCount\tRate/s\n")
				for e := list.Front(); e != nil; e = e.Next() {
					i := e.Value.(Bucket)
					_, _ = fmt.Fprintf(totals, "%s\t%d\t%f\n", i.bucketTime.String(), *i.count, *i.rate)
				}
				_, _ = tm.Println(totals)
				tm.Flush()
			} else {
				fmt.Print(cursor.ClearEntireLine())
				if list.Len() > 1 {
					fmt.Printf("\rCurrent: %.2f/s, %ds weighted avg: %.2f/s", *bucket.rate, BucketSeconds*2, smoothedRate)
				} else {
					fmt.Printf("\rCurrent: %.2f/s", *bucket.rate)
				}
			}
		}
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&Plot, "plot", "", false, "If set, it will plot a rate graph on the terminal")
	rootCmd.Flags().BoolVarP(&Table, "table", "", false, "If set, it will display a rate table on the terminal")
	rootCmd.Flags().IntVarP(&BucketSeconds, "bucketDuration", "b", 5, "Duration of a measurement bucket in seconds, 5s by default")
	rootCmd.Flags().IntVarP(&StorageDurationSeconds, "storageDuration", "s", 60, "How long to keep the counters around in seconds, 60s by default.")

	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
