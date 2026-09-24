import PrimaryButton from "../buttons/PrimaryButton";
import Input from "../primitives/Input";
import { useState } from "react";
import { apiFetch } from "../../api/apiFetch";
import { useAuth } from "../../context/AuthContext";
import { useNavigate, useLocation } from "react-router-dom";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema } from "../../schemas/loginSchema";
import { useForm } from "react-hook-form";

interface AuthenticationLoginProps {
    onBackToOptions: () => void;
}

export default function AuthenticationLogin ({
    onBackToOptions,
}: AuthenticationLoginProps) {
    const {login} = useAuth();
    const [error, setError] = useState("");
    const navigate = useNavigate();
    const location = useLocation();

    const {
        register,
        handleSubmit,
        formState: { isSubmitting, errors },
    } = useForm ({
        resolver: zodResolver(loginSchema),
    });

    const handleLogin = async (data: any) => {

        setError("");

        const body = new URLSearchParams({
            login: data.loginId,
            password: data.password,
        })

        const options: RequestInit = {
            method: "POST",
            body,
        };

        const result = await apiFetch("/api/auth/login", options);
        if (result === null) {
            setError("Unable to connect to the server. Please try again.");
            return;
        }
        if (result.status === 401) {
            setError(result.data.message);
            return;
        }
        if (result.status === 500) {
            setError(result.data.state);
            return;
        }
        if (result.status === 200) {
            await login();

            const from = location.state?.from?.pathname || "/home";

            navigate(from, { replace: true });
        }
    };

    return (
        <form
            onSubmit={handleSubmit(handleLogin)}
            className="
                flex
                flex-col
                w-full
                items-center
                gap-3
                py-5
                sm:px-6
                "
        >
            <Input
                label="login id"
                type="text"
                {...register("loginId")}
            />
            <Input
                label="password"
                type="password"
                {...register("password")}
            />
            {errors.loginId && (
                <p className="text-status-red text-body-primary">
                    {errors.loginId.message}
                </p>
            )}
            
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

            <PrimaryButton size="small" type="submit">
                {isSubmitting ? "Loging in..." : "Log in"}
            </PrimaryButton>
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
                onClick={onBackToOptions}
            >
                Back to sign-in options
            </button>
        </form>
    );
}