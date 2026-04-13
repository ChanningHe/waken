import { useState } from "react"
import { Power, Loader2, Check } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useWakeDevice } from "@/hooks/use-devices"

interface WakeButtonProps {
  deviceId: string
}

export function WakeButton({ deviceId }: WakeButtonProps) {
  const [success, setSuccess] = useState(false)
  const wake = useWakeDevice()

  const handleWake = () => {
    wake.mutate(deviceId, {
      onSuccess: () => {
        setSuccess(true)
        setTimeout(() => setSuccess(false), 2000)
      },
    })
  }

  if (success) {
    return (
      <Button size="icon" variant="outline" className="text-green-600 border-green-600" disabled>
        <Check className="h-4 w-4" />
      </Button>
    )
  }

  return (
    <Button
      size="icon"
      variant="outline"
      onClick={handleWake}
      disabled={wake.isPending}
      title="Wake device"
    >
      {wake.isPending ? (
        <Loader2 className="h-4 w-4 animate-spin" />
      ) : (
        <Power className="h-4 w-4" />
      )}
    </Button>
  )
}
