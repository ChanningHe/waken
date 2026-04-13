import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { useCreateDevice, useUpdateDevice } from "@/hooks/use-devices"
import type { Device, DeviceInput } from "@/lib/api"

interface DeviceFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  device?: Device | null
}

export function DeviceForm({ open, onOpenChange, device }: DeviceFormProps) {
  const [form, setForm] = useState<DeviceInput>({
    name: "",
    mac: "",
    broadcast_addr: "",
    port: 0,
  })
  const [error, setError] = useState("")

  const createDevice = useCreateDevice()
  const updateDevice = useUpdateDevice()
  const isEditing = !!device

  useEffect(() => {
    if (device) {
      setForm({
        name: device.name,
        mac: device.mac,
        broadcast_addr: device.broadcast_addr,
        port: device.port,
      })
    } else {
      setForm({ name: "", mac: "", broadcast_addr: "", port: 0 })
    }
    setError("")
  }, [device, open])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError("")

    const input: DeviceInput = {
      name: form.name.trim(),
      mac: form.mac.trim(),
      ...(form.broadcast_addr ? { broadcast_addr: form.broadcast_addr.trim() } : {}),
      ...(form.port ? { port: form.port } : {}),
    }

    if (!input.name) {
      setError("Name is required")
      return
    }
    if (!input.mac) {
      setError("MAC address is required")
      return
    }

    const mutation = isEditing
      ? updateDevice.mutateAsync({ id: device.id, data: input })
      : createDevice.mutateAsync(input)

    mutation
      .then(() => onOpenChange(false))
      .catch((err: Error) => setError(err.message))
  }

  const isPending = createDevice.isPending || updateDevice.isPending

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent onClose={() => onOpenChange(false)}>
        <DialogHeader>
          <DialogTitle>{isEditing ? "Edit Device" : "Add Device"}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              placeholder="e.g. Gaming PC"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="mac">MAC Address</Label>
            <Input
              id="mac"
              value={form.mac}
              onChange={(e) => setForm({ ...form, mac: e.target.value })}
              placeholder="AA:BB:CC:DD:EE:FF"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="broadcast">Broadcast Address</Label>
              <Input
                id="broadcast"
                value={form.broadcast_addr}
                onChange={(e) => setForm({ ...form, broadcast_addr: e.target.value })}
                placeholder="255.255.255.255"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="port">Port</Label>
              <Input
                id="port"
                type="number"
                value={form.port || ""}
                onChange={(e) => setForm({ ...form, port: parseInt(e.target.value) || 0 })}
                placeholder="9"
              />
            </div>
          </div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? "Saving..." : isEditing ? "Save" : "Add"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
