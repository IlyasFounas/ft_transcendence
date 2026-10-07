import RoundsSettingsCard from "../cards/RoundsSettingsCard";

interface CreateOperationCardProps {
	onRoundsChange: (value: number) => void;
  	onTimeChange: (value: number) => void;
	errorMessage: string | null;
}

export default function CreateOperationCard({onRoundsChange, onTimeChange, errorMessage} : CreateOperationCardProps)
{
	function formatTime(value: number)
	{
		if (value < 60)
			return `${value} sec`;

		const minutes = Math.floor(value / 60);
		const seconds = value % 60;

		if (seconds === 0)
			return `${minutes} min`;

		return `${minutes} min ${seconds}`;
	}

	return (
		<div className="w-full h-full flex flex-col justify-center items-center gap-6 lg:gap-8 text-heading-small lg:text-heading-medium">
			<div className="w-full flex flex-col items-center justify-center gap-6 lg:gap-12 lg:mt-8">
				<RoundsSettingsCard min={1} max={5} interval={1} defaultValue={3} settings="Rounds" onChange={onRoundsChange}></RoundsSettingsCard>
				<RoundsSettingsCard min={10} max={120} interval={10} defaultValue={60} settings="Time Limits" formatValue={formatTime} onChange={onTimeChange}></RoundsSettingsCard>
			</div>
				<p className={` text-center text-red-500 text-label-small transition-opacity ${
				errorMessage ? "opacity-100" : "opacity-0"
				}`}>
				{errorMessage ?? "\u00A0"}
				</p>
		</div>
	)
}