import { useState, useEffect } from "react";
import FriendListRow from "../rows/FriendListRow";
import FriendRequestRow from "../rows/FriendRequestRow";
import type { UserStatusType } from "../../types/user.types";
import { Search } from "lucide-react";
import PendingFriendRow from "../rows/PendingFriendRow";
import { useAuth } from "../../context/AuthContext";
import type { FriendResponse } from "../../types/friend.type";
import AddFriendRow from "../rows/AddFriendRow";
import FriendRow from "../rows/FriendRow";

interface Friend {
    id?: number;
    username: string;
    avatar: string;
    rank?: string;
    userId: number;
    status?: UserStatusType;
    friendStatus?: string;
}

export default function FriendListPanel() {
    
    const { authFetch, user } = useAuth();

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
            const formattedFriends: Friend[] = (result.data.friends ?? []).map(
                (friend: FriendResponse) => ({
                    userId: friend.user_id,
                    username: friend.name,
                    avatar: friend.avatar || "/avatars/avatar1.png",
                    rank: friend.rank,
                })
            );

            setFriends(formattedFriends);
        };
    
    const getFriendRequests = async () => {
        const result = await authFetch("api/friend/requests");
        
        if (!result) {
            return;
        }
        const formattedRequests: Friend[] = (result.data.friendRequests ?? []).map(
            (request: FriendResponse) => ({
                userId: request.user_id,
                username: request.name,
                avatar: request.avatar || "/avatars/avatar1.png",
                rank: request.rank,
            })
        );

        setRequests(formattedRequests);
    };

    const getPending = async () => {
        const result = await authFetch("api/friend/requests/pending");
        if (!result) {
            return;
        }
        const formattedPending: Friend[] = (result.data.friendPending ?? []).map(
            (pending: FriendResponse) => ({
                userId: pending.user_id,
                username: pending.name,
                avatar: pending.avatar || "/avatars/avatar1.png",
                rank: pending.rank,
            })
        );

        setPending(formattedPending);
    }

    const searchUser = async () => {
        const username = search.trim().toLowerCase();

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
        
        if (username.toLowerCase() === user?.Username?.toLowerCase()) {
            setSearchResult(null);
            setSearchNotFound(false);
            return;
        }
        
        const searchedUser = result.data.friendStatus;

        if (!searchedUser) {
            setSearchResult(null);
            setSearchNotFound(true);
            return;
        }
        setSearchResult({
            userId: searchedUser.user_id,
            username: searchedUser.name,
            avatar: searchedUser.avatar || "/avatars/avatar1.png",
            rank: searchedUser.rank,
            friendStatus: searchedUser.status,
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
        <div
            className="
                bg-surface-base
                border
                border-border-default
                rounded-lg
                p-4
                xl:p-10

                xl:h-full
                xl:min-h-0
                xl:flex
                xl:flex-col
            "
        >
            <h2 className="text-heading-small lg:text-heading-medium xl:text-heading-large text-center uppercase">
                Friends List
            </h2>

        <div className="grid grid-cols-3 mb-4">
            <button
            onClick={() => setActiveTab("friends")}
            className={`
                py-3
                xl:py-6
                text-label-small
                md:text-heading-small
                xl:text-heading-medium
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
                xl:py-6
                text-label-small
                md:text-heading-small
                xl:text-heading-medium
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
                    xl:py-6
                    text-label-small
                    md:text-heading-small
                    xl:text-heading-medium
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
        <div className="flex flex-col min-h-0 xl:flex-1">

        {activeTab === "add" && (
            <div className="
                flex
                items-center
                gap-3
                bg-surface-elevated
                px-4
                py-3
                xl:p-4
                border-b
                border-border-default
                rounded-lg
                mb-2
                "
            >
                <input
                    id="friend-search"
                    type="text"
                    autoComplete="off"
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
                        xl:text-heading-medium
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

            <div
                className="
                    h-64
                    xl:flex-1
                    xl:min-h-0
                    overflow-y-auto
                    custom-scrollbar
                "
            >
                {activeTab === "friends" &&
                    friends.map((friend) => (
                    <FriendListRow
                        key={friend.userId}
                        role={friend.rank || ""}
                        {...friend}
                        onRemove={() => handleRemoveFriend(friend.userId)}
                    />
                ))}

                {activeTab === "requests" &&
                    requests.map((request) => (
                    <FriendRequestRow
                        key={request.username}
                        role={request.rank || ""}
                        {...request}
                        onAccept={() => handleAcceptFriend(request.userId)}
                        onReject={() => handleRejectRequest(request.userId)}
                    />
                ))}

                {activeTab === "add" && search.trim() === "" &&
                    pending.map((pending) => (
                        <PendingFriendRow
                            key={pending.userId}
                            role={pending.rank || ""}
                            {...pending}
                            onCancel={() => handleCancelRequest(pending.userId)}
                        />
                    ))
                }

                {activeTab === "add" &&
                    search.trim() !== "" &&
                    searchNotFound && (
                        <p className="text-text-secondary text-body-primary lg:text-heading-small xl:text-heading-medium text-center py-6">
                            User not found
                        </p>
                    )
                }

                {activeTab === "add" &&
                    search.trim() !== "" &&
                    searchResult?.friendStatus === "none" && (
                        <AddFriendRow
                            {...searchResult}
                            role={searchResult.rank || ""}
                            onAdd={() => handleAddFriend(searchResult)}
                        />
                    )
                }

                {activeTab === "add" &&
                    search.trim() !== "" &&
                    searchResult?.friendStatus === "friends" && (
                        <FriendRow
                            {...searchResult}
                            role={searchResult.rank || ""}
                        />
                    )
                }
                
                {activeTab === "add" &&
                    search.trim() !== "" &&
                    searchResult?.friendStatus === "pending_sent" && 
                        pending.map((searchResult) => (
                        <PendingFriendRow
                            key={searchResult.userId}
                            role={searchResult.rank || ""}
                            {...searchResult}
                            onCancel={() => handleCancelRequest(searchResult.userId)}
                        />
                    ))
                }

                {activeTab === "add" &&
                    search.trim() !== "" &&
                    searchResult?.friendStatus === "pending_received" && 
                        requests.map((searchResult) => (
                        <FriendRequestRow
                            key={searchResult.username}
                            role={searchResult.rank || ""}
                            {...searchResult}
                            onAccept={() => handleAcceptFriend(searchResult.userId)}
                            onReject={() => handleRejectRequest(searchResult.userId)}
                        />
                    ))
                }
               
            </div>
        </div>
    </div>
    );
}
