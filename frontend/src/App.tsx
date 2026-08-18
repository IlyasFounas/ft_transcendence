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
          type="sdsdsl"
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

    </div>
  );
}

export default App;