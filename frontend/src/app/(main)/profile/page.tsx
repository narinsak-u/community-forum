"use client";

import { useEffect } from "react";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { ProfileContent } from "./ProfileContent";

export default function MyProfilePage() {
  const { requireAuth } = useRequireAuth();

  useEffect(() => {
    requireAuth({ redirect: "/profile" });
  }, [requireAuth]);

  return <ProfileContent />;
}
