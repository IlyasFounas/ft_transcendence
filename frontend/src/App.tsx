import OperationHistory from './components/cards/OperationHistory.tsx'
import { operations } from "./data/operations";
import UserProfile from './components/user/UserProfile.tsx';
import AddFriendRow from './components/rows/AddFriendRow.tsx';
import SmallButton from './components/buttons/SmallButton.tsx';
import PrimaryButton from './components/buttons/PrimaryButton.tsx';
function App() {

  return (
	< div className="w-full min-h-screen bg-black">
      <OperationHistory operations={operations}/>
      <UserProfile username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png'/>
      <AddFriendRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png'/>
      <SmallButton onClick={() => console.log("clicked")}>Add Friend</SmallButton>
      <SmallButton variant='danger'>remove</SmallButton>
      <SmallButton
          variant="primary"
          disabled
      >
          Add Friend
      </SmallButton>
      <SmallButton
          variant="danger"
          disabled
      >
          Add Friend
      </SmallButton>
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
		</div>
  </div>
  )
}

export default App
