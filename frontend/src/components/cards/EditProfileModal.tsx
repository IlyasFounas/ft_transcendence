import { useState } from "react";
import { X, Upload, Pencil } from "lucide-react";
import Avatar from "../primitives/Avatar";
import SecondaryButton from "../buttons/SecondaryButton";

interface EditProfileModalProps {
    username: string;
    avatar: string;
    avatars: string[];
    onClose: () => void;
    onSave: (username: string, avatar: string) => void;
}

export default function EditProfileModal({
    username,
    avatar,
    avatars,
    onClose,
    onSave,
}: EditProfileModalProps) {

    const [newUsername, setNewUsername] = useState(username);
    const [selectedAvatar, setSelectedAvatar] = useState(avatar);

    const handleUpload = (
        event: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = event.target.files?.[0];

        if (!file) {
            return;
        }

        const imageUrl = URL.createObjectURL(file);
        setSelectedAvatar(imageUrl);
    };

    const handleSave = () => {
        onSave(newUsername, selectedAvatar);
    };

    return (
        <div className="
            fixed
            inset-0
            z-50
            flex
            items-center
            justify-center
            bg-black/70
            px-4
            "
        >
            <div className="
                w-full
                max-w-2xl
                max-h-[90vh]
                overflow-y-auto
                custom-scrollbar
                bg-surface-base
                border
                border-border-default
                rounded-lg
                p-6
                shadow-[0_0_20px_rgba(53,223,255,0.35)]
                "
            >

                <div className="flex items-start justify-between mb-6">
                    
                    <div>
                        <h2 className="text-heading-medium lg:text-hero uppercase">
                            Edit Profile
                        </h2>

                        <p className="text-label-small lg:text-body-primary text-text-secondary uppercase">
                            Customize your agent identity
                        </p>
                    </div>

                    <button
                        type="button"
                        onClick={onClose}
                        className="text-text-secondary hover:text-primary transition"
                    >
                        <X size={28} />
                    </button>
                </div>

                <section className="border border-border-default rounded-lg p-5 mb-5">

                    <div className="flex items-center justify-between mb-4">
                        <h3 className="text-heading-small uppercase">
                            Avatar
                        </h3>
                        <span className="text-label-small lg:text-body-primary text-text-secondary uppercase">
                            Choose preset or upload your own
                        </span>
                    </div>

                    <div className="flex flex-col lg:flex-row gap-6">
                        <div className="flex flex-col items-center shrink-0">
                            <Avatar
                                avatar={selectedAvatar}
                                size={140}
                            />
                            <span className="text-label-small text-text-secondary uppercase mt-2">
                                Current Avatar
                            </span>
                        </div>

                        <div className="flex-1 min-w-0">
                            <div className="max-h-40 overflow-y-auto custom-scrollbar pr-2">
                                <div className="grid grid-cols-3 sm:grid-cols-4 gap-3">
                                    {avatars.map((avatarOption) => (
                                        <button
                                            key={avatarOption}
                                            type="button"
                                            onClick={() =>
                                                setSelectedAvatar(avatarOption)
                                            }
                                            className={`
                                                flex
                                                items-center
                                                justify-center
                                                rounded-md
                                                p-1
                                                transition
                                                ${
                                                    selectedAvatar === avatarOption
                                                        ? "border-2 border-primary drop-shadow-[0_0_8px_#35DFFF]"
                                                        : "border border-border-default hover:border-primary"
                                                }
                                            `}
                                        >
                                            <Avatar
                                                avatar={avatarOption}
                                                size={64}
                                            />
                                        </button>
                                    ))}
                                </div>
                            </div>

                            <label className="
                                mt-4
                                flex
                                items-center
                                justify-center
                                gap-3
                                border
                                border-border-default
                                rounded-lg
                                px-4
                                py-4
                                cursor-pointer
                                hover:border-primary
                                transition
                                "
                            >
                                <Upload
                                    size={22}
                                    className="text-primary"
                                />
                                <div>
                                    <p className="text-heading-small uppercase">
                                        Upload Avatar
                                    </p>

                                    <p className="text-label-small text-text-secondary">
                                        PNG, JPG, WEBP
                                    </p>
                                </div>

                                <input
                                    type="file"
                                    accept="image/png,image/jpeg,image/webp"
                                    onChange={handleUpload}
                                    className="hidden"
                                />
                            </label>
                        </div>
                    </div>

                </section>


                <section className="border border-border-default rounded-lg p-5">

                    <div className="flex items-center justify-between mb-3">
                        <h3 className="text-heading-small uppercase">
                            Display Name
                        </h3>
                    </div>

                    <div className="
                        flex
                        items-center
                        gap-3
                        border
                        border-primary
                        rounded-md
                        px-4
                        py-3
                        "
                    >
                        <input
                            type="text"
                            value={newUsername}
                            onChange={(event) =>
                                setNewUsername(event.target.value)
                            }
                            minLength={3}
                            maxLength={30}
                            className="
                                flex-1
                                min-w-0
                                bg-transparent
                                outline-none
                                text-text-primary
                                text-heading-small
                                uppercase
                            "
                        />
                        <Pencil
                            size={20}
                            className="text-primary shrink-0"
                        />
                    </div>
                    <p className="text-label-small lg:text-body-primary text-text-secondary mt-2">
                        {newUsername.length} / 30 characters
                    </p>
                </section>

                <div className="flex justify-end gap-4 mt-6">
                   
                    <SecondaryButton onClick={onClose}>
                        Cancel
                    </SecondaryButton>
                    <button
                        type="button"
                        onClick={handleSave}
                        className="
                            px-6
                            py-3
                            bg-primary
                            text-surface-base
                            rounded-md
                            uppercase
                            text-heading-medium
                            hover:drop-shadow-[0_0_8px_#35DFFF]
                            hover:scale-[1.05]
                            transition
                        "
                    >
                        Save Changes
                    </button>
                </div>
            </div>
        </div>
    );
}