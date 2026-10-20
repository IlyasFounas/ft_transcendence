import type { LeaderboardEntry } from "../services/leaderboard";

export const mockLeaderboard: LeaderboardEntry[] = [
  { position: 1, username: "Nov",   rank: "Senior Agent", avatar: "avatars/avatar1.png", total: 24500 },
  { position: 2, username: "Orion",  rank: "Agent",        avatar: "avatars/avatar2.png", total: 21300 },
  { position: 3, username: "Lyra",   rank: "Agent",        avatar: "avatars/avatar4.png", total: 19800 },
  { position: 4, username: "Atlas",  rank: "Junior Agent", avatar: "avatars/avatar3.png", total: 15200 },
  { position: 5, username: "Zephyr", rank: "Junior Agent", avatar: "avatars/avatar5.png", total: 12100 },
  { position: 6, username: "Vega",   rank: "Senior Agent", avatar: "avatars/avatar6.png", total: 12100 },
  { position: 7, username: "Sirius", rank: "Agent",        avatar: "avatars/avatar7.png", total: 9875 },
  { position: 8, username: "Draco",  rank: "Agent",        avatar: "avatars/avatar8.png", total: 9850 },
  { position: 9, username: "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmm", rank: "Junior Agent", avatar: "avatars/avatar9.png", total: 8457 },
  { position: 10, username: "Pulsar", rank: "Junior Agent", avatar: "avatars/avatar10.png", total: 10 },
]