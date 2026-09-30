import { useEffect, useState } from "react";
import { CountdownRing } from "./CountdownRing";

interface NextRoundTimerProps {
	targetTime: string | null;
	children: React.ReactNode;
}

export default function NextRoundTimer({ targetTime, children }: NextRoundTimerProps) {
	const [countdown, setCountdown] = useState(0);

	useEffect(() => {
		if (!targetTime) return;

		const endsAt = new Date(targetTime).getTime();

		const updateCountdown = () => {
			const remaining = Math.max(0, Math.round((endsAt - Date.now()) / 1000));
			setCountdown(remaining);
		};

		updateCountdown();
		const timer = setInterval(updateCountdown, 1000);

		return () => clearInterval(timer);
	}, [targetTime]);

	if (!targetTime) return null;
	return (
		<div className="flex flex-row items-center justify-around w-full outline outline-2 rounded-xl
			outline-offset-[-2px] 
			outline-cyan-400/60 bg-linear-298 from-slate-900 via-sky-900 to-slate-900 w-80">
			<h1 className="text-heading-medium uppercase text-primary">
				{children}
			</h1>
			<CountdownRing countdown={countdown}></CountdownRing>
		</div>
	)
}