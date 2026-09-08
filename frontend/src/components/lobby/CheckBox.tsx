import { useState } from "react"

export default function CheckBox({children}: { children: string})
{
	const [isChecked, setIsChecked] = useState(true);

	return (
		<div>
			<label className="flex gap-2">
					<input	type="checkbox" 
						checked={isChecked}
						onChange={(e) => setIsChecked(e.target.checked)}>
					</input>
					{children}
			   </label>
		</div>
	)
}