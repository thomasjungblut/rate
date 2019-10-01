package cmd

import (
	"bufio"
	"fmt"
	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/spf13/cobra"
	"io"
	"os"
	"runtime"
	"time"
)

var (
	Plot                   bool
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

	ringBuffer := make([]int, StorageDurationSeconds/BucketSeconds)

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
	lastBucket := 0
	lastBucketCount := 0
	reader := bufio.NewReader(os.Stdin)
	lineCount := 0
	lastPrintedSeconds := 0
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

		nowSeconds := time.Now().Round(time.Second).Second()
		bucketSecond := nowSeconds % BucketSeconds
		diffPercent := float64(bucketSecond) / bucketSecondsFloat
		bucket := (nowSeconds - bucketSecond) % len(ringBuffer)
		if bucket != lastBucket {
			lastBucketCount = ringBuffer[lastBucket]
			ringBuffer[bucket] = 1
		} else {
			ringBuffer[bucket]++
		}
		lastBucket = bucket

		// only sample every 1k lines or after 2s
		if lineCount%1000 == 0 || (nowSeconds-lastPrintedSeconds) >= 2 {
			rate := 0.0
			if lastBucketCount != 0 {
				// this is doing a smart average over the last two buckets to smooth out incomplete windows.
				// it is based on how much of the window already passed and it is weighted based on that.
				rate = (1.0-diffPercent)*(float64(lastBucketCount)/bucketSecondsFloat) + diffPercent*float64(ringBuffer[bucket])/bucketSecondsFloat
			} else {
				rate = float64(ringBuffer[bucket]) / bucketSecondsFloat
			}

			fmt.Print(cursor.ClearEntireLine())
			fmt.Printf("\rRate: %.2f/s", rate)
			lastPrintedSeconds = nowSeconds
		}
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&Plot, "plot", "", false, "If set, it will plot a rate graph on the terminal")
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
