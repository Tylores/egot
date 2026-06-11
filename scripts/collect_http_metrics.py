#!/usr/bin/env python3
import os
import re
import json

def parse_nginx_logs(log_path):
    if not os.path.exists(log_path):
        print(f"⚠️ Access log not found: {log_path}")
        return None

    # Log format: '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent ...'
    # Example: 127.0.0.1 - - [30/May/2026:09:20:51 +0000] "GET /dcap HTTP/1.1" 200 842 "-" "Go-http-client/1.1" "-"
    log_pattern = re.compile(
        r'(?P<ip>\S+)\s+-\s+(?P<user>\S+)\s+\[(?P<time>[^\]]+)\]\s+"(?P<method>[A-Z]+)\s+(?P<path>\S+)\s+HTTP/[0-9.]+"\s+(?P<status>\d+)\s+(?P<bytes>\d+)'
    )

    total_requests = 0
    total_bytes = 0
    
    phases = {
        "Registration & Discovery": {"requests": 0, "bytes_sent": 0},
        "Scheduling": {"requests": 0, "bytes_sent": 0},
        "Operation": {"requests": 0, "bytes_sent": 0},
        "Verification": {"requests": 0, "bytes_sent": 0},
        "Settlement": {"requests": 0, "bytes_sent": 0}
    }

    endpoints = {}

    with open(log_path, 'r') as f:
        for line in f:
            match = log_pattern.match(line)
            if not match:
                continue

            method = match.group('method')
            path = match.group('path')
            status = int(match.group('status'))
            bytes_sent = int(match.group('bytes'))

            total_requests += 1
            total_bytes += bytes_sent

            # Determine endpoint signature
            # Normalize path variables e.g., /edev/123/rg -> /edev/{id}/rg
            norm_path = re.sub(r'/\d+(/|$)', '/{id}/', path)
            norm_path = norm_path.rstrip('/')
            endpoint_sig = f"{method} {norm_path}"

            if endpoint_sig not in endpoints:
                endpoints[endpoint_sig] = {"count": 0, "bytes_sent": 0, "statuses": {}}
            
            endpoints[endpoint_sig]["count"] += 1
            endpoints[endpoint_sig]["bytes_sent"] += bytes_sent
            endpoints[endpoint_sig]["statuses"][status] = endpoints[endpoint_sig]["statuses"].get(status, 0) + 1

            # Map endpoint to ESI phase
            phase = classify_phase(method, norm_path)
            phases[phase]["requests"] += 1
            phases[phase]["bytes_sent"] += bytes_sent

    return {
        "summary": {
            "total_requests": total_requests,
            "total_bytes_sent": total_bytes
        },
        "phases": phases,
        "endpoints": endpoints
    }

def classify_phase(method, path):
    # Order matters for subpaths
    if "/mup" in path or "/bill" in path:
        return "Settlement"
    elif "/rsps" in path or "/dstat" in path or "/lel" in path or "/ders" in path:
        return "Verification"
    elif "actderc" in path or "actedc" in path:
        return "Operation"
    elif "/frq" in path or "/frp" in path or "/dr" in path or "/derp" in path:
        # Note: default control queries represent scheduling
        return "Scheduling"
    elif "/dcap" in path or "/tm" in path or "/edev" in path or "/der" in path:
        return "Registration & Discovery"
    else:
        return "Registration & Discovery" # fallback

if __name__ == "__main__":
    script_dir = os.path.dirname(os.path.abspath(__file__))
    repo_root = os.path.dirname(script_dir)
    log_file = os.path.join(repo_root, "logs/nginx_access.log")
    output_file = os.path.join(repo_root, "data/http_metrics.json")

    print(f"📊 Analyzing gateway communication logs from {log_file}...")
    metrics = parse_nginx_logs(log_file)
    
    if metrics:
        os.makedirs(os.path.dirname(output_file), exist_ok=True)
        with open(output_file, 'w') as f:
            json.dump(metrics, f, indent=2)
        print(f"✅ Communication metrics successfully written to {output_file}")
        
        # Print a clean console table for verification
        print("\n=== ESI Lifecycle Communication Summary ===")
        print(f"{'ESI Lifecycle Phase':<30} | {'Request Count':<15} | {'Bytes Sent':<15}")
        print("-" * 68)
        for phase, data in metrics["phases"].items():
            print(f"{phase:<30} | {data['requests']:<15} | {data['bytes_sent']:<15}")
        print("-" * 68)
        print(f"{'Total':<30} | {metrics['summary']['total_requests']:<15} | {metrics['summary']['total_bytes_sent']:<15}\n")
    else:
        print("❌ Could not parse access logs. Make sure Nginx has been run.")
