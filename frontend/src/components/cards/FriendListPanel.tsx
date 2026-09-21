import { useState, useEffect } from "react";
import FriendListRow from "../rows/FriendListRow";
import FriendRequestRow from "../rows/FriendRequestRow";
import type { UserProps, UserStatusType } from "../../types/user.types";
import { Search } from "lucide-react";
import PendingFriendRow from "../rows/PendingFriendRow";
import { useAuth } from "../../context/AuthContext";
import type { FriendResponse } from "../../types/friend.type";

interface Friend extends UserProps {
  status: UserStatusType;
  rank: number;
}

interface FriendListPanelProps {
  requests: UserProps[];
  pending: UserProps[];
}

export default function FriendListPanel({
  requests,
  pending,
}: FriendListPanelProps) {
    
    const { authFetch } = useAuth();

    const [activeTab, setActiveTab] = useState("friends");
    const [search, setSearch] = useState("");
    const [friends, setFriends] = useState<Friend[]>([]);

    useEffect(() => {
        const getFriends = async () => {

            const result = await authFetch("api/friend");

            if (!result) {
                return;
            }

            const formattedFriends: Friend[] = result.data.friends.map(
                (friend: FriendResponse) => ({
                    username: friend.name,
                    avatar: friend.avatar || "/avatars/avatar1.png"
                })
            );

            setFriends(formattedFriends);
        };

        getFriends();
    }, [authFetch]);

    return (
        <div className="bg-surface-base border border-border-default rounded-lg px-4 py-4">
            <h2 className="text-heading-medium text-center uppercase">
                Friends List
            </h2>

        <div className="grid grid-cols-3 mb-4">
            <button
            onClick={() => setActiveTab("friends")}
            className={`
                py-3
                text-heading-small
                uppercase
                border-b
                cursor-pointer
                ${
                activeTab === "friends"
                    ? "text-primary border-primary"
                    : "text-text-secondary border-border-default"
                }
            `}
        >
            Friends
        </button>

        <button
            onClick={() => setActiveTab("requests")}
            className={`
                py-3
                text-heading-small
                uppercase
                border-b
                cursor-pointer
                ${
                activeTab === "requests"
                    ? "text-primary border-primary"
                    : "text-text-secondary border-border-default"
                }
            `}
            >
                Requests
            </button>

            <button
                onClick={() => setActiveTab("add")}
                className={`
                    py-3
                    text-heading-small
                    uppercase
                    border-b
                    cursor-pointer
                    ${
                    activeTab === "add"
                        ? "text-primary border-primary"
                        : "text-text-secondary border-border-default"
                    }
                `}
            >
                Add Friends
            </button>
        </div>
        <div className="flex flex-col min-h-0">

        {activeTab === "add" && (
            <div className="
                flex
                items-center
                gap-3
                bg-surface-elevated
                px-4
                py-3
                border-b
                border-border-default
                rounded-lg
                mb-2
                "
            >
                <input
                    type="text"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="Search"
                    className="
                        flex-1
                        bg-transparent
                        text-text-primary
                        text-heading-small
                        uppercase
                        outline-none
                        placeholder:text-text-secondary
                    "
                />

                <Search
                    size={28}
                    className="text-primary shrink-0 cursor-pointer"
                />
            </div>
            )}

            <div className={`
                ${activeTab === "add" ? "h-48" : "h-64"}
                overflow-y-auto
                custom-scrollbar
                `}
            >
                {activeTab === "friends" &&
                    friends.map((friend) => (
                    <FriendListRow
                        key={friend.username}
                        {...friend}
                    />
                ))}

                {activeTab === "requests" &&
                    requests.map((request) => (
                    <FriendRequestRow
                        key={request.username}
                        {...request}
                    />
                ))}

                {activeTab === "add" &&
                    pending.map((pending) => (
                    <PendingFriendRow
                        key={pending.username}
                        {...pending}
                    />
                ))}
                
            </div>
        </div>
    </div>
    );
}