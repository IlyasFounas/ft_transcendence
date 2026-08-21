export interface UserProps {
    username: string;
    role: string;
    avatar: string;
}

export type UserStatusType = "online" | "in-mission" | "offline";