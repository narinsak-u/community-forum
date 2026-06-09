"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useMe } from "@/hooks/use-auth";

export default function MyProfilePage() {
  const router = useRouter();
  const { data: currentUser, isLoading } = useMe();

  useEffect(() => {
    if (isLoading) return;
    if (currentUser?.username) {
      router.replace("/profile/" + currentUser.username);
    } else {
      router.replace("/login?redirect=/profile");
    }
  }, [currentUser, isLoading, router]);

  return null;
}
