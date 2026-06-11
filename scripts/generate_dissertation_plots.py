import pandas as pd
import matplotlib.pyplot as plt
import os

# Dissertation Style Settings
plt.rcParams.update({
    "font.family": "serif",
    "font.size": 11,
    "axes.labelsize": 12,
    "axes.titlesize": 13,
    "legend.fontsize": 9,
    "xtick.labelsize": 9,
    "ytick.labelsize": 9,
    "figure.dpi": 300,
    "savefig.bbox": 'tight',
    "grid.alpha": 0.3,
    "grid.linestyle": "--"
})

# Colors - Academic / Professional Palette
COLOR_COORD = "#1A365D"    # Deep Navy Blue
COLOR_BASE = "#C53030"     # Dark Red
COLOR_GRID = "#718096"     # Slate Gray
COLOR_LOSS = "#D69E2E"     # Dark Amber/Yellow
COLOR_LIMIT = "#E53E3E"    # Warning Red

def generate_hourly_scheduling_plot(eim_results_csv, output_path):
    """Generates the Hourly Energy Scheduling grid impact plot (Scenario 2.1)."""
    if not os.path.exists(eim_results_csv):
        print(f"Error: {eim_results_csv} not found.")
        return

    df = pd.read_csv(eim_results_csv)
    
    # Create the figure
    fig, ax1 = plt.subplots(figsize=(8, 5))
    
    # Plot Voltage on Axis 1 (Left)
    ax1.set_xlabel('Time of Day (Hours)')
    ax1.set_ylabel('Voltage (pu)', color=COLOR_COORD)
    
    # x-axis is 288 steps over 24 hours, so step / 12 = hour
    hours = df['Step'] / 12.0
    
    lns1 = ax1.plot(hours, df['Avg_Voltage_pu'], color=COLOR_COORD, linewidth=2, label='Average Voltage')
    lns2 = ax1.plot(hours, df['Min_Voltage_pu'], color="#3182CE", linewidth=1.5, linestyle='-.', label='Minimum Voltage')
    ax1.tick_params(axis='y', labelcolor=COLOR_COORD)
    
    # Add ANSI C84.1 limit
    lns3 = [ax1.axhline(0.95, color=COLOR_LIMIT, linestyle=':', label='ANSI C84.1 Limit (0.95 pu)')]
    
    # Set y-limit to show voltage variation clearly
    ax1.set_ylim(0.90, 1.05)
    
    # Plot Losses on Axis 2 (Right)
    ax2 = ax1.twinx()
    ax2.set_ylabel('System Losses (kW)', color=COLOR_LOSS)
    lns4 = ax2.plot(hours, df['Losses_kW'], color=COLOR_LOSS, linewidth=1.5, linestyle='--', label='Line Losses')
    ax2.tick_params(axis='y', labelcolor=COLOR_LOSS)
    
    # Combined Legend
    lns = lns1 + lns2 + lns3 + lns4
    labs = [l.get_label() for l in lns]
    ax1.legend(lns, labs, loc='lower left')
    
    plt.title('Grid Impact of Coordinated Hourly Energy Scheduling (EGoT)')
    ax1.grid(True)
    
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    plt.savefig(output_path)
    plt.close()
    print(f"✅ Saved hourly_scheduling_impact.png to {output_path}")

