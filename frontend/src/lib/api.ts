const BASE = "/api"

function getToken(): string {
  return localStorage.getItem("waken_token") ?? ""
}

export function setToken(token: string) {
  localStorage.setItem("waken_token", token)
}

export function getStoredToken(): string {
  return getToken()
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...((options.headers as Record<string, string>) ?? {}),
  }
  if (token) {
    headers["Authorization"] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE}${path}`, {
    ...options,
    headers,
  })

  if (res.status === 204) {
    return undefined as T
  }

  const data = await res.json()

  if (!res.ok) {
    throw new Error(data.error ?? `Request failed: ${res.status}`)
  }

  return data
}

export interface Device {
  id: number
  name: string
  mac: string
  broadcast_addr: string
  port: number
  created_at: string
  updated_at: string
}

export interface DeviceInput {
  name: string
  mac: string
  broadcast_addr?: string
  port?: number
}

export const api = {
  listDevices: () =>
    request<{ devices: Device[] }>("/devices").then((r) => r.devices),

  createDevice: (data: DeviceInput) =>
    request<{ device: Device }>("/devices", {
      method: "POST",
      body: JSON.stringify(data),
    }).then((r) => r.device),

  updateDevice: (id: number, data: DeviceInput) =>
    request<{ device: Device }>(`/devices/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }).then((r) => r.device),

  deleteDevice: (id: number) =>
    request<void>(`/devices/${id}`, { method: "DELETE" }),

  wakeDevice: (id: number) =>
    request<{ message: string }>(`/wake/${id}`, { method: "POST" }),

  wakeByMac: (mac: string, broadcastAddr?: string, port?: number) =>
    request<{ message: string }>("/wake", {
      method: "POST",
      body: JSON.stringify({
        mac,
        ...(broadcastAddr ? { broadcast_addr: broadcastAddr } : {}),
        ...(port ? { port } : {}),
      }),
    }),

  scanNetwork: () =>
    request<{ hosts: ScannedHost[] }>("/scan").then((r) => r.hosts),

  health: () => request<{ status: string; version: string }>("/health"),
}

export interface ScannedHost {
  ip: string
  mac: string
  hostname: string
}
