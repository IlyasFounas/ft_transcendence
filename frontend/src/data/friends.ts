import type { UserProps, UserStatusType } from "../types/user.types";

interface Friend extends UserProps {
  status: UserStatusType;
  rank: number;
}

export const friends: Friend[] = [
  {
    username: "Mark",
    role: "Field Analyst",
    avatar: "/avatars/avatar1.png",
    status: "online",
    rank: 12,
  },
  {
    username: "Louis",
    role: "Senior Analyst",
    avatar: "/avatars/avatar2.png",
    status: "offline",
    rank: 27,
  },
  {
    username: "Tom",
    role: "Intelligence Officer",
    avatar: "/avatars/avatar3.png",
    status: "online",
    rank: 8,
  },
  {
    username: "Lorane",
    role: "Senior Analyst",
    avatar: "/avatars/avatar2.png",
    status: "offline",
    rank: 27,
  },
  {
    username: "Cipher",
    role: "Intelligence Officer",
    avatar: "/avatars/avatar3.png",
    status: "online",
    rank: 8,
  },
];

export const requests: UserProps[] = [
  {
    username: "Louis",
    role: "Junior Agent",
    avatar: "/avatars/avatar4.png",
  },
  {
    username: "Mark",
    role: "Field Analyst",
    avatar: "/avatars/avatar5.png",
  },
  {
    username: "Dean",
    role: "Junior Agent",
    avatar: "/avatars/avatar4.png",
  },
  {
    username: "Mike",
    role: "Field Analyst",
    avatar: "/avatars/avatar5.png",
  },{
    username: "Minnie",
    role: "Junior Agent",
    avatar: "/avatars/avatar4.png",
  },
];

export const pending: UserProps[] = [
  {
    username: "Mark",
    role: "Senior Analyst",
    avatar: "/avatars/avatar6.png",
  },
];