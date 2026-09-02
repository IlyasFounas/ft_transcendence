import RoundsSettingsCard from "../cards/RoundsSettingsCard";
import CheckBox from "./CheckBox";

export default function CreateOperationCard()
{

	function formatTime(value : number)
	{
		if (value % 1 === 0.5)
			 return `${Math.floor(value)} min 30`;

   		return `${value} min`;
	}

	return (
		<div className="w-full h-full flex flex-col justify-center items-center gap-6 lg:gap-8 text-heading-small lg:text-heading-medium uppercase">
				<RoundsSettingsCard min={1} max={10} interval={1} settings="Rounds"></RoundsSettingsCard>
				<RoundsSettingsCard min={0.5} max={5} interval={0.5} settings="Time Limits" formatValue={formatTime}></RoundsSettingsCard>
				<div className="flex flex-row gap-6"> 
					<CheckBox>Zoom</CheckBox>
					<CheckBox>Move</CheckBox>
				</div>
		</div>
	)
}