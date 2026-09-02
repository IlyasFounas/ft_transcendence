import { useAuth } from "../context/AuthContext"

export default function HomePage() {

    const {isAuthenticated, accessToken} = useAuth();

    return (
        <div>
            <h1>Home page</h1>
            <span>isAuthenticated =</span>
            {isAuthenticated ? (
                <span>true</span>
                ) : (
                <span>false</span>)}
            <h1>AccessToken =</h1><h1>{accessToken}</h1>
        </div>
    );
}