import AuthenticationCard from "../components/cards/AuthenticationCard";
import AuthenticationLogin from "../components/authentication/AuthenticationLogin";

export default function AuthenticationPage () {
    return (
        <AuthenticationCard>
            <AuthenticationLogin onBackToOptions={() => {}}/>
        </AuthenticationCard>
    );
}