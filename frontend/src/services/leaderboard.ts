export type LeaderboardEntry = {
	position: number;
	username: string;
	rank: string;
	avatar: string;
	total: number;
};

type ApiFetch = (url: string, options?: RequestInit) => Promise<{ status: number; data: any } | null>;

export async function getGlobalLeaderboard(
  apiFetch: ApiFetch
): Promise<LeaderboardEntry[]> {

	const result = await apiFetch("/api/leaderboard");

	if (result === null || result.status !== 200) {
		return [];
	}

	return result.data?.leaderboard ?? [];
}
