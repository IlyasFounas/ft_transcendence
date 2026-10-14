import { CircleQuestionMark, LogOut } from 'lucide-react'
import LogoHorizontal from '../primitives/LogoHorizontal'
import Countdown from './Countdown'
import GameInfosCard from './GameInfosCard'
import SignalIntegrityCard from './SignalIntegrityCard'
import LogoIcon from "../primitives/LogoIcon"

type Props = {
	countdown: number;
	totalTime: number;
	gameId: number;
	round: number;
	totalRounds: number;
	onLeave: () => void;
	onTuto: () => void;
};

export default function GameBar({countdown, totalTime, gameId, round, totalRounds, onLeave, onTuto} : Props)
{
	return (
	<>
		<div
			className="
				hidden lg:grid
				grid-cols-[1fr_auto_1fr]
				w-full
				h-[100px]
				items-center
				px-8
				gap-4
			"
		>
			<div className="flex justify-start items-center gap-4">
				<LogoHorizontal className="hidden xl:flex min-w-[100px] h-[100px]" />
				<LogoIcon className="flex xl:hidden w-[75px] h-[75px]" />

				<GameInfosCard
					round={round}
					totalRounds={totalRounds}
					id={gameId}
				/>
			</div>

			<div className="flex justify-center">
				<Countdown countdown={countdown} />
			</div>

			<div className="flex justify-end gap-6">
				<SignalIntegrityCard
					countdown={countdown}
					totalTime={totalTime}
				/>

				<div className="flex flex-row gap-6">
					<button
						className="cursor-pointer transition-transform duration-200 ease-out hover:scale-[1.2] active:scale-[0.9]"
						onClick={onTuto}
					>
						<CircleQuestionMark />
					</button>

					<button
						className="cursor-pointer transition-transform duration-200 ease-out hover:scale-[1.2] hover:text-red-500 active:scale-[0.9]"
						onClick={onLeave}
					>
						<LogOut />
					</button>
				</div>
			</div>
		</div>

		<div
			className="
				flex lg:hidden
				items-center
				justify-between
				w-full
				px-4
				py-3
			"
		>
			<button
				className="cursor-pointer transition-transform duration-200 ease-out hover:scale-[1.15] active:scale-[0.9]"
				onClick={onTuto}
			>
				<CircleQuestionMark />
			</button>

			<Countdown countdown={countdown} />

			<button
				className="cursor-pointer transition-transform duration-200 ease-out hover:scale-[1.15] hover:text-red-500 active:scale-[0.9]"
				onClick={onLeave}
			>
				<LogOut />
			</button>
		</div>
	</>
	)
}