import { useMemo, useState } from "react";
import { Table, Tag, Input, Space, Button, Tooltip } from "antd";
import type { ColumnsType } from "antd/es/table";
import { DownloadOutlined, ReloadOutlined } from "@ant-design/icons";
import dayjs from "dayjs";
import type { HostResult } from "../lib/types";

function pingColor(ms?: number) {
  if (ms === undefined || ms === null) return "default";
  if (ms <= 5) return "green";
  if (ms <= 30) return "blue";
  if (ms <= 80) return "gold";
  return "red";
}

function toCSV(rows: HostResult[]) {
  const header = ["ip","hostname","mac","vendor","ping_ms","open_ports","scan_started","scan_duration_ms"];
  const body = rows.map(r => [
    r.ip, r.hostname ?? "", r.mac ?? "", r.vendor ?? "", r.ping_ms ?? "",
    (r.open_ports?.join(";") ?? ""), r.scan_started, r.scan_duration_ms
  ]);
  return [header, ...body].map(a => a.map(v => `"${String(v).replace(/"/g,'""')}"`).join(",")).join("\n");
}

export default function DataTable({
  rows,
  onSelect,
  onRefresh,
}: {
  rows: HostResult[];
  onSelect: (r: HostResult) => void;
  onRefresh?: () => void;
}) {
  const [q, setQ] = useState("");

  const data = useMemo(() => {
    if (!q.trim()) return rows;
    const needle = q.toLowerCase();
    return rows.filter(r =>
      r.ip.toLowerCase().includes(needle) ||
      (r.hostname ?? "").toLowerCase().includes(needle) ||
      (r.vendor ?? "").toLowerCase().includes(needle) ||
      (r.mac ?? "").toLowerCase().includes(needle)
    );
  }, [rows, q]);

  const columns: ColumnsType<HostResult> = [
    {
      title: "IP",
      dataIndex: "ip",
      sorter: (a, b) => a.ip.localeCompare(b.ip, undefined, { numeric: true }),
      fixed: "left",
      width: 140,
      render: (ip: string, r) => (
        <a onClick={() => onSelect(r)} style={{ fontFamily: "ui-monospace, Menlo" }}>{ip}</a>
      ),
    },
    {
      title: "Hostname",
      dataIndex: "hostname",
      ellipsis: true,
      width: 220,
      sorter: (a, b) => (a.hostname ?? "").localeCompare(b.hostname ?? ""),
      render: (v?: string) => v || <span style={{opacity:.5}}>—</span>,
    },
    {
      title: "Vendor",
      dataIndex: "vendor",
      width: 220,
      sorter: (a, b) => (a.vendor ?? "").localeCompare(b.vendor ?? ""),
      render: (v?: string) => v ? <Tag color="geekblue">{v}</Tag> : <span style={{opacity:.5}}>—</span>,
    },
    {
      title: "MAC",
      dataIndex: "mac",
      width: 190,
      render: (v?: string) => <span style={{ fontFamily: "ui-monospace, Menlo" }}>{v || "—"}</span>,
    },
    {
      title: "Ping",
      dataIndex: "ping_ms",
      width: 110,
      sorter: (a, b) => (a.ping_ms ?? Infinity) - (b.ping_ms ?? Infinity),
      render: (ms?: number) => ms !== undefined
        ? <Tag color={pingColor(ms)}>{ms} ms</Tag>
        : <span style={{opacity:.5}}>—</span>,
    },
    {
      title: "Open Ports",
      dataIndex: "open_ports",
      width: 240,
      sorter: (a, b) => (a.open_ports?.length ?? 0) - (b.open_ports?.length ?? 0),
      render: (ports?: number[]) =>
        ports?.length ? (
          <Space wrap size={4}>
            {ports.slice(0, 12).map(p => <Tag key={p}>{p}</Tag>)}
            {ports.length > 12 && <Tag>+{ports.length - 12}</Tag>}
          </Space>
        ) : <span style={{opacity:.5}}>—</span>,
    },
    {
      title: "Scanned",
      dataIndex: "scan_started",
      width: 160,
      sorter: (a, b) => dayjs(a.scan_started).valueOf() - dayjs(b.scan_started).valueOf(),
      render: (ts: string) => dayjs(ts).format("MMM D, HH:mm:ss"),
    },
    {
      title: "Duration",
      dataIndex: "scan_duration_ms",
      width: 120,
      sorter: (a, b) => a.scan_duration_ms - b.scan_duration_ms,
      render: (ms: number) => `${ms} ms`,
    },
  ];

  return (
    <div>
      <div style={{ display: "flex", gap: 8, margin: "8px 0 12px", alignItems: "center" }}>
        <Input
          placeholder="Search IP / hostname / MAC / vendor"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          style={{ maxWidth: 420 }}
          allowClear
        />
        <Tooltip title="Refresh last scan (re-run form to change params)">
          <Button icon={<ReloadOutlined />} onClick={onRefresh} />
        </Tooltip>
        <Button
          icon={<DownloadOutlined />}
          onClick={() => {
            const blob = new Blob([toCSV(data)], { type: "text/csv;charset=utf-8" });
            const url = URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url; a.download = `netscan_${Date.now()}.csv`;
            a.click(); URL.revokeObjectURL(url);
          }}
        >
          Export CSV
        </Button>
      </div>
      <Table
        rowKey={(r: { ip: any; }) => r.ip}
        columns={columns}
        dataSource={data}
        size="middle"
        pagination={{ pageSize: 15, showSizeChanger: true, pageSizeOptions: [10,15,25,50,100] }}
        scroll={{ x: 1100, y: 520 }}
        sticky
      />
    </div>
  );
}