import { cn } from "@/lib/utils"
import { MonitorSmartphone, Settings } from "lucide-react"

interface HeaderProps {
  page: "devices" | "settings"
  onNavigate: (page: "devices" | "settings") => void
}

export function Header({ page, onNavigate }: HeaderProps) {
  return (
    <header className="border-b">
      <div className="container mx-auto max-w-4xl px-4 flex items-center justify-between h-14">
        <div className="flex items-center gap-2 font-semibold">
          <MonitorSmartphone className="h-5 w-5" />
          <span>Waken</span>
        </div>
        <nav className="flex gap-1">
          <button
            onClick={() => onNavigate("devices")}
            className={cn(
              "px-3 py-1.5 text-sm rounded-md transition-colors cursor-pointer",
              page === "devices"
                ? "bg-secondary text-secondary-foreground"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            Devices
          </button>
          <button
            onClick={() => onNavigate("settings")}
            className={cn(
              "px-3 py-1.5 text-sm rounded-md transition-colors cursor-pointer",
              page === "settings"
                ? "bg-secondary text-secondary-foreground"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <Settings className="h-4 w-4" />
          </button>
        </nav>
      </div>
    </header>
  )
}
