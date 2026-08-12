import OperationHistory from './components/cards/OperationHistory.tsx'
import { operations } from "./data/operations";
import UserProfile from './components/user/UserProfile.tsx';

function App() {

  return (
	< div className="w-full min-h-screen bg-black">
      <OperationHistory operations={operations}/>
      <UserProfile username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png' />
	</div>
  )
}

export default App
