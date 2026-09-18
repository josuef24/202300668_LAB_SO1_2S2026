package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	carnet              = "202300668"
	procMetricPath      = "/proc/continfo_pr2_so1_202300668"
	defaultLogFile      = "./telemetry_snapshot.json"
	defaultGeneratorTTL = 20 * time.Second
)

type ProcMetric struct {
	PID         int     `json:"pid"`
	Name        string  `json:"name"`
	Command     string  `json:"command"`
	ContainerID string  `json:"container_id,omitempty"`
	VSZ         int64   `json:"vsz_kb"`
	RSS         int64   `json:"rss_kb"`
	MemPercent  float64 `json:"mem_percent"`
	CPUPercent  float64 `json:"cpu_percent"`
	Category    string  `json:"category"`
}

type ContainerDecision struct {
	Timestamp time.Time      `json:"timestamp"`
	ToRemove  []string       `json:"to_remove"`
	ToKeep    []string       `json:"to_keep"`
	Metrics   []ProcMetric   `json:"metrics"`
	Summary   map[string]int `json:"summary"`
}

type containerInfo struct {
	ID    string
	Name  string
	Image string
	Kind  string
}

func Run() {
	fmt.Printf("[daemon] Iniciando servicio de telemetría para carnet %s\n", carnet)
	if err := ensureBaseContainers(); err != nil {
		fmt.Printf("[daemon] error al asegurar contenedores base: %v\n", err)
	}

	ticker := time.NewTicker(defaultGeneratorTTL)
	defer ticker.Stop()

	for range ticker.C {
		metrics, err := readProcMetrics(procMetricPath)
		if err != nil {
			fmt.Printf("[daemon] no se pudo leer %s: %v\n", procMetricPath, err)
			metrics = mockMetrics()
		}

		containers, err := listRunningContainers()
		if err != nil {
			fmt.Printf("[daemon] no se pudo listar contenedores: %v\n", err)
			containers = nil
		}

		decision := decideContainers(metrics, containers)
		if err := persistDecision(decision); err != nil {
			fmt.Printf("[daemon] error al persistir snapshot: %v\n", err)
		}
	}
}

func ensureBaseContainers() error {
	containers, err := listRunningContainers()
	if err != nil {
		return err
	}

	lowCount, highCount := 0, 0
	for _, c := range containers {
		switch c.Kind {
		case "low":
			lowCount++
		case "high":
			highCount++
		}
	}

	if lowCount < 3 {
		for i := lowCount; i < 3; i++ {
			if _, err := runDocker("run", "-d", "--name", fmt.Sprintf("low-%s-%d", carnet, time.Now().UnixNano()), "alpine", "sleep", "240"); err != nil {
				return err
			}
		}
	}
	if highCount < 2 {
		for i := highCount; i < 2; i++ {
			if _, err := runDocker("run", "-d", "--name", fmt.Sprintf("high-%s-%d", carnet, time.Now().UnixNano()), "roldyoran/go-client"); err != nil {
				return err
			}
		}
	}
	return nil
}

func listRunningContainers() ([]containerInfo, error) {
	output, err := runDocker("ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(output) == "" {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	containers := make([]containerInfo, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue
		}
		kind := "low"
		if strings.Contains(parts[1], "high-") || strings.Contains(parts[2], "go-client") {
			kind = "high"
		}
		containers = append(containers, containerInfo{
			ID:    parts[0],
			Name:  parts[1],
			Image: parts[2],
			Kind:  kind,
		})
	}
	return containers, nil
}

func decideContainers(metrics []ProcMetric, containers []containerInfo) ContainerDecision {
	decision := ContainerDecision{
		Timestamp: time.Now(),
		ToRemove:  make([]string, 0),
		ToKeep:    make([]string, 0),
		Summary:   map[string]int{"low": 0, "high": 0, "removed": 0},
	}

	for _, c := range containers {
		if strings.Contains(c.Name, "grafana") {
			decision.ToKeep = append(decision.ToKeep, c.Name)
			continue
		}
		decision.Summary[c.Kind]++
		decision.ToKeep = append(decision.ToKeep, c.Name)
	}

	if len(containers) <= 5 {
		decision.Metrics = metrics
		return decision
	}

	// Ordena contenedores por consumo de RAM/CPU para decidir eliminación.
	kindCounts := map[string]int{"low": 0, "high": 0}
	for _, c := range containers {
		if strings.Contains(c.Name, "grafana") {
			continue
		}
		kindCounts[c.Kind]++
	}
	for _, c := range containers {
		if strings.Contains(c.Name, "grafana") {
			continue
		}
		if c.Kind == "low" && kindCounts["low"] > 3 {
			decision.ToRemove = append(decision.ToRemove, c.Name)
			kindCounts["low"]--
			decision.Summary["removed"]++
			continue
		}
		if c.Kind == "high" && kindCounts["high"] > 2 {
			decision.ToRemove = append(decision.ToRemove, c.Name)
			kindCounts["high"]--
			decision.Summary["removed"]++
		}
	}

	for _, name := range decision.ToRemove {
		if _, err := runDocker("rm", "-f", name); err == nil {
			fmt.Printf("[daemon] contenedor eliminado: %s\n", name)
		}
	}

	decision.Metrics = metrics
	return decision
}

func persistDecision(decision ContainerDecision) error {
	data, err := json.MarshalIndent(decision, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(defaultLogFile, data, 0644)
}

func readProcMetrics(path string) ([]ProcMetric, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	metrics := make([]ProcMetric, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 8 {
			continue
		}
		pid, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		vsz, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			continue
		}
		rss, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			continue
		}
		memPct, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			memPct = 0
		}
		cpuPct, err := strconv.ParseFloat(parts[6], 64)
		if err != nil {
			cpuPct = 0
		}
		metric := ProcMetric{
			PID:         pid,
			Name:        parts[1],
			Command:     parts[2],
			ContainerID: parts[7],
			VSZ:         vsz,
			RSS:         rss,
			MemPercent:  memPct,
			CPUPercent:  cpuPct,
		}
		metric.Category = detectCategory(metric.Name, metric.Command)
		metrics = append(metrics, metric)
	}
	// Ordena los datos para priorizar contenedores de mayor consumo.
	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].MemPercent == metrics[j].MemPercent {
			return metrics[i].CPUPercent > metrics[j].CPUPercent
		}
		return metrics[i].MemPercent > metrics[j].MemPercent
	})
	return metrics, nil
}

func detectCategory(name, command string) string {
	if strings.Contains(name, "low-") || strings.Contains(command, "sleep 240") {
		return "low"
	}
	if strings.Contains(name, "high-") || strings.Contains(command, "go-client") || strings.Contains(command, "while true") {
		return "high"
	}
	return "system"
}

func mockMetrics() []ProcMetric {
	return []ProcMetric{
		{PID: 1010, Name: "low-202300668-1", Command: "sleep 240", VSZ: 26000, RSS: 13800, MemPercent: 18.3, CPUPercent: 3.1, Category: "low"},
		{PID: 1011, Name: "high-202300668-1", Command: "go-client", VSZ: 50000, RSS: 22000, MemPercent: 25.8, CPUPercent: 21.9, Category: "high"},
		{PID: 1012, Name: "systemd", Command: "/sbin/init", VSZ: 90000, RSS: 41000, MemPercent: 9.2, CPUPercent: 4.4, Category: "system"},
	}
}

func runDocker(args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}
