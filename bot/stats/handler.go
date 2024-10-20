package stats

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"math"
	"runtime"
	"time"
)

var uptime = time.Now()

func Handler(s *dgo.Session, i *dgo.MessageCreate) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	embed := &dgo.MessageEmbed{
		Type: dgo.EmbedTypeRich,
		Fields: []*dgo.MessageEmbedField{
			{
				Name:  "Total allocated memory (ever allocated for heap objects)",
				Value: formatBytes(m.TotalAlloc),
			},
			{
				Name:  "Total memory obtained from the OS (includes heap, stack, etc.)",
				Value: formatBytes(m.Sys),
			},
			{
				Name:  "Currently allocated heap memory",
				Value: formatBytes(m.Alloc),
			},
			{
				Name:  "Heap memory reserved but not currently used",
				Value: formatBytes(m.HeapIdle),
			},
			{
				Name:  "Heap memory in-use",
				Value: formatBytes(m.HeapInuse),
			},
			{
				Name:  "Stack memory in-use",
				Value: formatBytes(m.StackInuse),
			},
			{
				Name:  "Memory obtained from system via mmap (span and cache)",
				Value: formatBytes(m.MSpanSys + m.MCacheSys),
			},
			{
				Name:  "Memory used for GC metadata",
				Value: formatBytes(m.GCSys),
			},
			{
				Name:  "Uptime",
				Value: fmt.Sprintf("<t:%d:R>", uptime.UTC().Unix()),
			},
		},
	}

	_, _ = s.ChannelMessageSendEmbed(i.ChannelID, embed)
}

func formatBytes(s uint64) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanizeBytes(s, 1000, sizes)
}

func humanizeBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
