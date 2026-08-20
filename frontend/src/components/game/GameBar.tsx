import {useState} from "react"

import { Cog, LogOut } from 'lucide-react'

import LogoHorizontal from '../primitives/LogoHorizontal'
import Countdown from './Countdown'
import GameInfosCard from './GameInfosCard'
import SignalIntegrityCard from './SignalIntegrityCard'
import LogoIcon from "../primitives/LogoIcon"

export default function GameBar()
{
	const totalTime = 60;
	const [countDown, setCountDown] = useState(totalTime);

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
				<GameInfosCard round={3} totalRounds={5} id={4242}></GameInfosCard>
			</div>
			<div className="flex justify-center">
				<Countdown countDown={countDown} setCountDown={setCountDown}
				></Countdown>
			</div>
			<div className="flex justify-end gap-6">
				<SignalIntegrityCard countDown={countDown} totalTime={totalTime}></SignalIntegrityCard>
				<div className="flex flex-row gap-6">
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						hover:-rotate-60
						
						active:scale-[0.9]">
						<Cog></Cog>
					</button>
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						
						active:scale-[0.9]
						">
						<LogOut></LogOut>
					</button>
				</div>
			</div>
		</div>
	)
}