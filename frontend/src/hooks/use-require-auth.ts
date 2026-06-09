import { useCallback } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useAuthStore } from "@/stores/auth-store";

interface RequireAuthOptions {
  redirect?: string;
  toast?: boolean;
  description?: string;
}

export function useRequireAuth() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const router = useRouter();

  const requireAuth = useCallback(
    (opts?: RequireAuthOptions): boolean => {
      if (isAuthenticated) return true;

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
    [isAuthenticated, router],
  );

  return { requireAuth };
}
