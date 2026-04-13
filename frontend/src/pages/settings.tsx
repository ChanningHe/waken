import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { getStoredToken, setToken, api } from "@/lib/api"
import { Check, X } from "lucide-react"

export function SettingsPage() {
  const [tokenInput, setTokenInput] = useState(getStoredToken())
  const [saved, setSaved] = useState(false)
  const [health, setHealth] = useState<{ status: string; version: string } | null>(null)
  const [healthError, setHealthError] = useState("")

  useEffect(() => {
    api.health()
      .then(setHealth)
      .catch((err: Error) => setHealthError(err.message))
  }, [])

  const handleSaveToken = () => {
    setToken(tokenInput)
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  return (
    <div className="container mx-auto max-w-2xl px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">Settings</h1>

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">API Token</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label htmlFor="token">Bearer Token</Label>
              <div className="flex gap-2">
                <Input
                  id="token"
                  type="password"
                  value={tokenInput}
                  onChange={(e) => setTokenInput(e.target.value)}
                  placeholder="Enter your API token"
                />
                <Button onClick={handleSaveToken}>
                  {saved ? <Check className="h-4 w-4" /> : "Save"}
                </Button>
              </div>
              <p className="text-xs text-muted-foreground">
                Token is stored in your browser and sent with every API request.
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Server Status</CardTitle>
          </CardHeader>
          <CardContent>
            {healthError ? (
              <div className="flex items-center gap-2 text-sm text-destructive">
                <X className="h-4 w-4" />
                <span>Connection failed: {healthError}</span>
              </div>
            ) : health ? (
              <div className="space-y-1 text-sm">
                <div className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-green-600" />
                  <span>Connected</span>
                </div>
                <p className="text-muted-foreground">Version: {health.version}</p>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Checking...</p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
