import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ThreadItem, ThreadDetail, ThreadsResponse } from "@/lib/mock-data";

export function useTrendingThreads() {
  return useQuery({
    queryKey: ["threads", "trending"],
    queryFn: () => api.get<{ threads: ThreadItem[] }>("/threads/trending"),
    staleTime: 5 * 60 * 1000,
  });
}

export function useFeaturedThread() {
  return useQuery({
    queryKey: ["threads", "featured"],
    queryFn: () => api.get<ThreadDetail>("/threads/featured"),
    staleTime: 5 * 60 * 1000,
  });
}

export function useThreads(page = 1, pageSize = 5, sort = "latest") {
  return useQuery({
    queryKey: ["threads", { page, pageSize, sort }],
    queryFn: () =>
      api.get<ThreadsResponse>(
        `/threads?page=${page}&pageSize=${pageSize}&sort=${sort}`
      ),
    staleTime: 60 * 1000,
  });
}
