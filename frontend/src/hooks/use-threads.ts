import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ThreadItem, ThreadDetail, ThreadsResponse } from "@/lib/mock-data";
import { queryKeys } from "./use-query-keys";

export function useTrendingThreads() {
  return useQuery({
    queryKey: queryKeys.threads.trending,
    queryFn: () => api.get<{ threads: ThreadItem[] }>("/threads/trending"),
    staleTime: 5 * 60 * 1000,
  });
}

export function useFeaturedThread() {
  return useQuery({
    queryKey: queryKeys.threads.featured,
    queryFn: () => api.get<ThreadDetail>("/threads/featured"),
    staleTime: 5 * 60 * 1000,
  });
}

export function useThreads(page = 1, pageSize = 5, sort = "latest") {
  return useQuery({
    queryKey: queryKeys.threads.list({ page, pageSize, sort }),
    queryFn: () =>
      api.get<ThreadsResponse>(
        `/threads?page=${page}&pageSize=${pageSize}&sort=${sort}`
      ),
    staleTime: 60 * 1000,
  });
}

export function useMyThreads(username: string, page = 1, pageSize = 5) {
  return useQuery({
    queryKey: [...queryKeys.users.threads(username), { page, pageSize }],
    queryFn: () =>
      api.get<ThreadsResponse>(
        `/users/${username}/threads?page=${page}&pageSize=${pageSize}`
      ),
    staleTime: 60 * 1000,
    enabled: !!username,
  });
}
