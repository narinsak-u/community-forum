import { useMe } from "@/hooks/use-auth";

export function SessionLoader({ children }: { children: React.ReactNode }) {
  useMe();

  return <>{children}</>;
}
