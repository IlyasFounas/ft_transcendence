import { Link } from 'react-router-dom';

export default function TermsPage () {
	
  return (
    <main className="mx-auto max-w-3xl px-6 py-10 text-gray-200">
      <Link to="/" className="mb-6 inline-block text-sm text-gray-400 hover:text-white">
        ← Back to home
      </Link>

      <h1 className="mb-2 text-3xl font-bold">Terms of Service</h1>
      <p className="mb-8 text-sm text-gray-400">Last updated: October 8, 2026</p>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">1. Acceptance of the terms</h2>
        <p>
          By creating an account or using ATLAS, you agree to these Terms of Service and
          to our Privacy Policy. If you do not agree, please do not use the service.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">2. About the service</h2>
        <p>
          ATLAS is a multiplayer geography guessing game developed as a school project
          at École 42. It is provided for educational and entertainment purposes, free of
          charge and "as is".
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">3. Your account</h2>
        <ul className="list-disc space-y-1 pl-6">
          <li>You must provide accurate information when creating your account.</li>
          <li>You are responsible for keeping your password confidential.</li>
          <li>You are responsible for all activity that happens under your account.</li>
          <li>One person should use only one account.</li>
        </ul>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">4. Acceptable use</h2>
        <p>When using ATLAS, you agree not to:</p>
        <ul className="mt-2 list-disc space-y-1 pl-6">
          <li>Cheat, exploit bugs, or use automated tools to gain an unfair advantage.</li>
          <li>Harass, insult or threaten other players, including in the chat.</li>
          <li>Post offensive, illegal or inappropriate content, including in usernames and avatars.</li>
          <li>Try to access other users' accounts or disrupt the service or its security.</li>
          <li>Impersonate another person.</li>
        </ul>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">5. User content</h2>
        <p>
          You remain responsible for the content you share (username, avatar, chat
          messages). By sharing it, you allow ATLAS to display it to other users as part
          of the service. We may remove content that breaks these terms.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">6. Suspension and termination</h2>
        <p>
          We may suspend or delete an account that breaks these terms. You can stop using
          the service at any time.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">7. Intellectual property</h2>
        <p>
          The ATLAS name, design and code belong to the project team. Map and imagery
          data used in the game belong to their respective providers and are used under
          their own terms.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">8. Availability and liability</h2>
        <p>
          As a school project, ATLAS may be unavailable, modified or shut down at any
          time, and bugs may occur. We do not guarantee that the service will be
          uninterrupted or error-free, and we are not liable for any loss resulting from
          its use, to the extent permitted by law.
        </p>
      </section>

      <section className="mb-6">
        <h2 className="mb-2 text-xl font-semibold">9. Privacy</h2>
        <p>
          How we handle your personal data is described in our{' '}
          <Link to="/privacy" className="underline hover:text-white">
            Privacy Policy
          </Link>
          .
        </p>
      </section>

      <section>
        <h2 className="mb-2 text-xl font-semibold">10. Changes to these terms</h2>
        <p>
          We may update these terms. Continuing to use ATLAS after a change means you
          accept the new version.
        </p>
      </section>
    </main>
  );
};
