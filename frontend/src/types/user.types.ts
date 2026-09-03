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