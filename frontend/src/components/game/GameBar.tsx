import {useState} from "react"

import { Cog, LogOut } from 'lucide-react'

import LogoHorizontal from '../primitives/LogoHorizontal'
import Countdown from './Countdown'
import GameInfosCard from './GameInfosCard'
import SignalIntegrityCard from './SignalIntegrityCard'

export default function GameBar()
{
	const totalTime = 60;
	const [countDown, setCountDown] = useState(totalTime);

	return (
		<div className="
			w-full
			flex flex-row justify-center items-center px-8 gap-8">
			<LogoHorizontal className="min-w-[150px]"></LogoHorizontal>
			<GameInfosCard round={3} totalRounds={5} id={4242}></GameInfosCard>
			<Countdown countDown={countDown} setCountDown={setCountDown}></Countdown>
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
	)
}