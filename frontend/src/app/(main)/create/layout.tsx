import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Create Entry | Midnight Forge",
  description:
    "Create a new thread or contribution in the Midnight Forge community.",
};

export default function CreateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
