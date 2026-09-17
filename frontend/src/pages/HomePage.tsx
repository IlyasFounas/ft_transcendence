import { useAuth } from "../context/AuthContext";

export default function HomePage() {
    const {
        isAuthenticated,
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
        </div>
    );
}
