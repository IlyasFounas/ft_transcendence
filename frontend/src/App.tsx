import { useState } from 'react'
import Button from './components/Button/Button.tsx'

import './App.css'

function App() {

  return (
    <div className="w-500 h-500 bg-black">
		<Button>Play</Button>
		<Button disabled>Play</Button>

		<Button icon={<span>→</span>}>
        Avec icon
      </Button>
	  <Button disabled icon={<span>→</span>}>
        Avec icon
      </Button>
    </div>
  )
}

export default App
