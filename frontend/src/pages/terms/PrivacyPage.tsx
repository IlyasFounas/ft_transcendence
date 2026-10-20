import { Link } from 'react-router-dom';

export default function PrivacyPage () {
  return (
		<main className="mx-auto max-w-3xl px-6 py-10 text-gray-200">
		<Link to="/" className="mb-6 inline-block text-sm text-gray-400 hover:text-white">
			← Back to home
		</Link>
      <h1 className="mb-2 text-3xl font-bold">Privacy Policy</h1>
      <p className="mb-8 text-sm text-gray-400">Last updated: October 8, 2026</p>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">1. Introduction</h2>
        <p>
          ATLAS is a multiplayer geography guessing game created as a school project at
          École 42. This policy explains what personal data we collect, why we collect it,
          and what rights you have over it.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">2. Data we collect</h2>
        <ul className="list-disc space-y-1 pl-6">
          <li>Account data: username, email address and a hashed password.</li>
          <li>Profile data: your avatar and your title (rank).</li>
          <li>Game data: your games, scores, guesses and leaderboard position.</li>
          <li>Social data: your friends list and the chat messages you send.</li>
          <li>Technical data: session information needed to keep you logged in.</li>
        </ul>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">3. How we use your data</h2>
        <ul className="list-disc space-y-1 pl-6">
          <li>To create your account and let you log in.</li>
          <li>To run games and display results and leaderboards to other players.</li>
          <li>To let you chat and play with other users in real time.</li>
          <li>To keep the service secure and prevent abuse.</li>
        </ul>
        <p className="mt-2">We do not sell your data and we do not use it for advertising.</p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">4. What other users can see</h2>
        <p>
          Your username, avatar, title and game scores are visible to other players
          (for example on the leaderboard and in lobbies). Chat messages are visible to
          their recipients. Your email address and password are never shown to other users.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">5. Storage and security</h2>
        <p>
          Your data is stored in our database. Passwords are hashed and never stored in
          plain text. We use HTTPS and standard security measures, but no system is
          completely secure.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">6. Cookies and sessions</h2>
        <p>
          We only use cookies or local storage that are strictly necessary to keep you
          logged in. We do not use tracking or advertising cookies.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">7. Data retention</h2>
        <p>
          We keep your data as long as your account exists. If you delete your account,
          your personal data is removed from our database.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">8. Your rights</h2>
        <p>
          Under the GDPR, you have the right to access, correct, export or delete your
          personal data, and to object to its processing. To exercise these rights,
          contact us at the address below.
        </p>
      </section>

      <section className="mb-6">
		<h2 className="mb-2 text-xl font-semibold">9. Contact</h2>
		<p>
			For any question about this policy or to exercise your rights, open an issue on
			our GitHub repository or contact the project team through the school.
		</p>
		</section>

      <section>
        <h2 className="mb-2 text-xl font-semibold">10. Changes</h2>
        <p>
          We may update this policy. The date at the top of the page shows the latest
          revision.
        </p>
      </section>
    </main>
  );
};
