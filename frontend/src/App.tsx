import { useState } from 'react'
import Avatar from './components/primitives/Avatar'
import LobbyPage from './pages/LobbyPage'
import RoomCodeCard from './components/cards/RoomCodeCard'
import Panel from './components/primitives/Panel'
import LobbyTestPage from './components/cards/LobbyTestPage'
import PrimaryButton from './components/buttons/PrimaryButton'
import Podium from './components/cards/podium/Podium'

function App() {

  return (
	<div>
		<Podium color="bronze" username="lulu" title="la goat" points={2584}></Podium>
		<Podium color="silver" username="lulu" title="la goat" points={2584}></Podium>
		<Podium color="gold" username="lulu" title="la goat" points={2584}></Podium>
	</div>
  )
}
export default App
