import { useEffect, useState } from "react";


export default function CountdownCard({ countDown , setCountDown}: { countDown?: number; setCountDown: React.Dispatch<React.SetStateAction<number>>; }) 
{

	useEffect(() => {
		if (countDown <= 0)
				return;

		const timer = setInterval(() => {
			setCountDown((prev) => prev - 1);
		}, 1000);

		return () => clearInterval(timer);
	}, [countDown]);


	function formatTime(value: number)
	{
		const minutes = Math.floor(value / 60);
		const seconds = value % 60;

		return `${minutes}:${seconds.toString().padStart(2, "0")}`;
	}

	return (
		<div className="
			w-[200px] h-[75px] shrink-0
			flex
			rounded-lg
			px-6 py-4
			shadow-[inset_0px_0px_20px_0px_rgba(0,175,255,0.05)] outline outline-1 outline-sky-950
			uppercase text-hero
			justify-center items-center
 			">
				0{formatTime(countDown)}
		</div>
	)
}