package main

import (
 "encoding/json"
 "flag"
 "fmt"
 "net"
 "os"
 "strconv"
 "strings"
 "sync"
 "sync/atomic"
 "time"
)

// PortResult holds the scan result for a single port.
type PortResult struct {
 Port int `json:"port"`
 Open bool `json:"open"`
 Banner string `json:"banner,omitempty"`
}

// TargetScanResult holds the scan summary for one target.
type TargetScanResult struct {
 Target string `json:"target"`
 Results []PortResult `json:"results"`
 OpenCount int `json:"open_count"`
 TotalPorts int `json:"total_ports"`
 Duration string `json:"duration"`
}

// Task represents a port scan job for a specific target.
type Task struct {
 Target string
 Port int
}

// worker function processes tasks from the jobs channel.
func worker(wg *sync.WaitGroup, jobs <-chan Task, results chan<- PortResult, dialer net.Dialer, progress *int32, totalJobs int) {
 defer wg.Done()
 for task := range jobs {
  result := scanPort(task.Target, task.Port, dialer)
  results <- result

  // Update progress indicator.
  newVal := atomic.AddInt32(progress, 1)
  fmt.Printf("\rScanning target %s: port %d/%d", task.Target, newVal, totalJobs)
 }
}

// scanPort dials a target port with timeout and attempts banner grabbing.
func scanPort(target string, port int, dialer net.Dialer) PortResult {
 address := net.JoinHostPort(target, strconv.Itoa(port))
 // Attempt to connect.
 conn, err := dialer.Dial("tcp", address)
 if err != nil {
  return PortResult{Port: port, Open: false}
 }
 defer conn.Close()

 // Mark as open.
 result := PortResult{Port: port, Open: true}

 // Banner grabbing: set a short read deadline.
 conn.SetReadDeadline(time.Now().Add(2 * time.Second))
 buffer := make([]byte, 1024)
 n, err := conn.Read(buffer)
 if err == nil && n > 0 {
  result.Banner = strings.TrimSpace(string(buffer[:n]))
 }
 return result
}

// parsePorts returns a slice of ports based on the -ports flag or start/end range.
func parsePorts(portsFlag string, startPort, endPort int) ([]int, error) {
 var ports []int
 if portsFlag != "" {
  // Use specific ports provided.
  parts := strings.Split(portsFlag, ",")
  for _, pStr := range parts {
   pStr = strings.TrimSpace(pStr)
   p, err := strconv.Atoi(pStr)
   if err != nil {
    return nil, fmt.Errorf("invalid port: %s", pStr)
   }
   ports = append(ports, p)
  }
 } else {
  // Use port range.
  if startPort > endPort {
   return nil, fmt.Errorf("start-port cannot be greater than end-port")
  }
  for p := startPort; p <= endPort; p++ {
   ports = append(ports, p)
  }
 }
 return ports, nil
}

func main() {
 // Command-line flags.
 targetFlag := flag.String("target", "", "Specify a single target IP address or hostname")
 targetsFlag := flag.String("targets", "", "Comma-separated list of targets (overrides -target if provided)")
 startPort := flag.Int("start-port", 1, "Start port (default: 1)")
 endPort := flag.Int("end-port", 1024, "End port (default: 1024)")
 workers := flag.Int("workers", 100, "Number of concurrent workers (default: 100)")
 timeout := flag.Int("timeout", 5, "Connection timeout in seconds (default: 5)")
 jsonOutput := flag.Bool("json", false, "Output results in JSON format")
 portsFlag := flag.String("ports", "", "Comma-separated list of specific ports to scan (overrides start-port/end-port)")
 flag.Parse()

 // Determine the list of targets.
 var targets []string
 if *targetsFlag != "" {
  for _, t := range strings.Split(*targetsFlag, ",") {
   t = strings.TrimSpace(t)
   if t != "" {
    targets = append(targets, t)
   }
  }
 } else if *targetFlag != "" {
  targets = []string{*targetFlag}
 } else {
  fmt.Println("Error: You must specify either -target or -targets")
  flag.Usage()
  os.Exit(1)
 }

 // Build the list of ports.
 ports, err := parsePorts(*portsFlag, *startPort, *endPort)
 if err != nil {
  fmt.Println("Error parsing ports:", err)
  os.Exit(1)
 }

 // Prepare the dialer with the specified timeout.
 dialer := net.Dialer{
  Timeout: time.Duration(*timeout) * time.Second,
 }

 // For each target, perform the scan.
 var allTargetResults []TargetScanResult
 for _, target := range targets {
  fmt.Printf("\nScanning target: %s\n", target)
  startTime := time.Now()

  // Create channels and progress counter.
  jobs := make(chan Task, len(ports))
  resultsChan := make(chan PortResult, len(ports))
  var progress int32 = 0

  // Launch worker pool.
  var wg sync.WaitGroup
  for i := 0; i < *workers; i++ {
   wg.Add(1)
   go worker(&wg, jobs, resultsChan, dialer, &progress, len(ports))
  }

  // Enqueue tasks.
  for _, port := range ports {
   jobs <- Task{Target: target, Port: port}
  }
  close(jobs)
  wg.Wait()
  close(resultsChan)

  // Collect results.
  var targetResults []PortResult
  openCount := 0
  for res := range resultsChan {
   targetResults = append(targetResults, res)
   if res.Open {
    openCount++
   }
  }

  duration := time.Since(startTime)
  fmt.Printf("\nScan Summary for %s:\n", target)
  fmt.Printf(" Open ports: %d\n", openCount)
  fmt.Printf(" Total ports scanned: %d\n", len(ports))
  fmt.Printf(" Time taken: %s\n", duration)

  allTargetResults = append(allTargetResults, TargetScanResult{
   Target: target,
   Results: targetResults,
   OpenCount: openCount,
   TotalPorts: len(ports),
   Duration: duration.String(),
  })
 }

 // Output JSON if requested.
 if *jsonOutput {
  jsonData, err := json.MarshalIndent(allTargetResults, "", " ")
  if err != nil {
   fmt.Println("Error marshaling JSON:", err)
   os.Exit(1)
  }
  fmt.Println("\nJSON Output:")
  fmt.Println(string(jsonData))
 }
}