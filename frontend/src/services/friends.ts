export type FriendsEntry = {
	user_id: number;
  	name: string;
  	avatar: string;
}

type AuthFetch = (url: string, options?: RequestInit) => Promise<{
        status: number,
        data: any,
    } | null >;

export async function getFriends(
	authFetch: AuthFetch
): Promise<FriendsEntry[]> {
	const result = await authFetch("/api/friend");

	if (result == null || result.status != 200)
		return [];

	return result.data?.friends ?? [];
}