import { useState } from 'react'
import SecondaryButton from './components/buttons/SecondaryButton.tsx'
import PrimaryButton from './components/buttons/PrimaryButton.tsx'
import ChatBoxCard from './components/cards/ChatBox.tsx'
import GameResultCard from './components/cards/GameResultCard.tsx'
//import LobbyCard from './components/LobbyCard.tsx'
import Avatar from './components/primitives/Avatar.tsx'

import './App.css'
import StatusBadge from './components/primitives/StatusBadge.tsx'

function App() {

  return (
    <div className="w-200 h-2000 bg-black g-50 p-25 inline-flex flex-col" >
		<PrimaryButton>Join Operation</PrimaryButton>
		<SecondaryButton>icon</SecondaryButton>
		<SecondaryButton disabled>icon</SecondaryButton>
    </div>
  )
}

export default App
