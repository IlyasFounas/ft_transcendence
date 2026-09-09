// services/lobby.ts
import { mockGameRooms, mockLobbies } from "../mocks/gameroom.mock";
import type { GameRoom, Lobby } from "../types/game";

const USE_MOCK = true;

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
  if (USE_MOCK) {
    await new Promise((res) => setTimeout(res, 300));

    const activeLobbies = mockLobbies.filter(
      (l) => l.joinCode === code && l.leftAt === null
    );

    if (activeLobbies.length === 0) return { status: "not_found" };

    const gameId = activeLobbies[0].gameId;
    const gameRoom = mockGameRooms.find((g) => g.id === gameId);
    if (!gameRoom) return { status: "not_found" };

    if (activeLobbies.length >= gameRoom.playerNb) return { status: "full" };

    return { status: "ok" };
  }

  const body = new URLSearchParams({ join_code: code });
  const result = await authFetch("/ws/lobby/join", {
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