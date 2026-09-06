import { apiFetch } from "../../api/apiFetch";
import PrimaryButton from "../buttons/PrimaryButton";
import Input from "../primitives/Input";
import { useState } from "react";

interface AuthenticationSignUpProps {
    onLogin: () => void;
}

export default function AuthenticationSignUp ({
    onLogin,
}: AuthenticationSignUpProps) {
    const [email, setEmail] = useState("");
    const [loginId, setLoginId] = useState("");
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const handleSignUp = async () => {
        if (!loginId || !email || !username || !password) {
            setError("All fields are required");
            return;
        }

        setError("");
        setSuccess("");
        
        const body = new URLSearchParams({
            login: loginId,
            password,
            username,
            email,
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
                    text-heading-small
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
        
        <div className="
            flex
            flex-col
            w-full
            items-center
            gap-2
            py-5
            sm:px-6
            "
        >
            <Input label="email" type="email" value={email} onChange={setEmail}/>
            <Input label="login id" type="text" value={loginId} onChange={setLoginId}/>
            <Input label="username" type="text" value={username} onChange={setUsername}/>
            <Input label="password" type="password" value={password} onChange={setPassword}/>
            
            <div className="pt-2 text-status-red text-label-small sm:text-heading-small">
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
                <PrimaryButton size="small" onClick={handleSignUp}>sign up</PrimaryButton>
            </div>
        </div>
    );
}