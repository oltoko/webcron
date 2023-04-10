package cmd

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

const (
	INTERVAL                      = "interval"
	intervalDefault time.Duration = 5 * time.Second
)

// rootCmd represents the base command when called without any subcommands
var (
	intervalValue time.Duration
	cronEndpoint  string

	rootCmd = &cobra.Command{
		Use:   "webcron http:/...",
		Short: "A utility to call Webcron Jobs",
		Long: `A utility to call Webcron Jobs, so you don't have to
pay for a service in the Internet.`,
		Run: run_cron,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.PersistentFlags().DurationVarP(&intervalValue, INTERVAL, "i", intervalDefault, "The interval as Duration when the cron should be called.")
}

func run_cron(cmd *cobra.Command, args []string) {

	parse_endpoint(args)

	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGTERM)
	signal.Notify(signalC, syscall.SIGINT)
	go wait_for_term(signalC)

	s := gocron.NewScheduler(time.Local)
	s.Every(intervalValue).Do(do_request)
	s.StartAsync()
	log.Printf("Every %s calls %s\n", intervalValue, cronEndpoint)

	for {
		time.Sleep(1 * time.Second)
	}
}

func parse_endpoint(args []string) {

	if len(args) != 1 {
		log.Fatalln("The Address which should be called wasn't given!")
	}

	_, err := url.Parse(args[0])
	if err != nil {
		log.Fatalf("Failed to parse given URL %s\n", err)
	}

	cronEndpoint = args[0]
}

func do_request() {

	start := time.Now()
	res, err := http.Get(cronEndpoint)
	if err != nil {
		log.Printf("Failed to send Request: %s\n", err)
		return
	}
	defer res.Body.Close()
	end := time.Now()

	switch res.StatusCode {
	case 200:
		log.Printf("Request took %s\n", end.Sub(start))
	default:
		log.Printf("Unexpected Status Code: %d\n", res.StatusCode)
	}
}

func wait_for_term(signalC chan os.Signal) {
	<-signalC
	os.Exit(0)
}
