import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api, type DeviceInput } from "@/lib/api"

export function useDevices() {
  return useQuery({
    queryKey: ["devices"],
    queryFn: api.listDevices,
  })
}

export function useCreateDevice() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: DeviceInput) => api.createDevice(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["devices"] }),
  })
}

export function useUpdateDevice() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: DeviceInput }) =>
      api.updateDevice(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["devices"] }),
  })
}

export function useDeleteDevice() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.deleteDevice(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["devices"] }),
  })
}

export function useWakeDevice() {
  return useMutation({
    mutationFn: (id: number) => api.wakeDevice(id),
  })
}
