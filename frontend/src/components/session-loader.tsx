import { useEffect } from "react";
import { useMe } from "@/hooks/use-auth";
import { useAuthStore, User } from "@/stores/auth-store";

export function SessionLoader({ children }: { children: React.ReactNode }) {
  const { data } = useMe();

  useEffect(() => {
    if (data) {
      useAuthStore.getState().setUser(data);
    }
  }, [data]);

  return <>{children}</>;
}
