import { useState, useEffect } from "react";
import FriendListRow from "../rows/FriendListRow";
import FriendRequestRow from "../rows/FriendRequestRow";
import type { UserStatusType } from "../../types/user.types";
import { Search } from "lucide-react";
import PendingFriendRow from "../rows/PendingFriendRow";
import { useAuth } from "../../context/AuthContext";
import type { FriendResponse } from "../../types/friend.type";
import AddFriendRow from "../rows/AddFriendRow";

interface Friend {
    id?: number;
    username: string;
    avatar: string;
    role?: string;
    userId: number;
    status?: UserStatusType;
    rank?: number;
    friendStatus?: string;
}

export default function FriendListPanel() {
    
    const { authFetch } = useAuth();

    const [activeTab, setActiveTab] = useState("friends");
    const [search, setSearch] = useState("");
    const [friends, setFriends] = useState<Friend[]>([]);
    const [requests, setRequests] = useState<Friend[]>([]);
    const [pending, setPending] = useState<Friend[]>([]);
    const [searchResult, setSearchResult] = useState<Friend | null>(null);
    const [searchNotFound, setSearchNotFound] = useState(false);
    
    const getFriends = async () => {
        
            const result = await authFetch("api/friend");
            
            if (!result) {
                return;
            }
            console.log(result.data);
            const formattedFriends: Friend[] = result.data.friends.map(
                (friend: FriendResponse) => ({
                    userId: friend.user_id,
                    username: friend.name,
                    avatar: friend.avatar || "/avatars/avatar1.png"
                })
            );

            setFriends(formattedFriends);
        };
    
    const getFriendRequests = async () => {
        const result = await authFetch("api/friend/requests");
        
        if (!result) {
            return;
        }
        console.log(result.data);
        const formattedRequests: Friend[] = result.data.friendRequests.map(
            (request: FriendResponse) => ({
                userId: request.user_id,
                username: request.name,
                avatar: request.avatar || "/avatars/avatar1.png"
            })
        );

        setRequests(formattedRequests);
    };

    const getPending = async () => {
        const result = await authFetch("api/friend/requests/pending");
        if (!result) {
            return;
        }
        console.log(result.data);
        const formattedPending: Friend[] = result.data.friendPending.map(
            (pending: FriendResponse) => ({
                userId: pending.user_id,
                username: pending.name,
                avatar: pending.avatar || "/avatars/avatar1.png"
            })
        );

        setPending(formattedPending);
    }

    const searchUser = async () => {
        const username = search.trim();

        if (username === "") {
            setSearchResult(null);
            setSearchNotFound(false);
            return;
        }

        const result = await authFetch(`api/friend/status/${encodeURIComponent(username)}`);

        if (!result || result.status !== 200) {
            setSearchResult(null);
            setSearchNotFound(true);
            return;
        }

        const user = result.data.friendStatus;

        if (!user) {
            setSearchResult(null);
            setSearchNotFound(true);
            return;
        }

        setSearchResult({
            userId: user.user_id,
            username: user.name,
            avatar: user.avatar || "/avatars/avatar1.png",
            friendStatus: user.status,
        });

        setSearchNotFound(false);
    };

    useEffect(() => {
        getFriends();
        getFriendRequests();
        getPending();
    }, [authFetch]);

    const handleRemoveFriend = async (userId: number) => {
        const result = await authFetch(`api/friend/${userId}/remove`, {
            method: "DELETE",
        });

        if (!result) {
            return;
        }

        if (result.status !== 200) {
            return;
        }

        setFriends((currentFriends) =>
            currentFriends.filter((friend) => friend.userId !== userId)
        );
    };

    const handleAcceptFriend = async (userId: number) => {
        const result = await authFetch(`api/friend/requests/${userId}/accept`, {
            method: "PATCH",
        });
        
        if (!result) {
            return;
        }

        if (result.status !== 200) {
            return;
        }
        
        setRequests((currentRequests) => 
            currentRequests.filter((request) => request.userId !== userId)
        );

        await getFriends();

    }

    const handleRejectRequest = async (userId: number) => {
        const result = await authFetch(`api/friend/requests/${userId}/reject`, {
            method: "DELETE",
        });
        
        if (!result) {
            return;
        }

        if (result.status !== 200) {
            return;
        }
        
        setRequests((currentRequests) => 
            currentRequests.filter((request) => request.userId !== userId)
        );
    }

    const handleCancelRequest = async (userId: number) => {
        const result = await authFetch(`api/friend/requests/${userId}/cancel`, {
            method: "DELETE",
        });
        
        if (!result) {
            return;
        }

        if (result.status !== 200) {
            return;
        }

        setPending((currentPending) => 
            currentPending.filter((pending) => pending.userId !== userId)
        );
    }

    const handleAddFriend = async (friend: Friend) => {
        const body = new URLSearchParams({
            username: friend.username,
        });

        const result = await authFetch("api/friend/add", {
            method: "POST",
            body,
        });

        if (!result) {
            return;
        }

        if (result.status !== 201) {
            return;
        }

        setSearchResult(null);
        setSearch("");
        setSearchNotFound(false);

        await getPending();
    };

    return (
        <div className="bg-surface-base border border-border-default rounded-lg px-4 py-4">
            <h2 className="text-heading-small lg:text-heading-medium text-center uppercase">
                Friends List
            </h2>

        <div className="grid grid-cols-3 mb-4">
            <button
            onClick={() => setActiveTab("friends")}
            className={`
                py-3
                text-label-small
                md:text-heading-small
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
                text-label-small
                md:text-heading-small
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
                    text-label-small
                    md:text-heading-small
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
                    onChange={(e) => {
                        setSearch(e.target.value);
                        setSearchResult(null);
                        setSearchNotFound(false);
                    }}
                    onKeyDown={(e) => {
                        if (e.key === "Enter") {
                            searchUser();
                        }
                    }}
                    placeholder="Search"
                    className="
                        flex-1
                        bg-transparent
                        text-text-primary
                        text-body-primary
                        md:text-heading-small
                        uppercase
                        outline-none
                        placeholder:text-text-secondary
                    "
                />

                <Search
                    size={22}
                    onClick={searchUser}
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
                        key={friend.userId}
                        role={friend.role || "Agent"}
                        {...friend}
                        onRemove={() => handleRemoveFriend(friend.userId)}
                    />
                ))}

                {activeTab === "requests" &&
                    requests.map((request) => (
                    <FriendRequestRow
                        key={request.username}
                        role={request.role || "Agent"}
                        {...request}
                        onAccept={() => handleAcceptFriend(request.userId)}
                        onReject={() => handleRejectRequest(request.userId)}
                    />
                ))}

                {activeTab === "add" && search.trim() === "" &&
                    pending.map((pending) => (
                        <PendingFriendRow
                            key={pending.userId}
                            role={pending.role || ""}
                            {...pending}
                            onCancel={() => handleCancelRequest(pending.userId)}
                        />
                    ))
                }

                {activeTab === "add" &&
                    search.trim() !== "" &&
                    searchNotFound && (
                        <p className="text-text-secondary text-body-primary lg:text-heading-small text-center py-6">
                            User not found
                        </p>
                    )
                }

                {activeTab === "add" &&
                    search.trim() !== "" &&
                    searchResult?.friendStatus === "none" && (
                        <AddFriendRow
                            {...searchResult}
                            role={searchResult.role || ""}
                            onAdd={() => handleAddFriend(searchResult)}
                        />
                    )
                }
                
            </div>
        </div>
    </div>
    );
}