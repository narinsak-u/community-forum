import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "My Profile | The Lands Between",
  description:
    "View and manage your The Lands Between profile, contributions, and digital vault.",
};

export default function ProfileLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
