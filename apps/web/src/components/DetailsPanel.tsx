import { Drawer, Descriptions, Tag, Tabs, Typography } from "antd";
import type { HostResult } from "../lib/types";

const { Paragraph, Text } = Typography;

export default function DetailsPanel({
  host, open, onClose,
}: { host: HostResult | null; open: boolean; onClose: () => void }) {
  return (
    <Drawer title={host?.ip} width={440} open={open} onClose={onClose}>
      {!host ? null : (
        <Tabs
          defaultActiveKey="overview"
          items={[
            {
              key: "overview",
              label: "Overview",
              children: (
                <Descriptions column={1} size="small" bordered>
                  <Descriptions.Item label="Hostname">{host.hostname || "—"}</Descriptions.Item>
                  <Descriptions.Item label="Vendor">{host.vendor || "—"}</Descriptions.Item>
                  <Descriptions.Item label="MAC">
                    <Text code>{host.mac || "—"}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="Ping">
                    {host.ping_ms !== undefined ? <Tag color="blue">{host.ping_ms} ms</Tag> : "—"}
                  </Descriptions.Item>
                  <Descriptions.Item label="Open Ports">
                    {host.open_ports?.length ? host.open_ports.join(", ") : "—"}
                  </Descriptions.Item>
                  <Descriptions.Item label="Started">{host.scan_started}</Descriptions.Item>
                  <Descriptions.Item label="Duration">{host.scan_duration_ms} ms</Descriptions.Item>
                </Descriptions>
              ),
            },
            {
              key: "services",
              label: "Services",
              children: (
                <Paragraph>
                  {host.open_ports?.length
                    ? host.open_ports.map(p => <Tag key={p}>{p}</Tag>)
                    : <span style={{opacity:.6}}>No open ports detected</span>}
                </Paragraph>
              ),
            },
            {
              key: "json",
              label: "JSON",
              children: (
                <Paragraph copyable>
                  <pre style={{ whiteSpace: "pre-wrap" }}>
                    {JSON.stringify(host, null, 2)}
                  </pre>
                </Paragraph>
              ),
            },
          ]}
        />
      )}
    </Drawer>
  );
}