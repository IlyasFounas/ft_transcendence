import { useState } from 'react'
import LobbyCard from './components/cards/LobbyCard.tsx'
import SummaryResultCard from './components/cards/SummaryResultCard.tsx'
import GoldPodiumCard from './components/cards/podium/GoldPodiumCard.tsx'
import PrimaryButton from './components/buttons/PrimaryButton.tsx'
import SecondaryButton from './components/buttons/SecondaryButton.tsx'
import Avatar from './components/primitives/Avatar.tsx'
import ChatBoxCard from './components/cards/ChatBox.tsx'
import GameResultCard from './components/cards/GameResultCard.tsx'

function App() {

  return (
	< div className="w-2000 h-2000 bg-black">
		<div className = "flex flex-row">
			<div className="flex flex-col">
				<PrimaryButton size="large" subtitle="Launch a new operation">Join Operation</PrimaryButton>
				<PrimaryButton size="medium" subtitle="Launch a new operation">Join Operation</PrimaryButton>
				<PrimaryButton size="small" subtitle="Launch a new operation">Join Operation</PrimaryButton>
			</div>
			<div className="flex flex-col">
				<PrimaryButton size="large">Join Operation</PrimaryButton>
				<PrimaryButton size="medium">Join Operation</PrimaryButton>
				<PrimaryButton size="small">Join Operation</PrimaryButton>
			</div>
			<div className="flex flex-col">
				<PrimaryButton size="large" disabled>Join Operation</PrimaryButton>
				<PrimaryButton size="medium" disabled>Join Operation</PrimaryButton>
				<PrimaryButton size="small" disabled>Join Operation</PrimaryButton>
			</div>
		</div>
		<div className="flex flex-row">
				<SecondaryButton>Play</SecondaryButton>
				<SecondaryButton icon={<span>→</span>}>Play</SecondaryButton>
				<SecondaryButton disabled>Play</SecondaryButton>
				<SecondaryButton icon={<span>→</span>} disabled>Play</SecondaryButton>
			<div className="flex flex-row m-8 gap-8">
				<Avatar></Avatar>
				<Avatar size={90}></Avatar>
				<Avatar size={70} borderColor="#FFB800"></Avatar>
			</div>
		</div>
		<div className="flex flex_row">
		<div className= "flex flex-col gap-5 m-8">
			<GameResultCard rank={2} distance={256} points={1000} ></GameResultCard>
			<LobbyCard username="Agent Nova" role="Field Analyst"></LobbyCard>
			<SummaryResultCard rank={2}
  			player={{ username: "Agent Tommy", role: "Field Analyst", avatar: "/avatars/avatar1.png" }}
 			 rounds={[ {round: 1, points: 1820,}, {round: 1, points: 1245,},{round: 1, points: 180, },{round: 1, points: 1896,},{round: 1, points: 1200,}, ]}></SummaryResultCard>
		</div>

		<div className= "flex flex-col gap-5 m-8">
			<ChatBoxCard variant="received" message="Hello ! How are you doing?" time="02.12 PM"></ChatBoxCard>
			<ChatBoxCard variant="sent" message="I'm good, how r u?" time="02.11 PM"></ChatBoxCard>
		</div>
		</div>
	</div>
  )
}

export default App
