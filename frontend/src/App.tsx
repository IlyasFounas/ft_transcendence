import { useState } from 'react'
import LobbyCard from './components/cards/LobbyCard.tsx'
import SummaryResultCard from './components/cards/SummaryResultCard.tsx'
import GoldPodiumCard from './components/cards/podium/GoldPodiumCard.tsx'

function App() {

  return (
	< div className="w-2000 h-200 bg-black">
		<GoldPodiumCard></GoldPodiumCard>
	</div>
  )
}

export default App
