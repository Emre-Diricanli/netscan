import { useState } from "react";
import { ConfigProvider, theme } from "antd";
import type { HostResult } from "../lib/types";
import ScanForm from "../features/scan/ScanForm";
import DataTable from "../components/DataTable";
import DetailsPanel from "../components/DetailsPanel";

export default function App() {
  const [rows, setRows] = useState<HostResult[]>([]);
  const [sel, setSel] = useState<HostResult | null>(null);

  return (
    <ConfigProvider theme={{ algorithm: theme.defaultAlgorithm }}>
      <div style={{ padding: 18, maxWidth: 1200, margin: "0 auto" }}>
        <h2 style={{ margin: 0 }}>Alpha Data Tech — Netscan</h2>
        <p style={{ marginTop: 6, color: "#64748b" }}>Scan a CIDR/IP, then click a row for details.</p>
        <ScanForm onDone={(data)=>{
          const sorted = [...data].sort((a,b)=> a.ip.localeCompare(b.ip, undefined, { numeric: true }));
          setRows(sorted); setSel(null);
        }}/>
        <DataTable
          rows={rows}
          onSelect={setSel}
          onRefresh={()=>{
            // optional: re-run last scan, or no-op if you prefer
          }}
        />
        <DetailsPanel host={sel} open={!!sel} onClose={()=>setSel(null)} />
      </div>
    </ConfigProvider>
  );
}