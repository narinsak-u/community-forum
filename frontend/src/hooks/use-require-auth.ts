import { useCallback } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useMe } from "./use-auth";

interface RequireAuthOptions {
  redirect?: string;
  toast?: boolean;
  description?: string;
}

export function useRequireAuth() {
  const { data, isFetched } = useMe();
  const router = useRouter();

  const requireAuth = useCallback(
    (opts?: RequireAuthOptions): boolean => {
      if (!isFetched) return true;
      if (data) return true;

      if (opts?.toast) {
        toast("SIGN_IN_REQUIRED", {
          description: opts.description || "Authenticate to continue.",
          action: {
            label: "LOGIN",
            onClick: () => router.push("/login"),
          },
        });
      } else {
        const redirectPath = opts?.redirect || "/";
        router.push(`/login?redirect=${encodeURIComponent(redirectPath)}`);
      }

      return false;
    },
    [data, isFetched, router],
  );

  return { requireAuth };
}
