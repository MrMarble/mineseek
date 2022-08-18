package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xrjr/mcutils/pkg/ping"

	"github.com/mrmarble/mineseek/internal/scanner"
	"github.com/mrmarble/mineseek/internal/utils"
)

var (
	// Populated by goreleaser during build.
	version = "master"
	commit  = "?"
	date    = ""
)

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong) error {
	fmt.Printf("mineseek has version %s built from %s on %s\n", version, commit, date)
	app.Exit(0)

	return nil
}

type Cli struct {
	Debug   bool        `help:"Enable debug mode."`
	Version VersionFlag `name:"version" help:"Print version information and quit"`

	Address string `help:"Address to scan." required:"" arg:""`

	Rate    int           `help:"Rate to scan at." default:"1"`
	Ports   string        `help:"Ports to scan. Can be a range or a list." default:"25565"`
	Timeout time.Duration `help:"Timeout for each scan." default:"1s"`

	Ping    bool `help:"Server List Ping." default:"false"`
	Stealth bool `help:"Stealth scan. Requires root." default:"false"`
	Verbose bool `help:"Verbose output. Show offline server" default:"false" short:"v"`
}

var (
	msg    = lipgloss.NewStyle().Foreground(lipgloss.Color("#77d3e6")).Inline(true).Render
	ip     = lipgloss.NewStyle().Foreground(lipgloss.Color("#a381ef")).Inline(true).Render
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("#a3e77d")).Inline(true).Render
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#e7d77d")).Inline(true).Render
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("#e77d7d")).Inline(true).Render
)

func main() {
	var cli Cli

	printLogo()

	// Kong hack to show help by default.
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&cli,
		kong.UsageOnError(),
		kong.Name("mineseek"),
		kong.Description("Minecraft server scanner. Running without arguments is basically a port scanner."),
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cli.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	ctx.FatalIfErrorf(run(&cli))
}

func run(cli *Cli) error {
	ports, err := utils.ParsePorts(cli.Ports)
	if err != nil {
		return err
	}

	addr := cli.Address
	if !strings.ContainsRune(addr, '/') {
		addr += "/32"
	}

	opts := scanner.Options{
		Host:     addr,
		Protocol: "tcp",
		Ports:    ports,
		Rate:     cli.Rate,
		Stealth:  cli.Stealth,
		Timeout:  cli.Timeout,
		SLP:      cli.Ping,
	}

	fmt.Println(msg("[~]"), "Scanning", addr, "with", cli.Rate, "threads")
	fmt.Println(msg("[~]"), "Hosts to scan:", utils.AvailableHosts(addr))

	results, err := scanner.Scan(&opts)
	if err != nil {
		return err
	}

	for result := range results {
		printServer(cli, result)
	}

	return nil
}

func printServer(cli *Cli, result scanner.Result) {
	log.Debug().Interface("result", result).Msg("Got result")

	if result.Open {
		if !cli.Verbose && result.Latency == -1 {
			return // Skip offline servers if not verbose.
		}

		fmt.Printf("Found %s", ip(result.Host))

		if cli.Ping {
			if result.SLP == nil || result.Latency == -1 {
				fmt.Printf(" (%s) \n\n", red("offline"))
			} else {
				infos := result.SLP.Infos()
				fmt.Printf("\n  Latency: %sms\n", printLatency(result.Latency))
				fmt.Printf("  Version: %s (protocol %d)\n", infos.Version.Name, infos.Version.Protocol)
				fmt.Printf("  Players: %d/%d %s\n",
					infos.Players.Online,
					infos.Players.Max,
					printPlayers(&infos))
				fmt.Printf("  MOTD: %s\n", infos.Description)
				fmt.Printf("\n" + msg("+====================+") + "\n\n\n")
			}
		}
	}
}

func printLatency(latency int) string {
	switch {
	case latency < 150:
		return green(strconv.Itoa(latency))
	case latency < 300:
		return yellow(strconv.Itoa(latency))
	default:
		return red(strconv.Itoa(latency))
	}
}

func printPlayers(infos *ping.Infos) string {
	sample := []string{}

	for _, player := range infos.Players.Sample {
		if player.ID != "00000000-0000-0000-0000-000000000000" {
			sample = append(sample, player.Name)
		}
	}

	if len(sample) > 0 {
		return "[" + strings.Join(sample, ", ") + "]"
	}

	return ""
}
