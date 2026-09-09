import { useState, useRef, useEffect } from "react";
import { useAuth } from "../../context/AuthContext";
import { checkRoomCode } from "../../services/lobby";

export default function JoinOperationCard() {
	const [code, setCode] = useState("");
	const [isChecking, setIsChecking] = useState(false);
	const [isShaking, setIsShaking] = useState(false);
	const [errorMessage, setErrorMessage] = useState<string | null>(null);
	const inputRef = useRef<HTMLInputElement>(null);
	const { authFetch } = useAuth();

	useEffect(() => {
		inputRef.current?.focus();
	}, []);

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const value = e.target.value.replace(/[^a-zA-Z0-9]/g, "").slice(0, 7);
		setCode(value);
		setErrorMessage(null);
	};

	const triggerInvalidShake = (message: string) => {
		setErrorMessage(message);
		setIsShaking(true);
		setTimeout(() => {
		setCode("");
		setIsShaking(false);
		}, 400);
	};

	const handleKeyDown = async (e: React.KeyboardEvent<HTMLInputElement>) => {
		if (e.key !== "Enter") return;
		if (code.length !== 7 || isChecking) return;

		setIsChecking(true);
		const result = await checkRoomCode(code, authFetch);
		setIsChecking(false);

		if (result.status === "ok") {
		console.log("Code validé :", code);
		return;
		}

		if (result.status === "full") {
		triggerInvalidShake("Room is full");
		return;
		}

		triggerInvalidShake("Invalid code");
	};

  return (
	<div className="w-full h-full flex flex-col justify-around items-center">
		<h1 className="uppercase text-heading-medium flex justify-center">Enter Room Code</h1>
		<div className="relative w-fit">
			<p className={`absolute -top-10 left-0 right-0 text-center text-red-500 text-label-small transition-opacity ${
				errorMessage ? "opacity-100" : "opacity-0"
			}`}>
				{errorMessage ?? "\u00A0"}
			</p>

			<input
				ref={inputRef}
				type="text"
				inputMode="numeric"
				maxLength={7}
				value={code}
				onChange={handleChange}
				onKeyDown={handleKeyDown}
				disabled={isChecking}
				className="absolute inset-0 h-full w-full opacity-0"
			/>

			<div className={`flex gap-4 text-heading-small uppercase ${isShaking ? "animate-shake" : ""}`}>
				{[0, 1, 2, 3, 4, 5, 6].map((index) => (
					<span key={index}>{code[index] ?? "_"}</span>
				))}
			</div>
		</div>
	</div>
  );
}