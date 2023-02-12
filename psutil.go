package main

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func getMemoryUsage(ctx context.Context) chan float64 {
	ch := make(chan float64)
	go func(ctx context.Context) {
		tk := time.NewTicker(duration)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				v, _ := mem.VirtualMemory()
				ch <- v.UsedPercent
			}
		}
	}(ctx)
	return ch
}

func getCpuPercent(ctx context.Context) chan float64 {
	ch := make(chan float64)
	go func(ctx context.Context) {
		tk := time.NewTicker(duration)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				v, _ := cpu.PercentWithContext(ctx, 1*time.Second, false)
				ch <- v[0]
			}
		}
	}(ctx)
	return ch
}
