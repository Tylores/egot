import opendssdirect as dss
import pandas as pd
import matplotlib.pyplot as plt
import os

def run_simulation(csv_path, dss_file):
    # 1. Load the OpenDSS model
    dss.Basic.ClearAll()
    dss.Text.Command(f"Compile ({dss_file})")
    
    # 2. Load DER results
    if not os.path.exists(csv_path):
        print(f"Error: {csv_path} not found. Run egot data-export first.")
        return
    
    der_data = pd.read_csv(csv_path)
    # Pivot data: index=Timestamp, columns=Device, values=Power(W)
    # We assume timestamps are nanos and should be converted to step indices
    der_data['Step'] = range(len(der_data)) # Simplified mapping
    
    # 3. Map DERs to Nodes
    # In a real study, we would have a mapping file. 
    # For demonstration, we'll map the first N devices to random nodes.
    devices = der_data['Device'].unique()
    nodes = dss.Circuit.AllBusNames()
    
    mapping = {}
    for i, dev in enumerate(devices):
        if i < len(nodes):
            mapping[dev] = nodes[i]
            # Define Generator in OpenDSS for this DER
            dss.Text.Command(f"New Generator.{dev} Bus1={nodes[i]} kW=0 kV=0.24 Phases=1")
    
    # 4. Time-Series Simulation
    results = []
    
    # Configure simulation for time-series
    dss.Text.Command("Set mode=daily stepsize=1m number=1")
    
    for step in range(len(der_data['Step'].unique())):
        step_data = der_data[der_data['Step'] == step]
        
        # Update each generator power based on telemetry
        for _, row in step_data.iterrows():
            dev = row['Device']
            power_kw = row['Power(W)'] / 1000.0
            dss.Text.Command(f"Generator.{dev}.kW={power_kw}")
        
        # Solve power flow
        dss.Solution.Solve()
        
        # Collect results (e.g., total losses and voltage at a sensitive node)
        losses = dss.Circuit.Losses()[0] / 1000.0 # kW
        voltages = dss.Circuit.AllBusMagPu()
        avg_v = sum(voltages) / len(voltages)
        
        results.append({
            'Step': step,
            'Losses(kW)': losses,
            'AvgVoltage(pu)': avg_v
        })
    
    # 5. Analysis & Plotting
    res_df = pd.DataFrame(results)
    print("Simulation Complete.")
    print(res_df.describe())
    
    # Plotting
    fig, ax1 = plt.subplots()
    ax1.plot(res_df['Step'], res_df['AvgVoltage(pu)'], color='blue', label='Avg Voltage')
    ax1.set_xlabel('Time Step (min)')
    ax1.set_ylabel('Voltage (pu)', color='blue')
    
    ax2 = ax1.twinx()
    ax2.plot(res_df['Step'], res_df['Losses(kW)'], color='red', label='Losses')
    ax2.set_ylabel('Losses (kW)', color='red')
    
    plt.title('Grid Impact of EGoT DER Aggregation')
    plt.savefig('grid_impact.png')
    print("Results saved to grid_impact.png")

if __name__ == "__main__":
    # Example paths
    # We would use a standard model like IEEE13Nodeckt.dss
    run_simulation('der_results.csv', 'model/IEEE13Nodeckt.dss')
