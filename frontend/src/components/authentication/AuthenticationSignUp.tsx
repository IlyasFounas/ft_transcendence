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
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const handleSignUp = () => {
        console.log(email);
        console.log(username);
        console.log(password);
        //need to add functions later
        // validate fields / send data to back / handle success or error
    }

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
            <Input label="username" type="text" value={username} onChange={setUsername}/>
            <Input label="password" type="password" value={password} onChange={setPassword}/>
            <PrimaryButton size="small" onClick={handleSignUp}>sign up</PrimaryButton>
            <div 
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
                onClick={onLogin}
            >
                Already have an account?
            </div>
        </div>
    );
}