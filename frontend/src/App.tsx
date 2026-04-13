import { useState } from "react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { Header } from "@/components/header"
import { DevicesPage } from "@/pages/devices"
import { SettingsPage } from "@/pages/settings"

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

type Page = "devices" | "settings"

function App() {
  const [page, setPage] = useState<Page>("devices")

  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-background">
        <Header page={page} onNavigate={setPage} />
        {page === "devices" && <DevicesPage />}
        {page === "settings" && <SettingsPage />}
      </div>
    </QueryClientProvider>
  )
}

export default App
