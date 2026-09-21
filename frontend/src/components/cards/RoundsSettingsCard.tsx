import { useState, useEffect } from "react"
import { Minus, Plus } from "lucide-react"

interface RoundsSettingsCardProps{
	min: number;
	max: number; 
	interval: number;
	defaultValue: number;
	settings: string;
	formatValue?: (value: number) => string;
	onChange: (value: number) => void;
}

export default function RoundsSettingsCard({min, max, interval, defaultValue, settings, formatValue, onChange} : RoundsSettingsCardProps) 
{
	const [value, setValue] = useState(defaultValue);

	function updateValue(newValue: number) {
		setValue(newValue);
		onChange(newValue);
	}

	useEffect(() => {
    onChange(defaultValue);
  	}, []);

	return (
		<div>
			<div className="flex items-center text-heading-small lg:text-heading-medium">
				<div className="w-[150px] lg:w-[200px]">
					{settings}
				</div>
				<div className="flex flex-rows justify-between w-[150px] lg:w-[200px]">

					<button
						className={`
						transition-transform
						duration-200
						ease-out
						${value > min ? "text-white hover:scale-[1.2] active:scale-[0.98] cursor-pointer" : "text-gray-600"}
						`}
						onClick={() => {
							if (value > min) {
								updateValue(value - interval);
							}
						}}
						>
						<Minus />
					</button>

					<div>
						 {formatValue ? formatValue(value) : value}
					</div>
					<button
						className={`
						transition-transform
						duration-200
						ease-out
						${value < max ? "text-white hover:scale-[1.2] active:scale-[0.98] cursor-pointer" : "text-gray-600"}
						`}
						onClick={() => {
							if (value < max) {
								updateValue(value + interval);
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