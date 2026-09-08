import { useState, useRef, useEffect } from "react";

export default function JoinOperationCard()
{
	const [code, setCode] = useState("");
	const inputRef = useRef<HTMLInputElement>(null);

	useEffect(() => {
    inputRef.current?.focus();
	}, []);

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
   	 	const value = e.target.value.replace(/[^a-zA-Z0-9]/g, "").slice(0, 7);
    	setCode(value);
	};

	const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
		if (e.key !== "Enter") return;
		if (code.length !== 7) {
			console.log ("code invalide", code);
			return ;
		}

		console.log("Code validé :", code);
	}

	return (
		<div className="w-full h-full flex flex-col justify-around items-center">
			<h1 className="uppercase text-heading-medium flex justify-center" >Enter Room Code</h1>
			<div className="relative w-fit">
				<input
					ref={inputRef}
					type="text"
					inputMode="numeric"
					maxLength={7}
					value={code}
					onChange={handleChange}
					onKeyDown={handleKeyDown}
					className="absolute inset-0 h-full w-full opacity-0"
				/>

				<div className="flex gap-4 text-heading-small uppercase">
					{[0, 1, 2, 3, 4, 5, 6].map((index) => (
						<span key={index}>
							{code[index] ?? "_"}
						</span>
					))}
				</div>
			</div>
		</div>
	)
}