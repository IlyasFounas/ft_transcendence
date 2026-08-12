import UserProfile from "../user/UserProfile";
import SmallButton from "../buttons/SmallButton";

interface AddFriendRowProps {
    username: string;
    role: string;
    avatar: string;
}

export default function AddFriendRow({
    username,
    role,
    avatar,
}: AddFriendRowProps) {
    return (
        <div className="
            flex
            items-center
            justify-between
            bg-surface-base
            border
            border-border-default
            px-4
            h-16
            "
        >
            <UserProfile 
            username={username}
            role={role}
            avatar={avatar}
            />
            <SmallButton>Add Friend</SmallButton>
        </div>
    );
}