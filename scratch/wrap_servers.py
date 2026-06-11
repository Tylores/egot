import os
import re

cmd_dir = "/home/tylor/phd/egot/cmd"

for root, dirs, files in os.walk(cmd_dir):
    for file in files:
        if file == "main.go":
            filepath = os.path.join(root, file)
            with open(filepath, "r") as f:
                content = f.read()
            
            # Check if it has ListenAndServeTLS
            if "ListenAndServeTLS" in content and "http.Server" in content:
                # Find the server initialization
                # Pattern to match:
                # server := http.Server{
                #     Addr:      routes.<Name>,
                #     TLSConfig: cfg,
                # }
                pattern = r'(server\s*:=\s*http\.Server\s*\{\s*Addr:\s*routes\.[a-zA-Z0-9]+,\s*TLSConfig:\s*cfg,?\s*\})'
                match = re.search(pattern, content)
                if match:
                    original = match.group(1)
                    # replace to add Handler
                    # We need to construct the replacement
                    # Let's extract the Addr line
                    addr_match = re.search(r'Addr:\s*(routes\.[a-zA-Z0-9]+)', original)
                    if addr_match:
                        addr_field = addr_match.group(1)
                        replacement = f"""server := http.Server{{
		Addr:      {addr_field},
		TLSConfig: cfg,
		Handler:   tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}}"""
                        new_content = content.replace(original, replacement)
                        with open(filepath, "w") as f:
                            f.write(new_content)
                        print(f"Updated {filepath}")
                else:
                    print(f"Pattern not matched in {filepath}")
