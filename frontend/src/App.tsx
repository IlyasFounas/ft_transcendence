import { useState } from 'react'
import SecondaryButton from './components/Button/SecondaryButton.tsx'
import PrimaryButton from './components/Button/PrimaryButton.tsx'

import './App.css'

function App() {

  return (
    <div className="w-2000 h-2000 bg-black g-50 p-25 flex flex-row" >
		<div className="flex flex-col p-25">
			<PrimaryButton size="large" subtitle="Launch a new operation">Start Mission</PrimaryButton>
			<PrimaryButton size="medium" subtitle="Launch a new operation">Start Mission</PrimaryButton>
			<PrimaryButton size="small" subtitle="Launch a new operation">Start Mission</PrimaryButton>
		</div>
		<div className="flex flex-col p-25">
			<PrimaryButton size="large">Join Operation</PrimaryButton>
			<PrimaryButton size="medium">Join Operation</PrimaryButton>
			<PrimaryButton size="small">Join Operation</PrimaryButton>
		</div>
		
    </div>
  )
}

export default App
