// mocks/gameRoom.mock.ts
import type { GameRoom } from "../types/game";

export const mockGameRooms: GameRoom[] = [
  {
    id: 7,
    title: "Room de Léo",
    mode: "classic",
    difficulty: "medium",
    timer: 60,
    roundNb: 5,
    finished: false,
    playerNb: 4,
    createdAt: "2026-09-08T10:00:00Z",
  },
];

// mocks/lobby.mock.ts
import type { Lobby } from "../types/game";

export const mockLobbies: Lobby[] = [
  { id: 1, userId: 42, gameId: 7, joinCode: "AB12CDE", isCreator: true, leftAt: null },
  { id: 2, userId: 43, gameId: 7, joinCode: "AB12CDE", isCreator: false, leftAt: null },
  { id: 3, userId: 44, gameId: 7, joinCode: "AB12CDE", isCreator: false, leftAt: null },
  { id: 4, userId: 45, gameId: 7, joinCode: "AB12CDE", isCreator: false, leftAt: null },
];