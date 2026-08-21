import PrimaryButton from "../buttons/PrimaryButton";
import Input from "../primitives/Input";
import { useState } from "react";
import { apiFetch } from "../../api/apiFetch";

interface AuthenticationLoginProps {
    onBackToOptions: () => void;
}

export default function AuthenticationLogin ({
    onBackToOptions,
}: AuthenticationLoginProps) {
    const [loginId, setLogin] = useState("");
    const [password, setPassword] = useState("");

    const handleLogin = () => {
        console.log(loginId);
        console.log(password);
        //need to add functions later
        // validate fields / send data to back / handle success or error
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
            <Input label="login id" type="login" value={loginId} onChange={setLogin}/>
            <Input label="password" type="password" value={password} onChange={setPassword}/>
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