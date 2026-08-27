import { CountdownRing } from "./CircleCountDown";
import { useState, useEffect } from "react"

export default function NextRoundTimer()
{
	const [countDown, setCountDown] = useState(10);

	useEffect(() => {
		if (countDown <= 0)
				return;

		const timer = setInterval(() => {
			setCountDown((prev) => prev - 1);
		}, 1000);

		return () => clearInterval(timer);
	}, [countDown]);

	
	return (
		<div className="flex flex-row items-center justify-around w-full outline outline-2 rounded-xl
			outline-offset-[-2px] 
			outline-cyan-400/60 bg-linear-298 from-slate-900 via-sky-900 to-slate-900 w-80">
			<h1 className="text-heading-large uppercase text-primary">
				Next Round
			</h1>
			<CountdownRing countDown={countDown}></CountdownRing>
		</div>
	)
}