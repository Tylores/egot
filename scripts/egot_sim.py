import opendssdirect as dss
import pandas as pd
import matplotlib.pyplot as plt
import os
import json
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class EGoTGridSim:
    def __init__(self, dss_file, mapping_file=None):
        self.dss_file = os.path.abspath(dss_file)
        self.mapping_file = os.path.abspath(mapping_file) if mapping_file else None
        self.der_mapping = {}
        self.results = []
        
        self._initialize_dss()
        if self.mapping_file:
            self._load_mapping()

    def _initialize_dss(self):
        logger.info(f"Loading OpenDSS model: {self.dss_file}")
        dss.Basic.ClearAll()
        dss.Text.Command(f"Compile ({self.dss_file})")
        self.nodes = dss.Circuit.AllBusNames()
        logger.info(f"Model loaded with {len(self.nodes)} nodes.")

    def _load_mapping(self):
        if os.path.exists(self.mapping_file):
            with open(self.mapping_file, 'r') as f:
                self.der_mapping = json.load(f)
            logger.info(f"Loaded mapping for {len(self.der_mapping)} DERs.")
        else:
            logger.warning(f"Mapping file {self.mapping_file} not found. Will use fallback mapping.")

    def add_der(self, device_id, bus, kw=0, kv=0.24):
        """Dynamically add a DER to the OpenDSS circuit."""
        dss.Text.Command(f"New Generator.{device_id} Bus1={bus} kW={kw} kV={kv} Phases=1")
        logger.debug(f"Added DER {device_id} at bus {bus}")

    def run_time_series(self, telemetry_csv):
        telemetry_csv = os.path.abspath(telemetry_csv)
        if not os.path.exists(telemetry_csv):
            logger.error(f"Telemetry file {telemetry_csv} not found.")
            return
        
        df = pd.read_csv(telemetry_csv)
        devices = df['Device'].unique()
        
        # Ensure all devices are in the circuit
        for dev in devices:
            if dev not in self.der_mapping:
                # Map to the end of the feeder (node 671) to highlight voltage drop/rise
                target_node = '671'
                self.der_mapping[dev] = target_node
                logger.info(f"Mapping device: {dev} -> {target_node}")
            
            self.add_der(dev, self.der_mapping[dev], kv=2.4)

        logger.info("Starting time-series simulation...")
        dss.Text.Command("Set mode=daily stepsize=1m number=1")
        
        # Sort by timestamp/step if available
        if 'Step' not in df.columns:
            df['Step'] = range(len(df)) # Fallback

        steps = df['Step'].unique()
        for step in steps:
            step_data = df[df['Step'] == step]
            
            for _, row in step_data.iterrows():
                dev = row['Device']
                # Go's Power(W) is positive for load, negative for generation.
                # OpenDSS Generator.kW is positive for generation (injection), negative for load (absorption).
                power_kw = -row['Power(W)'] / 1000.0
                dss.Text.Command(f"Generator.{dev}.kW={power_kw}")
            
            dss.Solution.Solve()
            
            # Record metrics
            self.results.append({
                'Step': step,
                'Losses_kW': dss.Circuit.Losses()[0] / 1000.0,
                'Avg_Voltage_pu': sum(dss.Circuit.AllBusMagPu()) / len(dss.Circuit.AllBusMagPu()),
                'Min_Voltage_pu': min(dss.Circuit.AllBusMagPu())
            })

        logger.info("Simulation complete.")
        return pd.DataFrame(self.results)

    def plot_results(self, res_df, output_file='egot_results.png'):
        fig, ax1 = plt.subplots(figsize=(10, 6))
        
        ax1.plot(res_df['Step'], res_df['Avg_Voltage_pu'], color='blue', label='Avg Voltage')
        ax1.plot(res_df['Step'], res_df['Min_Voltage_pu'], color='cyan', linestyle='--', label='Min Voltage')
        ax1.set_xlabel('Time Step (min)')
        ax1.set_ylabel('Voltage (pu)', color='blue')
        ax1.grid(True, alpha=0.3)
        
        ax2 = ax1.twinx()
        ax2.plot(res_df['Step'], res_df['Losses_kW'], color='red', label='Losses')
        ax2.set_ylabel('Losses (kW)', color='red')
        
        plt.title('EGoT Grid Impact Analysis')
        fig.tight_layout()
        plt.savefig(output_file)
        logger.info(f"Plot saved to {output_file}")

if __name__ == "__main__":
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    dss_path = os.path.join(base_dir, 'model/IEEE13Nodeckt.dss')

    # EIM (Scenario 2.1)
    print("⚡ Running OpenDSS simulation for Coordinated EIM...")
    sim_eim = EGoTGridSim(dss_path)
    res_eim = sim_eim.run_time_series(os.path.join(base_dir, 'data/scenario_eim.csv'))
    if res_eim is not None:
        res_eim.to_csv(os.path.join(base_dir, 'data/scenario_eim_results.csv'), index=False)
        print("Generated data/scenario_eim_results.csv")
        
    print("⚡ Running OpenDSS simulation for Baseline EIM...")
    sim_eim_base = EGoTGridSim(dss_path)
    res_eim_base = sim_eim_base.run_time_series(os.path.join(base_dir, 'data/scenario_eim_baseline.csv'))
    if res_eim_base is not None:
        res_eim_base.to_csv(os.path.join(base_dir, 'data/scenario_eim_baseline_results.csv'), index=False)
        print("Generated data/scenario_eim_baseline_results.csv")

    # Blackstart (Scenario 2.2)
    print("⚡ Running OpenDSS simulation for Coordinated Blackstart...")
    sim_bs = EGoTGridSim(dss_path)
    res_bs = sim_bs.run_time_series(os.path.join(base_dir, 'data/scenario_blackstart.csv'))
    if res_bs is not None:
        res_bs.to_csv(os.path.join(base_dir, 'data/scenario_blackstart_results.csv'), index=False)
        print("Generated data/scenario_blackstart_results.csv")
        
    print("⚡ Running OpenDSS simulation for Baseline Blackstart...")
    sim_bs_base = EGoTGridSim(dss_path)
    res_bs_base = sim_bs_base.run_time_series(os.path.join(base_dir, 'data/scenario_blackstart_baseline.csv'))
    if res_bs_base is not None:
        res_bs_base.to_csv(os.path.join(base_dir, 'data/scenario_blackstart_baseline_results.csv'), index=False)
        print("Generated data/scenario_blackstart_baseline_results.csv")
