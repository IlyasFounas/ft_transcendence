import { useAuth } from "../context/AuthContext";

export default function HomePage() {
    const {
        isAuthenticated,
        accessToken,
        refreshAccessToken,
    } = useAuth();

    return (
        <div>
            <h1>Home page</h1>

            <div>
                <span>isAuthenticated = </span>

                {isAuthenticated ? (
                  <span>true</span>
                ) : (
                  <span>false</span>
                )}
            </div>

            <div>
                <span>AccessToken = </span>
                <p>{accessToken}</p>
            </div>

            <button
                type="button"
                onClick={() => {
                  console.log("refresh button clicked");
                  refreshAccessToken();
                }}
            >
                Refresh Access Token
            </button>
        </div>
    );
}
