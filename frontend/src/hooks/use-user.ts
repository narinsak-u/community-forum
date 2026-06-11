import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { UserProfile, ThreadsResponse } from "@/lib/mock-data";
import { queryKeys } from "./use-query-keys";

interface UserResponse {
  user: UserProfile;
}

export function useUserProfile(username: string) {
  return useQuery({
    queryKey: queryKeys.users.profile(username),
    queryFn: () => api.get<UserResponse>(`/users/${username}`).then((r) => r.user),
    staleTime: 5 * 60 * 1000,
  });
}

interface UpdateProfileData {
  avatar?: string;
  bio?: string;
  stacks?: string[];
}

export function useUpdateProfile(username: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateProfileData) =>
      api.patch<UserResponse>(`/users/${username}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.profile(username) });
    },
  });
}

export function useUserThreads(username: string, page = 1, pageSize = 5) {
  return useQuery({
    queryKey: queryKeys.users.threads(username),
    queryFn: () =>
      api.get<ThreadsResponse>(
        `/users/${username}/threads?page=${page}&pageSize=${pageSize}`
      ),
    staleTime: 60 * 1000,
  });
}
