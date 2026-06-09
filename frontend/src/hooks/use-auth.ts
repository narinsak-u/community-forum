import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useAuthStore, User } from "@/stores/auth-store";
import { queryKeys } from "./use-query-keys";

export function useSignin() {
  const { setUser } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { login: string; password: string }) =>
      api.post<User>("/auth/signin", data),
    onSuccess: (user) => {
      setUser(user);
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.me });
    },
  });
}

export function useSignup() {
  return useMutation({
    mutationFn: (data: { username: string; email: string; password: string }) =>
      api.post<{ message: string }>("/auth/signup", data),
  });
}

export function useSignout() {
  const { logout } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.post<{ message: string }>("/auth/signout"),
    onSuccess: () => {
      logout();
      queryClient.invalidateQueries();
    },
  });
}

export function useMe() {
  return useQuery({
    queryKey: queryKeys.auth.me,
    queryFn: () => api.get<User>("/auth/me"),
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}
