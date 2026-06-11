#!/usr/bin/env python3
import ssl
import urllib.request
import os
import sys

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    repo_root = os.path.dirname(script_dir)
    os.chdir(repo_root)

    print("🔑 Configuring mTLS SSL Context...")
    ca_file = "ssl/ca.crt"
    client_cert = "ssl/client.crt"
    client_key = "ssl/client.key"

    if not all(os.path.exists(f) for f in [ca_file, client_cert, client_key]):
        print("❌ SSL certificate files missing in ssl/ directory!")
        sys.exit(1)

    # Setup mTLS SSL Context
    context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    context.load_verify_locations(cafile=ca_file)
    context.load_cert_chain(certfile=client_cert, keyfile=client_key)
    
    # Disable hostname checking since we hit 127.0.0.1 directly
    context.check_hostname = False
    context.verify_mode = ssl.CERT_NONE  # Avoid certificate name mismatch errors on localhost

    # List of endpoints mapping to the 5 ESI phases
    endpoints = [
        # Phase 1: Discovery & Registration
        ("GET", "/dcap", None),
        ("GET", "/tm", None),
        ("POST", "/edev", b'<EndDevice xmlns="http://ieee.org/2030.5"><sFDI>12345</sFDI></EndDevice>'),
        ("POST", "/edev/1/rg", b'<Registration xmlns="http://ieee.org/2030.5"><dateTimeRegistered>1483228800</dateTimeRegistered></Registration>'),
        
        # Phase 2: Scheduling
        ("POST", "/frq", b'<FlowReservationRequest xmlns="http://ieee.org/2030.5"><durationRequested>900</durationRequested></FlowReservationRequest>'),
        ("GET", "/frp/1", None),
        ("GET", "/dr/1/edc", None),
        
        # Phase 3: Operation
        ("GET", "/derp/1/actderc", None),
        
        # Phase 4: Verification (Responses & Statuses)
        ("POST", "/rsps/1/rsp", b'<Response xmlns="http://ieee.org/2030.5"><status>2</status></Response>'),
        ("PUT", "/edev/1/dstat", b'<DeviceStatus xmlns="http://ieee.org/2030.5"><changedTime>1483228800</changedTime></DeviceStatus>'),
        ("POST", "/edev/1/lel", b'<LogEvent xmlns="http://ieee.org/2030.5"><createdDateTime>1483228800</createdDateTime></LogEvent>'),
        
        # Phase 5: Telemetry & Settlement
        ("POST", "/mup/1", b'<MirrorUsagePoint xmlns="http://ieee.org/2030.5"><deviceLFDI>0000</deviceLFDI></MirrorUsagePoint>'),
        ("GET", "/mup/1", None)
    ]

    print("🚀 Sending mTLS requests to Nginx API Gateway...")
    
    # Send multiple requests for each to simulate periodic polling
    for i in range(5):  # Run 5 iterations to generate substantial log entries
        print(f"\n--- Traffic Generation Iteration {i+1}/5 ---")
        for method, path, data in endpoints:
            url = f"https://127.0.0.1:8443{path}"
            req = urllib.request.Request(url, data=data, method=method)
            req.add_header("Content-Type", "application/xml")
            
            try:
                with urllib.request.urlopen(req, context=context, timeout=5) as response:
                    print(f"✅ {method:<5} {path:<20} -> Status: {response.status}")
            except Exception as e:
                # If Nginx returns bad gateway because microservice is starting up, it still logs it!
                print(f"⚠️  {method:<5} {path:<20} -> Error: {e}")

    print("\n🎉 Traffic generation complete.")

if __name__ == "__main__":
    main()
