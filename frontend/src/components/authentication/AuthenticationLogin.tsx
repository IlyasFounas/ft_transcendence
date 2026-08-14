import PrimaryButton from "../buttons/PrimaryButton";
import Input from "../primitives/Input";
import { useState } from "react";

interface AuthenticationLoginProps {
    onBackToOptions: () => void;
}

export default function AuthenticationLogin ({
    onBackToOptions,
}: AuthenticationLoginProps) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const handleLogin = () => {
        console.log(email);
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
            <Input label="email" type="email" value={email} onChange={setEmail}/>
            <Input label="password" type="password" value={password} onChange={setPassword}/>
            <PrimaryButton size="small" onClick={handleLogin}>log in</PrimaryButton>
            {/* <div 
                 className="
                    text-text-primary
                    underline
                    text-body-primary
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
            </div> */}
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