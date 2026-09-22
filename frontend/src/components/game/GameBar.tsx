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
};

export default function GameBar({countdown, totalTime, gameId, round, totalRounds} : Props)
{
	return (
		<div className="
			grid
			grid-cols-[1fr_auto_1fr]
			w-full
			h-[120px] md:[100px]
			items-center
			px-8
			gap-8 md:gap-4
		">
			<div className="flex justify-start items-center gap-4">
				<div className="
					cursor-pointer
					transition-transform
					duration-200
					ease-out
					hover:scale-[1.2]
						
					active:scale-[0.9]
					">
						<LogoHorizontal className="hidden xl:flex min-w-[100px] h-[100px] md:w-[50px]"></LogoHorizontal>
						<LogoIcon className="flex xl:hidden w-[75px] h-[75px]"></LogoIcon>
				</div>
				<GameInfosCard round={round} totalRounds={totalRounds} id={gameId}></GameInfosCard>
			</div>
			<div className="flex justify-center">
				<Countdown countdown={countdown}/>
			</div>
			<div className="flex justify-end gap-6">
				<SignalIntegrityCard countdown={countdown} totalTime={totalTime}></SignalIntegrityCard>
				<div className="flex flex-row gap-6">
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]

						active:scale-[0.9]">
						<CircleQuestionMark></CircleQuestionMark>
					</button>
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						transition-colors
						hover:scale-[1.2]
						hover:text-red-900
						hover:brightness-100
						
						active:scale-[0.9]
						">
						<LogOut></LogOut>
					</button>
				</div>
			</div>
			</div>
	)
}