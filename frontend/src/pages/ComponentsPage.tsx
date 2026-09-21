import {useState} from 'react'
import OperationHistory from '../components/cards/OperationHistory';
import { operations } from "../data/operations";
import UserProfile from '../components/user/UserProfile.tsx';
import AddFriendRow from '../components/rows/AddFriendRow.tsx';
import SmallButton from '../components/buttons/SmallButton.tsx';
import FriendRequestRow from '../components/rows/FriendRequestRow.tsx';
import StatusBadge from '../components/primitives/StatusBadge.tsx';
import UserStatus from '../components/user/UserStatus.tsx';
import FriendListRow from '../components/rows/FriendListRow.tsx';
import AuthenticationCard from '../components/cards/AuthenticationCard.tsx';
import AuthenticationWelcome from '../components/authentication/AuthenticationWelcome.tsx';
import Input from '../components/primitives/Input.tsx';
import AuthenticationSignUp from '../components/authentication/AuthenticationSignUp.tsx';
import AuthenticationLogin from '../components/authentication/AuthenticationLogin.tsx';
import PrimaryButton from '../components/buttons/PrimaryButton.tsx';
import SecondaryButton from '../components/buttons/SecondaryButton.tsx';
import Avatar from '../components/primitives/Avatar.tsx';
import ChatBoxCard from '../components/cards/ChatBox.tsx';
import Podium from '../components/cards/podium/Podium.tsx';

export default function ComponentPage() {

	const [email, setEmail] = useState("");

  return (
      <div className="w-full min-h-screen bg-black p-8 text-white">
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Operation History</h2>
          <OperationHistory operations={operations} />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">User Profile</h2>
          <UserProfile
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Small Button — Default</h2>
          <SmallButton onClick={() => console.log("clicked")}>
            Add Friend
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Small Button — Danger</h2>
          <SmallButton variant="danger">
            remove
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Small Button — Primary Disabled</h2>
          <SmallButton variant="primary" disabled>
            Add Friend
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Small Button — Danger Disabled</h2>
          <SmallButton variant="danger" disabled>
            Add Friend
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Add Friend Row</h2>
          <AddFriendRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Add Friend Row — Second Example</h2>
          <AddFriendRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Friend Request Row</h2>
          <FriendRequestRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Friend Request Row — Long Username
          </h2>
          <FriendRequestRow
            username="agent tommyyyyyyyyyy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Status Badge — Ready</h2>
          <StatusBadge status="ready" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">User Status — Online</h2>
          <UserStatus status="online" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">User Status — In Mission</h2>
          <UserStatus status="in-mission" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">User Status — Offline</h2>
          <UserStatus status="offline" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Friend List Row — Online</h2>
          <FriendListRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            status="online"
            rank={233}
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Friend List Row — In Mission</h2>
          <FriendListRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            status="in-mission"
            rank={1212}
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Friend List Row — Offline</h2>
          <FriendListRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            status="offline"
            rank={33}
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">Input</h2>
          <Input
            label="Email"
            type="email"
            placeholder="example@mail.com"
			value={email}
    		onChange={setEmail}
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Authentication Card — Welcome
          </h2>
          <AuthenticationCard>
            <AuthenticationWelcome />
          </AuthenticationCard>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Authentication Card — Sign Up
          </h2>
          <AuthenticationCard>
            <AuthenticationSignUp onLogin={() => {}}/>
          </AuthenticationCard>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Authentication Card — Login
          </h2>
          <AuthenticationCard>
            <AuthenticationLogin onBackToOptions={() => {}}/>
          </AuthenticationCard>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Primary Button - without subtitles
          </h2>
          <PrimaryButton className='w-[500px] h-[125px]'>Join Operation</PrimaryButton>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold ">
            Primary Button - with subtitles
          </h2>
          <PrimaryButton subtitle="Launch a new operation" className='w-[500px] h-[125px]'>Start Mission</PrimaryButton>
        </section>
  
       <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Secondary Button - Inactive
          </h2>
          <SecondaryButton>Play</SecondaryButton>
        </section>
  
       <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Secondary Button - Active
          </h2>
          <SecondaryButton active={true}>Play</SecondaryButton>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Secondary Button - Active
          </h2>
          <SecondaryButton active={true}>Play</SecondaryButton>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Avatar
          </h2>
        <Avatar avatar="/avatars/avatar1.png"></Avatar>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            ChatBox - Received
          </h2>
            <ChatBoxCard variant="received" message="Hello ! How are you doing?" time="02.12 PM"></ChatBoxCard>
        </section>
  
  
      <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            ChatBox - Sent
          </h2>
            <ChatBoxCard variant="sent" message="I'm good, how r u?" time="02.11 PM"></ChatBoxCard>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Podium Card - Gold
          </h2>
            <Podium color="gold" username="Agent Tommy" title="Field Analyst" points={2.563}></Podium>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Podium Card - Silver
          </h2>
            <Podium color="silver" username="Agent Tommy" title="Field Analyst" points={2.563}></Podium>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-lg font-bold">
            Podium Card - Bronze
          </h2>
            <Podium color="bronze" username="Agent Tommy" title="Field Analyst" points={2.563}></Podium>
        </section>
      </div>
    );
}