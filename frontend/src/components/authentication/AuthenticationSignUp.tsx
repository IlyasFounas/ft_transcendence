import { apiFetch } from "../../api/apiFetch";
import PrimaryButton from "../buttons/PrimaryButton";
import Input from "../primitives/Input";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { signupSchema } from "../../schemas/signupSchema";

interface AuthenticationSignUpProps {
    onLogin: () => void;
}

export default function AuthenticationSignUp ({
    onLogin,
}: AuthenticationSignUpProps) {
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const {
        register,
        handleSubmit,
        formState: { isSubmitting, errors },
    } = useForm({
        resolver: zodResolver(signupSchema),
    });

    const handleSignUp = async (data: any) => {
        
        setError("");
        setSuccess("");
        
        const body = new URLSearchParams({
            login: data.loginId,
            password: data.password,
            username: data.username,
            email: data.email,
        });

        const options: RequestInit = {
            method: "POST",
            body,
        }

        const result = await apiFetch("/api/user/create", options);
        if (result === null) {
            setError("Unable to connect to the server. Please try again.");
            return;
        }

        if (result.status === 400 || result.status === 500) {
            setError(result.data.message);
            return ;
        }

        if (result.status === 201) {
            setSuccess(result.data.message);
            return;
        }
    };

    if (success) {
        return (
            <div className="
                flex
                w-full
                flex-col
                items-center
                gap-8
                text-center
            ">
                <p className="
                    text-text-primary
                    text-body-primary
                    px-8
                    sm:text-heading-small
                ">
                    {success}
                </p>

                <PrimaryButton
                    size="small"
                    onClick={onLogin}
                    
                >
                    go to login
                </PrimaryButton>
            </div>
        );
    }

    return (
        <form
            onSubmit={handleSubmit(handleSignUp)}
            className="
                flex
                flex-col
                w-full
                items-center
                gap-2
                py-5
                sm:px-6
                "
        >
            <Input
                label="email"
                type="email"
                {...register("email")}
            />
            {errors.email && (
                <p className="text-status-red text-body-primary">
                    {errors.email.message}
                </p>
            )}

            <Input
                label="login id"
                type="text"
                {...register("loginId")}
            />

            {errors.loginId && (
                <p className="text-status-red text-body-primary">
                    {errors.loginId.message}
                </p>
            )}

            <Input
                label="username"
                type="text"
                {...register("username")}
            />

            {errors.username && (
                <p className="text-status-red text-body-primary">
                    {errors.username.message}
                </p>
            )}

            <Input
                label="password"
                type="password"
                {...register("password")}
            />

            {errors.password && (
                <p className="text-status-red text-body-primary">
                    {errors.password.message}
                </p>
            )}

            <div className="pt-2 text-status-red text-body-primary sm:text-heading-small">
                {error && 
                    <p>{error}</p>
                }
            </div>
            
            <div className="flex items-center justify-between w-full">
                <button
                    type="button"
                    className="
                        text-text-primary
                        underline
                        text-label-small
                        md:text-body-primary
                        uppercase
                        cursor-pointer
                        transition
                        duration-150
                        hover:brightness-125
                        hover:drop-shadow-[0_0_6px_rgba(53,223,255,0.5)]
                        hover:scale-[1.02]
                        active:scale-[0.98]
                    "
                    onClick={onLogin}
                >
                Already have an account?
                </button>
                <PrimaryButton size="small" type="submit">
                    {isSubmitting ? "Signing up..." : "Sign up"}
                </PrimaryButton>
            </div>
        </form>
    );
}