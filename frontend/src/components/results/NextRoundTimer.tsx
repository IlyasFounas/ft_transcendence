import { useEffect, useState } from "react";
import { CountdownRing } from "./CountdownRing";
import PrimaryButton from "../buttons/PrimaryButton";

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
		<PrimaryButton>
			<div className="flex flex-row justify-between items-center">
				<h1 className="text-heading-medium uppercase text-primary">
					{children}
				</h1>
				<CountdownRing countdown={countdown}></CountdownRing>
			</div>
		</PrimaryButton>
	)
}