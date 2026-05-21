import csv
import math

# Generate 24 hours of data at 1-min intervals (1440 steps)
steps = 1440
devices = ['ESS_001', 'LOAD_001', 'PV_001']

with open('der_results.csv', 'w', newline='') as f:
    writer = csv.writer(f)
    writer.writerow(['Timestamp', 'Step', 'Device', 'Power(W)'])
    
    for step in range(steps):
        hour = step / 60.0
        for dev in devices:
            power = 0
            if dev == 'PV_001':
                if 6 <= hour <= 18:
                    power = 5000 * math.exp(-((hour-12)**2)/4)
            elif dev == 'LOAD_001':
                power = -(1000 + 2000 * math.exp(-((hour-8)**2)/2) + 4000 * math.exp(-((hour-19)**2)/4))
            elif dev == 'ESS_001':
                if 17 <= hour <= 21:
                    power = 3000
                elif 1 <= hour <= 5:
                    power = -2000
            
            writer.writerow([step * 60 * 10**9, step, dev, power])

print("Generated der_results.csv with 24 hours of data.")
