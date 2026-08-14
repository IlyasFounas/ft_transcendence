import OperationHistory from './components/cards/OperationHistory.tsx'
import { operations } from "./data/operations";
import UserProfile from './components/user/UserProfile.tsx';
import AddFriendRow from './components/rows/AddFriendRow.tsx';
import SmallButton from './components/buttons/SmallButton.tsx';
import FriendRequestRow from './components/rows/FriendRequestRow.tsx';
import StatusBadge from './components/primitives/StatusBadge.tsx';
import UserStatus from './components/user/UserStatus.tsx';
import FriendListRow from './components/rows/FriendListRow.tsx';
import FriendChatRow from './components/rows/FriendChatRow.tsx';
import AuthenticationCard from './components/cards/AuthenticationCard.tsx';
import AuthenticationWelcome from './components/authentication/AuthenticationWelcome.tsx';
import Input from './components/primitives/Input.tsx';
import AuthenticationSignUp from './components/authentication/AuthenticationSignUp.tsx';
import AuthenticationLogin from './components/authentication/AuthenticationLogin.tsx';

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
      <StatusBadge status='ready'/>   
      <UserStatus status='online'/>
      <UserStatus status='in-mission'/>   
      <UserStatus status='offline'/>
      <FriendListRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png' status='online' rank={233}/>  
      <FriendListRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png' status='in-mission' rank={1212}/>  
      <FriendListRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png' status='offline' rank={33}/>  
      <FriendChatRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png' lastMessage="make[1]: Entering directory '/home/spunyapr/commoncore/ft_transcendence/ft_transcendence'" timestamp='03.45pm' status='online'/>
      <FriendChatRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png' lastMessage="make[1]: Entering directory '/home/spunyapr/commoncore/ft_transcendence/ft_transcendence'" timestamp='03.45pm' status='offline'/>
      <FriendChatRow username="agent tommy" role="Field Analyst" avatar='/avatars/avatar1.png' lastMessage="make[1]: Entering directory '/home/spunyapr/commoncore/ft_transcendence/ft_transcendence'" timestamp='03.45pm' status='in-mission'/>
      <AuthenticationCard>
        <AuthenticationWelcome />
      </AuthenticationCard>
      <Input label="Email" type="sdsdsl" placeholder='example@mail.com'></Input>
      <AuthenticationCard>
        <AuthenticationSignUp />
      </AuthenticationCard>
      <AuthenticationCard>
        <AuthenticationLogin />
      </AuthenticationCard>
  </div>
  )
}

export default App
