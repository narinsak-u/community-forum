import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ThreadItem, ThreadDetail, ThreadsResponse } from "@/lib/mock-data";
import { MOCK_TRENDING, MOCK_THREADS, MOCK_FEATURED } from "@/lib/mock-data";

export function useTrendingThreads() {
  return useQuery({
    queryKey: ["threads", "trending"],
    queryFn: async () => {
      try {
        const data = await api.get<{ threads: ThreadItem[] }>("/threads/trending");
        const threads = data?.threads?.length ? data.threads : MOCK_TRENDING;
        return { threads, isEmpty: !data?.threads?.length };
      } catch {
        return { threads: MOCK_TRENDING, isEmpty: true };
      }
    },
    staleTime: 5 * 60 * 1000,
  });
}

export function useFeaturedThread() {
  return useQuery({
    queryKey: ["threads", "featured"],
    queryFn: async () => {
      try {
        const data = await api.get<ThreadDetail>("/threads/featured");
        return data?.id ? data : MOCK_FEATURED;
      } catch {
        return MOCK_FEATURED;
      }
    },
    staleTime: 5 * 60 * 1000,
  });
}

export function useThreads(page = 1, pageSize = 5, sort = "latest") {
  return useQuery({
    queryKey: ["threads", { page, pageSize, sort }],
    queryFn: async () => {
      try {
        const data = await api.get<ThreadsResponse>(
          `/threads?page=${page}&pageSize=${pageSize}&sort=${sort}`
        );
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
      } catch {
        return {
          threads: MOCK_THREADS,
          pagination: { page: 1, pageSize: 5, total: 0, totalPages: 1 },
          isEmpty: true,
        };
      }
    },
    staleTime: 60 * 1000,
  });
}
