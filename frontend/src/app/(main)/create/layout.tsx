import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Create Entry | The Lands Between",
  description:
    "Create a new thread or contribution in The Lands Between community.",
};

export default function CreateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
