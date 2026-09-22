import PrimaryButton from "../components/buttons/PrimaryButton";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

export default function HomePage() {
    const {
        isAuthenticated,
    } = useAuth();

	const navigate = useNavigate();

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
				<PrimaryButton onClick={() => navigate("/game")}>Play</PrimaryButton>
            </div>
        </div>
    );
}
