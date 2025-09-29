import { useState } from "react";
import { Button, Input, InputNumber, Select, Space, Tooltip } from "antd";
import { PlayCircleOutlined, ThunderboltOutlined } from "@ant-design/icons";
import type { ScanRequest, HostResult } from "../../lib/types";
import { runScan } from "../../lib/api";

const PORT_PRESETS = [
  { label: "Top 1000 (nmap-ish)", value: "1-1024,3306,3389,5432,5900,8000-8100,8080,8443,9000-9100" },
  { label: "All (1-65535)", value: "all" },
  { label: "SSH+HTTP(S)", value: "22,80,443" },
];

export default function ScanForm({ onDone }:{ onDone:(rows:HostResult[])=>void }) {
  const [form, setForm] = useState<ScanRequest>({
    target: "192.168.1.0/24",
    ports: "1-1024",
    concurrency: 150,
    timeout_ms: 400,
    discovery: "arp-cache", // safe default (no root needed)
  } as any);

  const [loading, setLoading] = useState(false);

  function set<K extends keyof ScanRequest>(k:K, v: ScanRequest[K]) {
    setForm(prev => ({ ...prev, [k]: v }));
  }

  async function submit() {
    setLoading(true);
    try {
      onDone(await runScan(form));
    } catch (e:any) {
      alert(e.message || String(e));
    } finally {
      setLoading(false);
    }
  }

  return (
    <Space wrap size="middle" style={{ marginBottom: 12 }}>
      <div>
        <div style={{ fontSize: 12, opacity: .7 }}>Target</div>
        <Input style={{ width: 240 }} value={form.target} onChange={e=>set("target", e.target.value)} />
      </div>
      <div>
        <div style={{ fontSize: 12, opacity: .7 }}>Ports</div>
        <Input style={{ width: 260 }} value={form.ports} onChange={e=>set("ports", e.target.value)} />
      </div>
      <div>
        <div style={{ fontSize: 12, opacity: .7 }}>Preset</div>
        <Select
          style={{ width: 220 }}
          options={PORT_PRESETS}
          placeholder="Choose preset"
          onChange={(v)=>set("ports", v)}
          allowClear
        />
      </div>
      <div>
        <div style={{ fontSize: 12, opacity: .7 }}>Discovery</div>
        <Select
          style={{ width: 160 }}
          value={(form as any).discovery ?? "arp-cache"}
          onChange={(v: any)=>set("discovery" as any, v as any)}
          options={[
            { value: "arp-cache", label: "Safe (no root)" },
            { value: "arp-raw", label: "Raw ARP (faster)" },
          ]}
        />
      </div>
      <div>
        <div style={{ fontSize: 12, opacity: .7 }}>Concurrency</div>
        <InputNumber min={10} max={2000} value={form.concurrency} onChange={(v)=>set("concurrency", Number(v))} />
      </div>
      <div>
        <div style={{ fontSize: 12, opacity: .7 }}>Timeout (ms)</div>
        <InputNumber min={100} max={3000} value={form.timeout_ms} onChange={(v)=>set("timeout_ms", Number(v))} />
      </div>
      <Tooltip title="Fast scan runs with modest timeouts & higher concurrency">
        <Button
          icon={<ThunderboltOutlined />}
          onClick={()=>{
            set("ports", "1-1024");
            set("concurrency", 250);
            set("timeout_ms", 250);
          }}
        >
          Fast Scan
        </Button>
      </Tooltip>
      <Button
        type="primary"
        icon={<PlayCircleOutlined />}
        loading={loading}
        onClick={submit}
      >
        Start Scan
      </Button>
    </Space>
  );
}