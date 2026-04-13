import { useState } from "react"
import { Pencil, Trash2, Plus, Monitor } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { WakeButton } from "@/components/wake-button"
import { DeviceForm } from "@/components/device-form"
import { ScanPanel } from "@/components/scan-panel"
import { useDevices, useDeleteDevice } from "@/hooks/use-devices"
import type { Device } from "@/lib/api"

export function DeviceList() {
  const { data: devices, isLoading, error } = useDevices()
  const deleteDevice = useDeleteDevice()
  const [formOpen, setFormOpen] = useState(false)
  const [editingDevice, setEditingDevice] = useState<Device | null>(null)

  const handleEdit = (device: Device) => {
    setEditingDevice(device)
    setFormOpen(true)
  }

  const handleAdd = () => {
    setEditingDevice(null)
    setFormOpen(true)
  }

  const handleDelete = (device: Device) => {
    if (confirm(`Delete "${device.name}"?`)) {
      deleteDevice.mutate(device.id)
    }
  }

  if (isLoading) {
    return <div className="text-center py-12 text-muted-foreground">Loading devices...</div>
  }

  if (error) {
    return (
      <div className="text-center py-12 text-destructive">
        Failed to load devices: {error.message}
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Devices</h1>
        <Button onClick={handleAdd}>
          <Plus className="h-4 w-4 mr-2" />
          Add Device
        </Button>
      </div>

      {!devices?.length ? (
        <Card>
          <CardContent className="py-12 text-center text-muted-foreground">
            <Monitor className="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p>No devices yet. Add one to get started.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {devices.map((device) => (
            <Card key={device.id}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base">{device.name}</CardTitle>
                <WakeButton deviceId={device.id} />
              </CardHeader>
              <CardContent>
                <p className="text-sm font-mono text-muted-foreground">{device.mac}</p>
                <p className="text-xs text-muted-foreground mt-1">
                  {device.broadcast_addr}:{device.port}
                </p>
                <div className="flex gap-1 mt-3">
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => handleEdit(device)}
                  >
                    <Pencil className="h-3 w-3 mr-1" />
                    Edit
                  </Button>
                  <Button
                    size="sm"
                    variant="ghost"
                    className="text-destructive hover:text-destructive"
                    onClick={() => handleDelete(device)}
                  >
                    <Trash2 className="h-3 w-3 mr-1" />
                    Delete
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <div className="mt-8">
        <ScanPanel existingDevices={devices ?? []} />
      </div>

      <DeviceForm
        open={formOpen}
        onOpenChange={setFormOpen}
        device={editingDevice}
      />
    </div>
  )
}
