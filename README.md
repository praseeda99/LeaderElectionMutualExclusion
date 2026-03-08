# Leader Election and Mutual Exclusion (Distributed File Replication)

A Go-based distributed system demonstrating the **Bully Algorithm** for leader election, **Ricart-Agrawala Algorithm** for mutual exclusion, and **Consensus-based File Replication**.

## 🚀 Quick Start (Local Demo)

To run a **5-node cluster** on a single machine, open 5 terminals and run:

1.  **Node 1**: `go run ./distributed-file-replication -id 1 -port 8001 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003,4=localhost:8004,5=localhost:8005"`
2.  **Node 2**: `go run ./distributed-file-replication -id 2 -port 8002 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003,4=localhost:8004,5=localhost:8005"`
3.  **Node 3**: `go run ./distributed-file-replication -id 3 -port 8003 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003,4=localhost:8004,5=localhost:8005"`
4.  **Node 4**: `go run ./distributed-file-replication -id 4 -port 8004 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003,4=localhost:8004,5=localhost:8005"`
5.  **Node 5**: `go run ./distributed-file-replication -id 5 -port 8005 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003,4=localhost:8004,5=localhost:8005"`

## 💻 Multi-Laptop Demo Setup

To run across **5 physical laptops**:

1.  **Connect to same Wi-Fi** (Mobile Hotspot recommended).
2.  **Find the IP of each laptop** (`ipconfig` in CMD).
3.  **Prepare the Peers String**: It must contain the IPs of all 5 laptops.
    *   *Example*: `"1=IP1:8001,2=IP2:8002,3=IP3:8003,4=IP4:8004,5=IP5:8005"`
4.  **Run on each laptop** using its unique ID and the same Peers String:
    ```bash
    # Example for Laptop 1
    go run ./distributed-file-replication -id 1 -port 8001 -peers "YOUR_PEERS_STRING"
    ```

## 🛠 Interactive Commands

Once a node is running, type these into the terminal:

-   `status`: View Current Node ID, Leader, and Files.
-   `election`: Trigger a new leader election.
-   `cs`: Request access to the Critical Section.
-   `replicate`: Replicate `report.txt` across all active nodes.
-   `exitcs`: Release the Critical Section.

## 📁 System Architecture

-   **Leader Election**: Implemented using the **Chang–Roberts Algorithm** (Ring-based).
-   **Mutual Exclusion**: Implemented using the **Ricart–Agrawala Algorithm** (Distributed Lamport Timestamps).
-   **Consensus**: Files are only committed if a majority of nodes acknowledge the write.
-   **Storage**: Each node stores its files in the `node_storage/nodeX/` directory.

## ⚠️ Troubleshooting

-   **Firewall**: If nodes can't see each other, allow `go.exe` through Windows Firewall.
-   **Quorum**: In a 5-node setup, at least 3 nodes must be online for replication to succeed.
