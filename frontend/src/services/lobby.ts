export type JoinRoomResult =
  | { status: "ok" }
  | { status: "not_found" }
  | { status: "full" }
  | { status: "already_in_lobby" };

type AuthFetch = (url: string, options?: RequestInit) => Promise<{ status: number; data: any } | null>;

export async function checkRoomCode(
  code: string,
  authFetch: AuthFetch
): Promise<JoinRoomResult> {

  	const body = new URLSearchParams({ join_code: code });
	const result = await authFetch("/api/lobby/join", {
	method: "POST",
	headers: { "Content-Type": "application/x-www-form-urlencoded" },
	body: body.toString(),
	});

  if (result === null) return { status: "not_found" };

  if (result.status === 200) return { status: "ok" };

  const msg = result.data?.message ?? "";
  if (msg === "The lobby is full") return { status: "full" };
  if (msg.includes("already")) return { status: "already_in_lobby" };
  return { status: "not_found" };
}