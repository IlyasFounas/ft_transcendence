import OperationHistory from '../components/cards/OperationHistory';
import { operations } from "../data/operations";
import UserProfile from '../components/user/UserProfile.tsx';
import AddFriendRow from '../components/rows/AddFriendRow.tsx';
import SmallButton from '../components/buttons/SmallButton.tsx';
import FriendRequestRow from '../components/rows/FriendRequestRow.tsx';
import StatusBadge from '../components/primitives/StatusBadge.tsx';
import UserStatus from '../components/user/UserStatus.tsx';
import FriendListRow from '../components/rows/FriendListRow.tsx';
import FriendChatRow from '../components/rows/FriendChatRow.tsx';
import AuthenticationCard from '../components/cards/AuthenticationCard.tsx';
import AuthenticationWelcome from '../components/authentication/AuthenticationWelcome.tsx';
import Input from '../components/primitives/Input.tsx';
import AuthenticationSignUp from '../components/authentication/AuthenticationSignUp.tsx';
import AuthenticationLogin from '../components/authentication/AuthenticationLogin.tsx';
import PrimaryButton from '../components/buttons/PrimaryButton.tsx';
import SecondaryButton from '../components/buttons/SecondaryButton.tsx';
import Avatar from '../components/primitives/Avatar.tsx';
import GameResultCard from '../components/cards/GameResultCard.tsx';
import LobbyTestPage from '../components/cards/LobbyTestPage.tsx';
import ChatBoxCard from '../components/cards/ChatBox.tsx';
import SummaryResultCard from '../components/cards/SummaryResultCard.tsx'
import Podium from '../components/cards/podium/Podium.tsx';

export default function ComponentPage() {
  return (
      <div className="w-full min-h-screen bg-black p-8 text-white">
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Operation History</h2>
          <OperationHistory operations={operations} />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">User Profile</h2>
          <UserProfile
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Small Button — Default</h2>
          <SmallButton onClick={() => console.log("clicked")}>
            Add Friend
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Small Button — Danger</h2>
          <SmallButton variant="danger">
            remove
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Small Button — Primary Disabled</h2>
          <SmallButton variant="primary" disabled>
            Add Friend
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Small Button — Danger Disabled</h2>
          <SmallButton variant="danger" disabled>
            Add Friend
          </SmallButton>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Add Friend Row</h2>
          <AddFriendRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Add Friend Row — Second Example</h2>
          <AddFriendRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Friend Request Row</h2>
          <FriendRequestRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Friend Request Row — Long Username
          </h2>
          <FriendRequestRow
            username="agent tommyyyyyyyyyy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Status Badge — Ready</h2>
          <StatusBadge status="ready" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">User Status — Online</h2>
          <UserStatus status="online" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">User Status — In Mission</h2>
          <UserStatus status="in-mission" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">User Status — Offline</h2>
          <UserStatus status="offline" />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Friend List Row — Online</h2>
          <FriendListRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            status="online"
            rank={233}
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Friend List Row — In Mission</h2>
          <FriendListRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            status="in-mission"
            rank={1212}
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Friend List Row — Offline</h2>
          <FriendListRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            status="offline"
            rank={33}
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Friend Chat Row — Online</h2>
          <FriendChatRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            lastMessage="make[1]: Entering directory '/home/spunyapr/commoncore/ft_transcendence/ft_transcendence'"
            timestamp="03.45pm"
            status="online"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Friend Chat Row — Offline</h2>
          <FriendChatRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            lastMessage="make[1]: Entering directory '/home/spunyapr/commoncore/ft_transcendence/ft_transcendence'"
            timestamp="03.45pm"
            status="offline"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Friend Chat Row — In Mission
          </h2>
          <FriendChatRow
            username="agent tommy"
            role="Field Analyst"
            avatar="/avatars/avatar1.png"
            lastMessage="make[1]: Entering directory '/home/spunyapr/commoncore/ft_transcendence/ft_transcendence'"
            timestamp="03.45pm"
            status="in-mission"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">Input</h2>
          <Input
            label="Email"
            type="email"
            placeholder="example@mail.com"
          />
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Authentication Card — Welcome
          </h2>
          <AuthenticationCard>
            <AuthenticationWelcome />
          </AuthenticationCard>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Authentication Card — Sign Up
          </h2>
          <AuthenticationCard>
            <AuthenticationSignUp />
          </AuthenticationCard>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Authentication Card — Login
          </h2>
          <AuthenticationCard>
            <AuthenticationLogin />
          </AuthenticationCard>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Primary Button - without subtitles
          </h2>
          <PrimaryButton>Join Operation</PrimaryButton>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Primary Button - with subtitles
          </h2>
          <PrimaryButton subtitle="Launch a new operation">Start Mission</PrimaryButton>
        </section>
  
       <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Secondary Button - Inactive
          </h2>
          <SecondaryButton>Play</SecondaryButton>
        </section>
  
       <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Secondary Button - Active
          </h2>
          <SecondaryButton active={true}>Play</SecondaryButton>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Secondary Button - Active
          </h2>
          <SecondaryButton active={true}>Play</SecondaryButton>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Avatar
          </h2>
        <Avatar avatar="/avatars/avatar1.png"></Avatar>
        </section>
  
       <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Game Result Card
          </h2>
            <GameResultCard rank={2} distance={256} points={1000} ></GameResultCard>
        </section>
  
       <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Lobby Card
          </h2>
            <LobbyTestPage></LobbyTestPage>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            ChatBox - Received
          </h2>
            <ChatBoxCard variant="received" message="Hello ! How are you doing?" time="02.12 PM"></ChatBoxCard>
        </section>
  
  
      <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            ChatBox - Sent
          </h2>
            <ChatBoxCard variant="sent" message="I'm good, how r u?" time="02.11 PM"></ChatBoxCard>
        </section>
  
       <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Summary Result Roll
          </h2>
            <SummaryResultCard rank={2}
                player={{ username: "Agent Tommy", role: "Field Analyst", avatar: "/avatars/avatar1.png" }}
             rounds={[ {round: 1, points: 1820,}, {round: 1, points: 1245,},{round: 1, points: 180, },{round: 1, points: 1896,},{round: 1, points: 1200,}, ]}></SummaryResultCard>
        </section>
  
      <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Podium Card - Gold
          </h2>
            <Podium color="gold" username="Agent Tommy" title="Field Analyst" points={2.563}></Podium>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Podium Card - Silver
          </h2>
            <Podium color="silver" username="Agent Tommy" title="Field Analyst" points={2.563}></Podium>
        </section>
  
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-bold">
            Podium Card - Bronze
          </h2>
            <Podium color="bronze" username="Agent Tommy" title="Field Analyst" points={2.563}></Podium>
        </section>
      </div>
    );
}