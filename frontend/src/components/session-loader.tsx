import { useEffect } from "react";
import { useMe } from "@/hooks/use-auth";
import { useAuthStore, User } from "@/stores/auth-store";

export function SessionLoader({ children }: { children: React.ReactNode }) {
  const { data } = useMe();
  const setUser = useAuthStore((s) => s.setUser);

  useEffect(() => {
    if (data) setUser(data);
  }, [data, setUser]);

  return <>{children}</>;
}
