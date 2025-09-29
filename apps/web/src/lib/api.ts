import { API_BASE } from "./config";
import type { HostResult, ScanRequest } from "./types";

export async function runScan(req: ScanRequest): Promise<HostResult[]> {
  const res = await fetch(`${API_BASE}/scan`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}