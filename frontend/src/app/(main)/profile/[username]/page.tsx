import { ProfileContent } from "../ProfileContent";
import { fetchJSON, apiUrl } from "@/lib/server-fetch";

export default async function UserProfilePage({
  params,
}: {
  params: Promise<{ username: string }>;
}) {
  const { username } = await params;
  const data = await fetchJSON<any>(apiUrl(`/users/${username}`));
  return <ProfileContent username={username} initialProfile={data?.user} />;
}
