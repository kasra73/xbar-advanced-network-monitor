# Xbar Advanced Network Monitor

## Overview

This project monitors bandwidth usage for primary network interfaces and displays real-time statistics with colored bars. It uses the gopsutil library to read network IO counters and calculates the bandwidth changes over a sample interval.

## How It Works

- Captures network counters at two intervals.
- Computes differences to determine bandwidth usage.
- Displays aggregated bandwidth using colored bars.
- Stores a historical maximum value to scale the display dynamically.

## Compilation

Use the provided build script to compile the project:

```bash
./build.sh
```

This builds the binary and moves it to the xbar plugins directory.

## Installation & Usage

1. Ensure Go (version 1.22.7 or later) is installed.
2. Run `bash build.sh` from the project root.
3. The plugin binary will be installed to:
   `$HOME/Library/Application Support/xbar/plugins/`
4. Configure your xbar to load the new plugin. (Refresh the xbar menu or restart it if necessary.)

## Dependencies

- Go
- [gopsutil](https://github.com/shirou/gopsutil)
