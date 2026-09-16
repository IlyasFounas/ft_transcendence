import RoundsSettingsCard from "../cards/RoundsSettingsCard";

interface CreateOperationCardProps {
	onRoundsChange: (value: number) => void;
  	onTimeChange: (value: number) => void;
}

export default function CreateOperationCard({onRoundsChange, onTimeChange} : CreateOperationCardProps)
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
		<div className="w-full h-full flex flex-col justify-center items-center gap-6 lg:gap-8 text-heading-small lg:text-heading-medium uppercase">
				<RoundsSettingsCard min={1} max={5} interval={1} settings="Rounds" onChange={onRoundsChange}></RoundsSettingsCard>
				<RoundsSettingsCard min={10} max={120} interval={10} settings="Time Limits" formatValue={formatTime} onChange={onTimeChange}></RoundsSettingsCard>
		</div>
	)
}