import { useRef, useEffect, useState } from "react";

interface JoinOperationCardProps {
	code: string;
	onChange: (value: string) => void;
	onSubmit: () => void;
	isChecking: boolean;
	isShaking: boolean;
	errorMessage: string | null;
}

export default function JoinOperationCard({ code, onChange, onSubmit, isChecking, isShaking, errorMessage }: JoinOperationCardProps) {
	const inputRef = useRef<HTMLInputElement>(null);
	const [isFocused, setIsFocused] = useState(false);

	useEffect(() => {
		inputRef.current?.focus();
	}, []);

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const value = e.target.value.replace(/[^a-zA-Z0-9]/g, "").slice(0, 7).toUpperCase();
		onChange(value);
	};

	const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
		if (e.key !== "Enter") return;
		onSubmit();
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
			id="room-code"
			name="roomCode"
			type="text"
			inputMode="numeric"
			autoComplete="off"
			maxLength={7}
			value={code}
			onChange={handleChange}
			onKeyDown={handleKeyDown}
			onFocus={() => setIsFocused(true)}
			onBlur={() => setIsFocused(false)}
			disabled={isChecking}
			aria-label="Room code"
			className="absolute inset-0 h-full w-full opacity-0"
			/>

			<div className={`flex gap-4 text-heading-small uppercase ${isShaking ? "animate-shake" : ""}`}>
				{[0, 1, 2, 3, 4, 5, 6].map((index) => (
					<span
						key={index}
						className={isFocused && index === code.length && !isChecking ? "animate-blink" : ""}
					>
						{code[index] ?? "_"}
					</span>
				))}
			</div>
		</div>
		</div>
	);
}