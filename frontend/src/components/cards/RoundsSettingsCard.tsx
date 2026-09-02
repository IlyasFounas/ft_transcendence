import { useState } from "react"
import { Minus, Plus } from "lucide-react"

interface RoundsSettingsCardProps{
	min: number;
	max: number; 
	interval: number;
	settings: string;
	formatValue?: (value: number) => string;
}

export default function RoundsSettingsCard({min, max, interval, settings, formatValue} : RoundsSettingsCardProps) 
{
	const [value, setValue] = useState(min);

	return (
		<div>
			<div className="flex items-center text-heading-small lg:text-heading-medium">
				<div className="w-[150px] lg:w-[200px]">
					{settings}
				</div>
				<div className="flex flex-rows justify-between w-[150px] lg:w-[200px]">

					<button
						className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						active:scale-[0.98]
						"
						onClick={() => {
							if (value > min) {
								setValue(value - interval);
							}
						}}
						>
						<Minus />
					</button>

					<div>
						 {formatValue ? formatValue(value) : value}
					</div>
					<button
						className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						active:scale-[0.98]
						"
						onClick={() => {
							if (value < max) {
								setValue(value + interval);
							}
						}}
						>
						<Plus />
					</button>
				</div>
			</div>
		</div>
	)
}