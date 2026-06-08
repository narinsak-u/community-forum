import { useEffect } from "react";
import { useMe } from "@/hooks/use-auth";

export function SessionLoader({ children }: { children: React.ReactNode }) {
  const { data, isFetched } = useMe();

  useEffect(() => {
    if (data) {
      // useMe already calls setUser via the select function
    }
  }, [data]);

  // Don't block rendering — just restore session in background
  return <>{children}</>;
}
