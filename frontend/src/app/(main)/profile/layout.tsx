import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "My Profile | Midnight Forge",
  description: "View and manage your Midnight Forge profile, contributions, and digital vault.",
};

export default function ProfileLayout({ children }: { children: React.ReactNode }) {
  return children;
}
