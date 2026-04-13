import { useState } from "react"
import { Radar, Plus, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { api, type ScannedHost, type Device } from "@/lib/api"
import { useCreateDevice } from "@/hooks/use-devices"

interface ScanPanelProps {
  existingDevices: Device[]
}

export function ScanPanel({ existingDevices }: ScanPanelProps) {
  const [hosts, setHosts] = useState<ScannedHost[]>([])
  const [scanning, setScanning] = useState(false)
  const [error, setError] = useState("")
  const [addedMacs, setAddedMacs] = useState<Set<string>>(new Set())
  const createDevice = useCreateDevice()

  const existingMacs = new Set(
    existingDevices.map((d) => d.mac.toUpperCase())
  )

  const handleScan = async () => {
    setScanning(true)
    setError("")
    try {
      const result = await api.scanNetwork()
      setHosts(result)
      setAddedMacs(new Set())
    } catch (err) {
      setError(err instanceof Error ? err.message : "Scan failed")
    } finally {
      setScanning(false)
    }
  }

  const handleAdd = (host: ScannedHost) => {
    const name = host.hostname || host.ip
    createDevice.mutate(
      { name, mac: host.mac },
      {
        onSuccess: () => {
          setAddedMacs((prev) => new Set([...prev, host.mac.toUpperCase()]))
        },
      }
    )
  }

  const isExisting = (mac: string) => existingMacs.has(mac.toUpperCase())
  const isAdded = (mac: string) => addedMacs.has(mac.toUpperCase())

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="text-base">Network Scan</CardTitle>
        <Button size="sm" variant="outline" onClick={handleScan} disabled={scanning}>
          {scanning ? (
            <Loader2 className="h-4 w-4 mr-2 animate-spin" />
          ) : (
            <Radar className="h-4 w-4 mr-2" />
          )}
          {scanning ? "Scanning..." : "Scan LAN"}
        </Button>
      </CardHeader>
      <CardContent>
        {error && <p className="text-sm text-destructive mb-3">{error}</p>}

        {hosts.length === 0 && !scanning && !error && (
          <p className="text-sm text-muted-foreground">
            Click "Scan LAN" to discover devices on the local network.
          </p>
        )}

        {hosts.length > 0 && (
          <div className="space-y-2 max-h-80 overflow-y-auto">
            {hosts.map((host) => {
              const existing = isExisting(host.mac)
              const added = isAdded(host.mac)
              return (
                <div
                  key={host.mac}
                  className="flex items-center justify-between py-2 px-3 rounded-md bg-secondary/50 text-sm"
                >
                  <div className="min-w-0 flex-1">
                    <div className="font-medium truncate">
                      {host.hostname || host.ip}
                    </div>
                    <div className="text-xs text-muted-foreground font-mono">
                      {host.ip} &middot; {host.mac}
                    </div>
                  </div>
                  {existing || added ? (
                    <span className="text-xs text-muted-foreground ml-2 shrink-0">
                      {added ? "Added" : "Exists"}
                    </span>
                  ) : (
                    <Button
                      size="sm"
                      variant="ghost"
                      className="ml-2 shrink-0"
                      onClick={() => handleAdd(host)}
                      disabled={createDevice.isPending}
                    >
                      <Plus className="h-3 w-3 mr-1" />
                      Add
                    </Button>
                  )}
                </div>
              )
            })}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