def generate_eim_comparison_plot(baseline_csv, coordinated_csv, output_path):
    """Generates a comparison of EIM regulation between Baseline and EGoT Coordinated (Scenario 2.1)."""
    if not os.path.exists(baseline_csv) or not os.path.exists(coordinated_csv):
        print("Error: EIM result files not found.")
        return

    df_base = pd.read_csv(baseline_csv)
    df_coord = pd.read_csv(coordinated_csv)
    
    hours = df_base['Step'] / 12.0
    
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(8, 8), sharex=True)
    
    # Top Panel: Minimum Voltage comparison
    ax1.plot(hours, df_coord['Min_Voltage_pu'], color=COLOR_COORD, linewidth=2, label='Coordinated (EGoT)')
    ax1.plot(hours, df_base['Min_Voltage_pu'], color=COLOR_BASE, linewidth=1.5, linestyle='--', label='Baseline (Uncoordinated)')
    ax1.axhline(0.95, color=COLOR_LIMIT, linestyle=':', label='ANSI C84.1 Limit (0.95 pu)')
    
    ax1.set_ylabel('Minimum Voltage (pu)')
    ax1.set_ylim(0.90, 1.05)
    ax1.set_title('Minimum Node Voltage Comparison')
    ax1.legend(loc='lower left')
    ax1.grid(True)
    
    # Bottom Panel: Line Losses comparison
    ax2.plot(hours, df_coord['Losses_kW'], color=COLOR_COORD, linewidth=2, label='Coordinated (EGoT)')
    ax2.plot(hours, df_base['Losses_kW'], color=COLOR_BASE, linewidth=1.5, linestyle='--', label='Baseline (Uncoordinated)')
    
    ax2.set_xlabel('Time of Day (Hours)')
    ax2.set_ylabel('System Losses (kW)')
    ax2.set_title('Feeder System Losses Comparison')
    ax2.legend(loc='upper right')
    ax2.grid(True)
    
    plt.suptitle('5-Minute Energy Imbalance Market (EIM) Regulation Comparison', fontsize=14)
    plt.tight_layout()
    
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    plt.savefig(output_path)
    plt.close()
    print(f"✅ Saved eim_regulation_comparison.png to {output_path}")

def generate_blackstart_plot(baseline_csv, coordinated_csv, output_path):
    """Generates the Blackstart Coordination ramp-up comparison plot (Scenario 2.2)."""
    if not os.path.exists(baseline_csv) or not os.path.exists(coordinated_csv):
        print("Error: Blackstart telemetry files not found.")
        return

    df_base = pd.read_csv(baseline_csv)
    df_coord = pd.read_csv(coordinated_csv)
    
    # Calculate total active power per step
    # Group by Step and sum Power(W)
    total_coord = df_coord.groupby('Step')['Power(W)'].sum() / 1000.0
    total_base = df_base.groupby('Step')['Power(W)'].sum() / 1000.0
    
    steps = total_coord.index
    
    fig, ax = plt.subplots(figsize=(8, 5))
    
    ax.step(steps, total_coord, color=COLOR_COORD, linewidth=2, where='post', label='Coordinated Cold Load Pickup')
    ax.step(steps, total_base, color=COLOR_BASE, linewidth=1.5, linestyle='--', where='post', label='Baseline (Uncoordinated)')
    
    # Feeder capacity limit line (12 kW)
    ax.axhline(12.0, color=COLOR_LIMIT, linestyle=':', linewidth=2, label='Feeder Capacity Limit (12 kW)')
    
    ax.set_xlabel('Time Step (Minutes)')
    ax.set_ylabel('Total Feeder Demand (kW)')
    ax.set_title('Feeder Demand Ramp-up During Blackstart Recovery')
    ax.legend(loc='upper left')
    ax.grid(True)
    ax.set_ylim(-5, 25)
    
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    plt.savefig(output_path)
    plt.close()
    print(f"✅ Saved blackstart_ramp_rate.png to {output_path}")

if __name__ == "__main__":
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    FIGURES_DIR = "/home/tylor/phd/dissertation/Figures"
    
    # Data file paths
    eim_results = os.path.join(base_dir, "data/scenario_eim_results.csv")
    eim_base_results = os.path.join(base_dir, "data/scenario_eim_baseline_results.csv")
    blackstart_telemetry = os.path.join(base_dir, "data/scenario_blackstart.csv")
    blackstart_base_telemetry = os.path.join(base_dir, "data/scenario_blackstart_baseline.csv")
    
    print("🎨 Generating dissertation figures...")
    
    # 1. Hourly Scheduling grid impact plot
    generate_hourly_scheduling_plot(
        eim_results, 
        os.path.join(FIGURES_DIR, "hourly_scheduling_impact.png")
    )
    
    # 2. EIM regulation comparison plot
    generate_eim_comparison_plot(
        eim_base_results, 
        eim_results, 
        os.path.join(FIGURES_DIR, "eim_regulation_comparison.png")
    )
    
    # 3. Blackstart recovery comparison plot
    generate_blackstart_plot(
        blackstart_base_telemetry, 
        blackstart_telemetry, 
        os.path.join(FIGURES_DIR, "blackstart_ramp_rate.png")
    )
    
    print("🎉 All figures generated successfully.")
