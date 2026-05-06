import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { UserProfile, MOCK_USER } from "@/lib/mock-data";

interface UserResponse {
  user: UserProfile;
}

export function useUserProfile(username: string) {
  return useQuery({
    queryKey: ["user", username],
    queryFn: () => api.get<UserResponse>("/users/" + username),
    staleTime: 5 * 60 * 1000,
    select: (data) => data?.user?.id ? data.user : MOCK_USER,
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
      api.patch<UserResponse>("/users/" + username, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", username] });
    },
  });
}
