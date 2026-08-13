import { useState } from 'react'
import LobbyCard from './components/cards/LobbyCard.tsx'
import SummaryResultCard from './components/cards/SummaryResultCard.tsx'
import PrimaryButton from './components/buttons/PrimaryButton.tsx'
import SecondaryButton from './components/buttons/SecondaryButton.tsx'
import Avatar from './components/primitives/Avatar.tsx'
import ChatBoxCard from './components/cards/ChatBox.tsx'
import GameResultCard from './components/cards/GameResultCard.tsx'
import { PodiumShapeSvg } from './components/cards/podium/PodiumShape.tsx'
import { AvatarFrame } from './components/cards/podium/AvatarFrame.tsx'
import PodiumCard from './components/cards/podium/PodiumCard.tsx'

function App() {

  return (
	< div className="w-2000 h-2000 bg-black">

    <div className="flex gap-4">
      <PodiumCard
        variant="gold"
        rank={1}
        name="Agent Nova"
        role="Field Analyst"
        score={2025}
        imageSrc="avatars/avatar1.png"
      />
      <PodiumCard
        variant="silver"
        rank={2}
        name="Agent Orion"
        role="Field Analyst"
        score={1830}
        imageSrc="avatars/avatar1.png"
      />
      <PodiumCard
        variant="bronze"
        rank={3}
        name="Agent Vega"
        role="Field Analyst"
        score={1690}
        imageSrc="avatars/avatar1.png"
      />
    </div>
	</div>
  )
}

export default App
