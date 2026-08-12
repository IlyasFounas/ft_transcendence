import OperationHistory from './components/cards/OperationHistory.tsx'
import { operations } from "./data/operations";
import UserProfile from './components/user/UserProfile.tsx';
import AddFriendRow from './components/rows/AddFriendRow.tsx';
import SmallButton from './components/buttons/SmallButton.tsx';
import FriendRequestRow from './components/rows/FriendRequestRow.tsx';

function App() {

  return (
	< div className="w-full min-h-screen bg-black">
      <OperationHistory operations={operations}/>
      <UserProfile username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png'/>
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
      <AddFriendRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png'/>
      <AddFriendRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png'/>
      <FriendRequestRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png'/>
      <FriendRequestRow username="agent tommyyyyyyyyyy" role="Field Analyst" avatar='/avatars/avatar1.png'/>
      
  </div>
  )
}

export default App
