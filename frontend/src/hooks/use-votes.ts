import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { queryKeys } from "./use-query-keys";

interface VoteData {
  value: number;
}

interface VoteResponse {
  message: string;
  vote: {
    id: number;
    value: number;
  };
}

export function useVoteThread(slug: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: VoteData) =>
      api.post<VoteResponse>(`/threads/${slug}/vote`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.threads.detail(slug) });
      queryClient.invalidateQueries({ queryKey: queryKeys.threads.all });
    },
  });
}

export function useVoteComment(threadSlug: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ commentId, value }: { commentId: number; value: number }) =>
      api.post<VoteResponse>(`/comments/${commentId}/vote`, { value }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.threads.detail(threadSlug) });
    },
  });
}
