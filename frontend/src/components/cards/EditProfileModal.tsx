import { useState } from "react";
import { X, Upload, Pencil } from "lucide-react";
import Avatar from "../primitives/Avatar";
import { useAuth } from "../../context/AuthContext";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { editProfileSchema } from "../../schemas/editProfileSchema";
import type { EditProfileFormData } from "../../schemas/editProfileSchema";

interface EditProfileModalProps {
    username: string;
    avatar: string;
    avatars: string[];
    onClose: () => void;
}

export default function EditProfileModal({
    username,
    avatar,
    avatars,
    onClose,
}: EditProfileModalProps) {

    const { authFetch, getUser} = useAuth();

    const [selectedAvatar, setSelectedAvatar] = useState(avatar);
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const [uploadError, setUploadError] = useState("");

    const {
        register,
        handleSubmit,
        formState: { isSubmitting, errors },
    } = useForm<EditProfileFormData>({
        resolver: zodResolver(editProfileSchema),
        defaultValues: {
            username,
        },
    });

    const handleUpload = (
        event: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = event.target.files?.[0];

        if (!file) {
            return;
        }

        const maxSize = 5 * 1024 * 1024;

        if (file.size > maxSize) {
            setUploadError("Image must be smaller than 5MB.");
            return;
        }

        setUploadError("");
        setAvatarFile(file);
        setSelectedAvatar(URL.createObjectURL(file));
    };

    const handleSave = async (data: EditProfileFormData) => {
        const usernameBody = new URLSearchParams({
            name: data.username,
        });

        const usernameResult = await authFetch("api/user/name", {
            method: "PATCH",
            body: usernameBody,
        });

        if (!usernameResult){
            return;
        }

        if (avatarFile) {
            const formData = new FormData();
            formData.append("avatar", avatarFile);

            const avatarResult = await authFetch("api/user/avatar/upload", {
                method: "POST",
                body: formData,
            });

            if (!avatarResult) {
                return ;
            }

            if (avatarResult.status !== 200) {
                setUploadError(
                    avatarResult.data?.message || "Failed to upload avatar."
                );
                return;
            }

        } else {
            const avatarBody = new URLSearchParams({
                avatar: selectedAvatar,
            });
            
            const avatarResult = await authFetch("api/user/avatar", {
                method: "PATCH",
                body: avatarBody,
            });
            if (!avatarResult) {
                return ;
            }
        }
        await getUser();
        onClose();
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

                <div className="flex items-start justify-between mb-3 md:mb-6">
                    
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

                <section className="border border-border-default rounded-lg p-5 mb-2 md:mb-5">

                    <div className="flex items-center justify-between mb-4 gap-3">
                        <h3 className="text-heading-small uppercase">
                            Avatar
                        </h3>
                        <span className="text-label-small lg:text-body-primary text-text-secondary uppercase">
                            Choose preset or upload your own
                        </span>
                    </div>

                    <div className="flex flex-col lg:flex-row gap-6">
                        <div className="flex flex-col items-center shrink-0">
                            <div className="shrink-0 lg:hidden">
                                <Avatar
                                    avatar={selectedAvatar}
                                    size={100}
                                />
                            </div>
                            <div className="shrink-0 hidden lg:block">
                                <Avatar
                                    avatar={selectedAvatar}
                                    size={140}
                                />
                            </div>
                            <span className="text-label-small text-text-secondary uppercase mt-2">
                                Current Avatar
                            </span>
                        </div>

                        <div className="flex-1 min-w-0">
                            <div className="max-h-22 md:max-h-40 overflow-y-auto custom-scrollbar pr-2">
                                <div className="grid grid-cols-3 md:grid-cols-5 lg:grid-cols-4 gap-3">
                                    {avatars.map((avatarOption) => (
                                        <button
                                            key={avatarOption}
                                            type="button"
                                            onClick={() => {
                                                setSelectedAvatar(avatarOption);
                                                setAvatarFile(null);
                                                setUploadError("");
                                            }}
                                            className={`
                                                flex
                                                items-center
                                                justify-center
                                                rounded-md
                                                md:p-1
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
                                    <p className="text-body-primary md:text-heading-small uppercase">
                                        Upload Avatar
                                    </p>

                                    <p className="text-label-small text-text-secondary">
                                        PNG, JPG
                                    </p>
                                </div>

                                <input
                                    type="file"
                                    accept="image/png,image/jpeg"
                                    onChange={handleUpload}
                                    className="hidden"
                                />
                            </label>
                            {uploadError && (
                                <p className="text-red-400 text-body-primary text-center pt-2">
                                    {uploadError}
                                </p>
                            )}
                        </div>
                    </div>

                </section>


                <section className="border border-border-default rounded-lg p-4 md:p-5">

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
                        p-3
                        md:px-4
                        md:py-3
                        "
                    >
                        <input
                            type="text"
                            {...register("username")}
                            className="
                                flex-1
                                min-w-0
                                bg-transparent
                                outline-none
                                text-text-primary
                                text-body-primary
                                md:text-heading-small
                                uppercase
                            "
                        />
                        <Pencil
                            size={20}
                            className="text-primary shrink-0"
                        />
                    </div>
                    {errors.username && (
                        <p className="text-red-400 text-sm mt-2">
                            {errors.username.message}
                        </p>
                    )}
                </section>

                <div className="flex justify-center gap-4 mt-6">
                    <button
                        type="button"
                        onClick={onClose}
                        className="
                            px-6
                            py-3
                            bg-surface-elevated
                            text-text-secondary
                            border
                            border-primary
                            rounded-lg
                            uppercase
                            text-body-primary
                            md:text-heading-medium
                            hover:drop-shadow-[0_0_8px_#35DFFF]
                            transition
                        "
                    >
                        Cancel
                    </button>
                    <button
                        type="button"
                        onClick={handleSubmit(handleSave)}
                        disabled={isSubmitting}
                        className="
                            px-6
                            py-3
                            bg-primary
                            text-surface-base
                            rounded-lg
                            uppercase
                            text-body-primary
                            md:text-heading-medium
                            hover:drop-shadow-[0_0_8px_#35DFFF]
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