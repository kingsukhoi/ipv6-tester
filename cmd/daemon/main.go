package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/kingsukhoi/ipv6-tester/pkg/ipv6Info"
)

func main() {
	hostname, _ := os.Hostname()

	if hostname == "" {
		hostname = "localhost"
	}

	victoriaEndpoint, exists := os.LookupEnv("VICTORIA_ENDPOINT")
	if !exists {
		panic("VICTORIA_ENDPOINT environment variable not set")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := runTest(ctx, ipv6Info.GoogleIpV6Addr)
	if err != nil {
		slog.Error("runTest failed", "error", err)
	}

	err = metrics.PushMetrics(ctx, victoriaEndpoint, false,
		&metrics.PushOptions{
			ExtraLabels: fmt.Sprintf(`hostname="%s"`, hostname),
		})
	if err != nil {
		slog.Error("metrics push failed", "error", err)
	}
	slog.Info("metrics push finished")

	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-ticker.C:
			errL := runTest(ctx, ipv6Info.GoogleIpV6Addr)
			if errL != nil {
				slog.Error("runTest failed", "error", errL)
			}

			errL = metrics.PushMetrics(ctx, victoriaEndpoint, false,
				&metrics.PushOptions{
					ExtraLabels: fmt.Sprintf(`hostname="%s"`, hostname),
				})
			if errL != nil {
				slog.Error("metrics push failed", "error", errL)
			}
			slog.Info("metrics push finished")
		case <-ctx.Done():
			return
		}
	}

}
func runTest(ctx context.Context, serverToTest string) error {
	connRes, err := ipv6Info.TestIpV6Connection(ctx, serverToTest)
	if err != nil {
		slog.Error("error testing ipv6 connection", "error", err, "server", serverToTest)
	} else {
		slog.Info("connection test", "server", connRes.ServerAddress, "success", connRes.Success, "localIp", connRes.LocalAddress)
	}
	var successFloat float64
	if connRes.Success {
		successFloat = 1
	} else {
		successFloat = 0
	}

	ipv6WorkingGauge := metrics.GetOrCreateGauge(fmt.Sprintf(`ipv6_working{server="%s"}`,
		connRes.ServerAddress), nil)
	ipv6WorkingGauge.Set(successFloat)

	if err != nil && connRes.LocalAddress == "" {
		addressGauge := metrics.GetOrCreateGauge(fmt.Sprintf(`ipv6_addresses{addresses="%s", server="%s"}`,
			"", serverToTest), nil)
		addressGauge.Set(0)
		return err
	}

	currInterface, err := ipv6Info.GetInterfaceByIP(connRes.LocalAddress)
	if err != nil {
		return err
	}
	slog.Info("address for server", "server", connRes.ServerAddress, "local_address", connRes.LocalAddress,
		"interface", currInterface.Name)

	ipAddrs, err := ipv6Info.GetLocalIpv6Addresses(ctx, currInterface.Name)
	if err != nil {
		return fmt.Errorf("getting local ipv6 addresses: %w", err)
	}

	slog.Info("current local ipv6 addresses", "addresses", ipAddrs.Addresses, "interface", currInterface.Name)

	addressGauge := metrics.GetOrCreateGauge(fmt.Sprintf(`ipv6_addresses{addresses="%s", server="%s"}`,
		strings.Join(ipAddrs.Addresses, ","), connRes.ServerAddress), nil)
	addressGauge.Set(float64(len(ipAddrs.Addresses)))

	return nil
}
