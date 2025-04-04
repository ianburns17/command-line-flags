# command-line-flags
https://youtu.be/11DNBLcevz4

Port Scanner
This is a multi-threaded port scanner written in Go. The program scans specified target(s) for open ports, optionally grabs banners from services running on those ports, and outputs the results in either human-readable format or JSON format.

Features
Scan a single or multiple targets (IP addresses or hostnames).

Scan a range of ports or specific ports.

Supports banner grabbing for open ports.

Uses concurrent workers to speed up the scanning process.

Outputs results in either a human-readable format or JSON format.

Configurable timeout for connections and number of workers.

Installation
Ensure you have Go installed on your system.

Clone or download the repository to your local machine.

bash
Copy
git clone https://github.com/yourusername/port-scanner.git
cd port-scanner
Build the program.

bash
Copy
go build -o port_scanner .
Usage
You can run the port scanner from the command line with various options.

Command-line Options
-target: Specify a single target IP address or hostname (e.g., -target example.com).

-targets: Comma-separated list of target IPs or hostnames to scan (e.g., -targets example.com,192.168.1.1).

-start-port: Specify the starting port number (default is 1).

-end-port: Specify the ending port number (default is 1024).

-ports: Comma-separated list of specific ports to scan (overrides -start-port and -end-port).

-workers: Number of concurrent workers to use for scanning (default is 100).

-timeout: Timeout for each connection attempt in seconds (default is 5 seconds).

-json: Output the scan results in JSON format.

-h: Display help and usage information.

Examples
Scan a single target on specific ports:
Scan example.com on ports 80 and 443:

bash
Copy
./port_scanner -target example.com -ports 80,443
Scan a range of ports on a single target:
Scan example.com on ports 1-1024:

bash
Copy
./port_scanner -target example.com -start-port 1 -end-port 1024
Scan multiple targets:
Scan example.com and 192.168.1.1 on ports 80, 443, and 8080:

bash
Copy
./port_scanner -targets example.com,192.168.1.1 -ports 80,443,8080
Scan with a specific number of workers and timeout:
Scan example.com with 50 workers and a 3-second timeout:

bash
Copy
./port_scanner -target example.com -workers 50 -timeout 3
Scan and output results in JSON format:
Scan example.com and output the results in JSON format:

bash
Copy
./port_scanner -target example.com -json
Scan a range of ports with 20 workers:
bash
Copy
./port_scanner -targets example.com,192.168.1.1 -start-port 1 -end-port 1024 -workers 20
Output
By default, the program outputs a summary of the scan for each target, including:

Number of open ports

Total ports scanned

Time taken for the scan

If the -json flag is specified, the output will be a JSON representation of the results.

Example output:

yaml
Copy
Scanning target: example.com
Scanning target example.com: port 100/1024

Scan Summary for example.com:
Open ports: 10
Total ports scanned: 1024
Time taken: 2m30s

{
  "target": "example.com",
  "results": [
    {"port": 80, "open": true, "banner": "HTTP/1.1 200 OK"},
    {"port": 443, "open": true, "banner": "HTTP/2 200"}
  ],
  "open_count": 2,
  "total_ports": 1024,
  "duration": "2m30s"
}