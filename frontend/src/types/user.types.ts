export interface UserProps {
    username: string;
    role: string;
    avatar: string;
}

export type UserStatusType = "online" | "in-mission" | "offline";

export interface Friend {
    username: string;
    avatar: string;
    role: string;
    status: UserStatusType;
    lastMessage: string;
    timestamp: string;
}

export interface User {
    ID: number;
    Login: string | null;
    Username: string;
    GoogleId: string | null;
    Email: string;
    Is_verified: boolean;
    Avatar: string;
    XP: number;
    Rank: string;
    CreatedAt: string;
    Lobby: unknown[] | null;
}

export interface UserLeaderboard {
    position: number;
    username: string;
    rank: string;
    total: number;
}