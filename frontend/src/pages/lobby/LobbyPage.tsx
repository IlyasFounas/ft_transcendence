import { useState } from "react";
import SecondaryButton from "../../components/buttons/SecondaryButton";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import CreateOperationCard from "../../components/lobby/CreateOperationCard";
import JoinOperationCard from "../../components/lobby/JoinOperationCard";

export default function LobbyPage()
{
	const [mode, setMode] = useState<"create" | "join">("join");
	

	return (
		<div className="min-h-screen w-full bg-[url(/backgrounds/earth_background.png)] py-50">
			<div className="flex min-h-[50px] items-center justify-center gap-8">
				<SecondaryButton
					active={mode === "create"}
					className="px-[12px] text-heading-small lg:text-heading-large"
					onClick={() => setMode("create")}
				>
					Create Operation
				</SecondaryButton>

				<SecondaryButton
					active={mode === "join"}
					className="text-heading-small lg:text-heading-large"
					onClick={() => setMode("join")}
				>
					Join Operation
				</SecondaryButton>
			</div>
			<div className="
						flex
						h-50 lg:h-[300px]
						w-[calc(100%-2rem)]
						max-w-[500px]
						mx-auto m-8
						rounded-lg
						border
						border-surface-elevated
						bg-[#020B16]
						opacity-90">
					{mode === "create" ? (
						<CreateOperationCard/>
					) : (
						<JoinOperationCard/>
					)}
			</div>
			<div className="flex justify-center">
					{mode === "create" ? (
						<PrimaryButton size="small" className="w-300">Create Operation</PrimaryButton>
					) : (
						<PrimaryButton size="small">Join Operation</PrimaryButton>
					)}
			</div>
			
	</div>
	)
}