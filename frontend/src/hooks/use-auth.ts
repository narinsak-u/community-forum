import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { User } from "@/lib/mock-data";
import { queryKeys } from "./use-query-keys";

export function useSignin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { login: string; password: string }) =>
      api.post<User>("/auth/signin", data),
    onSuccess: (user) => {
      queryClient.setQueryData(queryKeys.auth.me, user);
    },
  });
}

export function useSignup() {
  return useMutation({
    mutationFn: (data: { username: string; email: string; password: string }) => {
      if (data.username.length < 3) throw new Error("Username must be at least 3 characters");
      if (data.password.length < 6) throw new Error("Password must be at least 6 characters");
      return api.post<{ message: string }>("/auth/signup", data);
    },
  });
}

export function useSignout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.post<{ message: string }>("/auth/signout"),
    onSuccess: () => {
      queryClient.setQueryData(queryKeys.auth.me, null);
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
