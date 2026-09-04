import PrimaryButton from "../buttons/PrimaryButton";
import Input from "../primitives/Input";
import { useState } from "react";
import { apiFetch } from "../../api/apiFetch";
import { useAuth } from "../../context/AuthContext";
import { useNavigate, useLocation } from "react-router-dom";

interface AuthenticationLoginProps {
    onBackToOptions: () => void;
}

export default function AuthenticationLogin ({
    onBackToOptions,
}: AuthenticationLoginProps) {
    const {login} = useAuth();
    const [loginId, setLoginId] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();
    const location = useLocation();

    const handleLogin = async () => {
        if (!loginId || !password) {
            setError("Login ID and password are required");
            return;
        }

        setError("");

        const body = new URLSearchParams({
            login: loginId,
            password,
        });

        const options: RequestInit = {
            method: "POST",
            body,
        };

        const result = await apiFetch("/api/auth/login", options);
        if (result === null) {
            setError("Unable to connect to the server. Please try again.");
            return;
        }
        if (result.status === 401 || result.status === 500) {
            setError(result.data.message);
            return;
        }
        if (result.status === 200) {
            const accessToken = result.data.Authorization;

            login(accessToken);

            const from = location.state?.from?.pathname || "/home";

            navigate(from, { replace: true });
        }
        
    };

    return (
        <div className="
            flex
            flex-col
            w-full
            items-center
            gap-4
            "
        >
            <Input label="login id" type="login" value={loginId} onChange={setLoginId}/>
            <Input label="password" type="password" value={password} onChange={setPassword}/>
            {error && <p>{error}</p>}
            <PrimaryButton size="small" onClick={handleLogin}>log in</PrimaryButton>
            <button
                type="button"
                className="
                    text-text-primary
                    underline
                    text-heading-small
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
        </div>
    );
}