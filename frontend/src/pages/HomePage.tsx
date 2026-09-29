import PrimaryButton from "../components/buttons/PrimaryButton";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

export default function HomePage() {
    const {isAuthenticated, authFetch} = useAuth();
	const handlePlaySolo = async () => {
	const body = new URLSearchParams({
			timer: "60",
			round_total: "5", 
		});

		const result = await authFetch("/api/lobby/create", {
			method: "POST",
			headers: { "Content-Type": "application/x-www-form-urlencoded" },
			body: body.toString(),
		});

		if (result?.status === 200) {
			navigate("/game/solo"); 
		}
	};

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
				<PrimaryButton onClick={handlePlaySolo}>Play</PrimaryButton>
            </div>
        </div>
    );
}
