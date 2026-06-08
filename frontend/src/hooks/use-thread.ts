import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ThreadDetail } from "@/lib/mock-data";
import { MOCK_FEATURED } from "@/lib/mock-data";

export function useThread(slug: string) {
  return useQuery({
    queryKey: ["thread", slug],
    queryFn: () => api.get<ThreadDetail>("/threads/" + slug),
    staleTime: 60 * 1000,
    select: (data) => (data && data.id ? data : MOCK_FEATURED),
  });
}

interface CreateThreadData {
  title: string;
  content: string;
  tags?: string[];
  status?: string;
}

interface UpdateThreadData {
  title?: string;
  content?: string;
  tags?: string[];
  status?: string;
}

export function useCreateThread() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateThreadData) =>
      api.post<ThreadDetail>("/threads", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["threads"] });
    },
  });
}

export function useUpdateThread(slug: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateThreadData) =>
      api.patch<ThreadDetail>("/threads/" + slug, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["thread", slug] });
      queryClient.invalidateQueries({ queryKey: ["threads"] });
    },
  });
}

export function useDeleteThread(slug: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.delete<{ message: string }>("/threads/" + slug),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["threads"] });
    },
  });
}
