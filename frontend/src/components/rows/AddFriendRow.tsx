import UserProfile from "../user/UserProfile";
import SmallButton from "../buttons/SmallButton";
import type { UserProps } from "../../types/user.types";

interface AddFriendRowProps extends UserProps {
    onAdd?: () => void;
}
export default function AddFriendRow({
    username,
    role,
    avatar,
    onAdd,
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
            xl:h-20
            xl:px-8
            "
        >
            <UserProfile 
                username={username}
                role={role}
                avatar={avatar}
            />
            <SmallButton onClick={onAdd} >Add Friend</SmallButton>
        </div>
    );
}