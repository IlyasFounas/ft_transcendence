import OperationHistory from './components/cards/OperationHistory.tsx'
import { operations } from "./data/operations";

function App() {

  return (
	< div className="w-full min-h-screen bg-white">
      <OperationHistory operations={operations}/>
	</div>
  )
}

export default App
