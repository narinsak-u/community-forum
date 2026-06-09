import { ProfileContent } from "../ProfileContent";
import { fetchJSON, apiUrl } from "@/lib/server-fetch";
import type { Metadata } from "next";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ username: string }>;
}): Promise<Metadata> {
  const { username } = await params;
  return {
    title: `${username} | The Lands Between`,
    description: `View ${username}'s profile and contributions on The Lands Between.`,
  };
}

export default async function UserProfilePage({
  params,
}: {
  params: Promise<{ username: string }>;
}) {
  const { username } = await params;
  const data = await fetchJSON<any>(apiUrl(`/users/${username}`));
  return <ProfileContent username={username} initialProfile={data?.user} />;
}
