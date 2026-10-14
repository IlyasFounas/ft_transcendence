interface ResultInfoRowProps{
	title: string;
	infos?: string;
	suffix?: string;
	variant: "primary" | "secondary";
	animate? : boolean;
}

const styles =
{
	primary:
	{
		title: "text-heading-small lg:text-heading-large",
		infos: "text-label-small lg:text-heading-medium"
	},

	secondary: 
	{
		title: "text-label-small lg:text-heading-medium",
		infos: "text-label-small lg:text-heading-small",
	}
}

export default function ResultInfoRow({title, infos, variant, suffix, animate} : ResultInfoRowProps)
{
	const Style = styles[variant];

	return (
		<div className="w-full">
			<div className="
				py-2 text-text-primary
				flex flex-row
				justify-between items-center px-[10px]
				uppercase
				outline 
				outline-1 
				outline-offset-[-1px] 
				outline-sky-950
				rounded-[10px]
				bg-surface-base
				">
				<h1 className={`${Style.title}`}>{title}</h1>
				<div className="flex flex-row gap-2">
					<p className={`${Style.infos} text-primary ${animate ? "animate-scale-pop" : ""}`}>{infos}</p>
					<p className={`${Style.infos}`} >{suffix}</p>
				</div>
			</div>
		</div>
	)
}