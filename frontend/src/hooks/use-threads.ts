import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ThreadItem, ThreadDetail, ThreadsResponse } from "@/lib/mock-data";
import { MOCK_TRENDING, MOCK_THREADS, MOCK_FEATURED } from "@/lib/mock-data";

export function useTrendingThreads() {
  return useQuery({
    queryKey: ["threads", "trending"],
    queryFn: () => api.get<{ threads: ThreadItem[] }>("/threads/trending"),
    staleTime: 5 * 60 * 1000,
    select: (data) => {
      const threads = data?.threads?.length ? data.threads : MOCK_TRENDING;
      return { threads, isEmpty: !data?.threads?.length };
    },
  });
}

export function useFeaturedThread() {
  return useQuery({
    queryKey: ["threads", "featured"],
    queryFn: () => api.get<ThreadDetail>("/threads/featured"),
    staleTime: 5 * 60 * 1000,
    select: (data) => data?.id ? data : MOCK_FEATURED,
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
    select: (data) => {
      const threads = data?.threads?.length ? data.threads : MOCK_THREADS;
      return {
        threads,
        pagination: data?.pagination ?? {
          page: 1,
          pageSize: 5,
          total: 0,
          totalPages: 1,
        },
        isEmpty: !data?.threads?.length,
      };
    },
  });
}
