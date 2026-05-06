import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { CommentItem } from "@/lib/mock-data";

interface CreateCommentData {
  content: string;
  parentId?: number;
}

interface CommentResponse {
  message: string;
  comment: CommentItem;
}

export function useCreateComment(slug: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCommentData) =>
      api.post<CommentResponse>("/threads/" + slug + "/comments", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["thread", slug] });
    },
  });
}

export function useDeleteComment(threadSlug: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (commentId: number) =>
      api.delete<{ message: string }>("/comments/" + commentId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["thread", threadSlug] });
    },
  });
}
