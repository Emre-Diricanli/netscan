export interface HostResult {
  ip: string;
  hostname?: string;
  mac?: string;
  vendor?: string;
  ping_ms?: number;
  open_ports: number[];
  scan_started: string;
  scan_duration_ms: number;
}
export interface ScanRequest {
  target: string;
  ports: string;
  concurrency?: number;
  timeout_ms?: number;
}